package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/fabric8-services/admin-console/app"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/configuration"
	"github.com/fabric8-services/admin-console/controller"
	"github.com/fabric8-services/admin-console/migration"
	"github.com/fabric8-services/fabric8-common/closeable"
	"github.com/fabric8-services/fabric8-common/goamiddleware"
	"github.com/fabric8-services/fabric8-common/log"
	"github.com/fabric8-services/fabric8-common/metric"
	"github.com/fabric8-services/fabric8-common/sentry"
	"github.com/fabric8-services/fabric8-common/token"

	"github.com/goadesign/goa"
	goalogrus "github.com/goadesign/goa/logging/logrus"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/google/gops/agent"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/lib/pq"
)

func main() {

	// --------------------------------------------------------------------
	// Parse flags
	// --------------------------------------------------------------------
	var configFilePath string
	var printConfig bool
	var migrateDB bool

	flag.StringVar(&configFilePath, "config", "", "Path to the config file to read")
	flag.BoolVar(&printConfig, "printConfig", false, "Prints the config (including merged environment variables) and exits")
	flag.BoolVar(&migrateDB, "migrateDatabase", false, "Migrates the database to the newest version and exits.")
	flag.Parse()

	config, err := configuration.New()
	if err != nil {
		log.Panic(context.TODO(), map[string]interface{}{
			"config_file_path": configFilePath,
			"err":              err,
		}, "failed to setup the configuration")
	}

	// Initialized developer mode flag and log level for the logger
	log.InitializeLogger(config.IsLogJSON(), config.GetLogLevel())

	var db *gorm.DB
	for {
		db, err = gorm.Open("postgres", config.GetPostgresConfigString())
		if err != nil {
			log.Logger().Errorf("ERROR: Unable to open connection to database %v", err)
			closeable.Close(context.TODO(), db)
			log.Logger().Infof("Retrying to connect in %v...", config.GetPostgresConnectionRetrySleep())
			time.Sleep(config.GetPostgresConnectionRetrySleep())
		} else {
			defer closeable.Close(context.TODO(), db)
			break
		}
	}
	// Migrate the schema
	err = migration.Migrate(db.DB(), config.GetPostgresDatabase())
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed migration")
	}
	// Nothing to here except exit, since the migration is already performed.
	if migrateDB {
		os.Exit(0)
	}

	// Initialize sentry client
	haltSentry, err := sentry.InitializeSentryClient(
		nil, // will use the `os.Getenv("Sentry_DSN")` instead
		sentry.WithRelease(app.Commit),
		sentry.WithEnvironment(config.GetEnvironment()),
	)
	if err != nil {
		log.Panic(context.TODO(), map[string]interface{}{
			"err": err,
		}, "failed to setup the sentry client")
	}
	defer haltSentry()

	// Create service
	service := goa.New("admin-console")

	// Mount middleware
	service.Use(middleware.RequestID())
	// Use our own log request to inject identity id and modify other properties
	service.Use(gzip.Middleware(9))
	service.Use(app.ErrorHandler(service, true))
	service.Use(middleware.Recover())
	// record HTTP request metrics in prometh
	service.Use(
		metric.Recorder(
			"admin_console",
			metric.WithRequestDurationBucket(prometheus.ExponentialBuckets(0.05, 2, 8))))
	service.WithLogger(goalogrus.New(log.Logger()))

	appDB := application.NewGormApplication(db)
	tokenManager, err := token.DefaultManager(config)
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to create token manager")
	}

	// Middleware that extracts and stores the token in the context
	jwtMiddlewareTokenContext := goamiddleware.TokenContext(tokenManager, app.NewJWTSecurity())
	service.Use(jwtMiddlewareTokenContext)

	service.Use(token.InjectTokenManager(tokenManager))
	service.Use(log.LogRequest(config.IsDeveloperModeEnabled()))
	app.UseJWTMiddleware(service, jwt.New(tokenManager.PublicKeys(), nil, app.NewJWTSecurity()))

	// Mount the '/status' controller
	statusCtrl := controller.NewStatusController(service)
	app.MountStatusController(service, statusCtrl)

	// Mount the '/search' controller
	searchCtrl := controller.NewSearchController(service, config, appDB)
	app.MountSearchController(service, searchCtrl)

	log.Logger().Infoln("Git Commit SHA: ", app.Commit)
	log.Logger().Infoln("UTC Build Time: ", app.BuildTime)
	log.Logger().Infoln("UTC Start Time: ", app.StartTime)
	log.Logger().Infoln("GOMAXPROCS:     ", runtime.GOMAXPROCS(-1))
	log.Logger().Infoln("NumCPU:         ", runtime.NumCPU())

	printUserInfo()
	// Nothing to here except exit.
	if printConfig {
		os.Exit(0)
	}

	http.Handle("/api/", service.Mux)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	if config.GetDiagnoseHTTPAddress() != "" {
		log.Logger().Infoln("Diagnose:       ", config.GetDiagnoseHTTPAddress())
		// Start diagnostic http
		if err := agent.Listen(agent.Options{Addr: config.GetDiagnoseHTTPAddress(), ConfigDir: "/tmp/gops/"}); err != nil {
			log.Error(context.TODO(), map[string]interface{}{
				"addr": config.GetDiagnoseHTTPAddress(),
				"err":  err,
			}, "unable to connect to diagnose server")
		}
	}

	// // Start/mount metrics http
	if config.GetHTTPAddress() == config.GetMetricsHTTPAddress() {
		http.Handle("/metrics", promhttp.Handler())
	} else {
		go func(metricAddress string) {
			mx := http.NewServeMux()
			mx.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(metricAddress, mx); err != nil {
				log.Error(nil, map[string]interface{}{
					"addr": metricAddress,
					"err":  err,
				}, "unable to connect to metrics server")
				service.LogError("startup", "err", err)
			}
		}(config.GetMetricsHTTPAddress())
	}

	// Start http
	if err := http.ListenAndServe(config.GetHTTPAddress(), nil); err != nil {
		log.Error(nil, map[string]interface{}{
			"addr": config.GetHTTPAddress(),
			"err":  err,
		}, "unable to connect to server")
		service.LogError("startup", "err", err)
	}

}

func printUserInfo() {
	u, err := user.Current()
	if err != nil {
		log.Warn(nil, map[string]interface{}{
			"err": err,
		}, "failed to get current user")
	} else {
		log.Info(nil, map[string]interface{}{
			"username": u.Username,
			"uuid":     u.Uid,
		}, "Running as user name '%s' with UID %s.", u.Username, u.Uid)
		g, err := user.LookupGroupId(u.Gid)
		if err != nil {
			log.Warn(nil, map[string]interface{}{
				"err": err,
			}, "failed to lookup group")
		} else {
			log.Info(nil, map[string]interface{}{
				"groupname": g.Name,
				"gid":       g.Gid,
			}, "Running as as group '%s' with GID %s.", g.Name, g.Gid)
		}
	}

}
