package configuration

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	commonconfig "github.com/fabric8-services/fabric8-common/configuration"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// String returns the current configuration as a string
func (c *Configuration) String() string {
	allSettings := c.v.AllSettings()
	y, err := yaml.Marshal(&allSettings)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"settings": allSettings,
			"err":      err,
		}).Panicln("Failed to marshall config to string")
	}
	return fmt.Sprintf("%s\n", y)
}

const (
	// Constants for viper variable names. Will be used to set
	// default values as well as to get each value

	// General
	varHTTPAddress                        = "http.address"
	varHeaderMaxLength                    = "header.maxlength"
	varMetricsHTTPAddress                 = "metrics.http.address"
	varDeveloperModeEnabled               = "developer.mode.enabled"
	varCleanTestDataEnabled               = "clean.test.data"
	varCleanTestDataErrorReportingEnabled = "clean.test.data.error.reporting"
	varDBLogsEnabled                      = "enable.db.logs"
	varLogLevel                           = "log.level"
	varLogJSON                            = "log.json"

	// Postgres
	varPostgresHost                 = "postgres.host"
	varPostgresPort                 = "postgres.port"
	varPostgresUser                 = "postgres.user"
	varPostgresDatabase             = "postgres.database"
	varPostgresPassword             = "postgres.password"
	varPostgresSSLMode              = "postgres.sslmode"
	varPostgresConnectionTimeout    = "postgres.connection.timeout"
	varPostgresTransactionTimeout   = "postgres.transaction.timeout"
	varPostgresConnectionRetrySleep = "postgres.connection.retrysleep"
	varPostgresConnectionMaxIdle    = "postgres.connection.maxidle"
	varPostgresConnectionMaxOpen    = "postgres.connection.maxopen"

	varDiagnoseHTTPAddress = "diagnose.http.address"

	// sentry
	varEnvironment = "environment"
	varSentryDSN   = "sentry.dsn"

	// other services
	varAuthURL   = "auth.url"
	varTenantURL = "tenant.url"
)

// Configuration encapsulates the Viper configuration object which stores the configuration data in-memory.
type Configuration struct {
	// Main Configuration
	v                         *viper.Viper
	defaultConfigurationError error
	mux                       sync.RWMutex
}

// New creates a configuration reader object using configurable configuration file paths
func New() *Configuration {
	c := &Configuration{
		v: viper.New(),
	}

	// Set up the main configuration
	c.v.SetEnvPrefix("ADMIN")
	c.v.AutomaticEnv()
	c.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.v.SetTypeByDefaultValue(true)
	c.setConfigDefaults()

	// Check sensitive default configuration
	hostname := c.validateURL(c.GetAuthServiceURL(), "Auth service")
	if hostname == "localhost" {
		c.appendDefaultConfigErrorMessage("Auth service url is localhost")
	}
	if c.GetSentryDSN() == "" {
		c.appendDefaultConfigErrorMessage("Sentry DSN is empty")
	}
	if c.GetPostgresPassword() == defaultDBPassword {
		c.appendDefaultConfigErrorMessage("default DB password is used")
	}
	if c.IsDeveloperModeEnabled() {
		c.appendDefaultConfigErrorMessage("developer mode is enabled")
	}

	return c
}

// returns the hostname of the given URL if this latter was not empty
func (c *Configuration) validateURL(serviceURL, serviceName string) string {
	if serviceURL == "" {
		c.appendDefaultConfigErrorMessage(fmt.Sprintf("%s url is empty", serviceName))
		return ""
	}

	u, err := url.Parse(serviceURL)
	if err != nil {
		c.appendDefaultConfigErrorMessage(fmt.Sprintf("invalid %s url: %s", serviceName, err.Error()))
		return ""
	}
	if u.Hostname() == "" { // probably missing the http/https scheme
		c.appendDefaultConfigErrorMessage(fmt.Sprintf("invalid %s url (missing scheme?)", serviceName))
		return ""
	}
	return u.Hostname()

}

func (c *Configuration) appendDefaultConfigErrorMessage(message string) {
	if c.defaultConfigurationError == nil {
		c.defaultConfigurationError = errors.New(message)
	} else {
		c.defaultConfigurationError = errors.Errorf("%s; %s", c.defaultConfigurationError.Error(), message)
	}
}

// DefaultConfigurationError returns an error if the default values is used
// for sensitive configuration like service account secrets or private keys.
// Error contains all the details.
// Returns nil if the default configuration is not used.
func (c *Configuration) DefaultConfigurationError() error {
	// Lock for reading because config file watcher can update config errors
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.defaultConfigurationError
}

// GetAuthServiceURL returns Auth Service URL
func (c *Configuration) GetAuthServiceURL() string {
	return c.v.GetString(varAuthURL)
}

// GetTenantServiceURL returns Tenant Service URL
func (c *Configuration) GetTenantServiceURL() string {
	return c.v.GetString(varTenantURL)
}

func (c *Configuration) setConfigDefaults() {
	//---------
	// Postgres
	//---------

	// We already call this in NewConfiguration() - do we need it again??
	c.v.SetTypeByDefaultValue(true)

	c.v.SetDefault(varPostgresHost, "localhost")
	c.v.SetDefault(varPostgresPort, 5435)
	c.v.SetDefault(varPostgresUser, "postgres")
	c.v.SetDefault(varPostgresDatabase, "postgres")
	c.v.SetDefault(varPostgresPassword, defaultDBPassword)
	c.v.SetDefault(varPostgresSSLMode, "disable")
	c.v.SetDefault(varPostgresConnectionTimeout, 5)
	c.v.SetDefault(varPostgresConnectionMaxIdle, -1)
	c.v.SetDefault(varPostgresConnectionMaxOpen, -1)

	// Number of seconds to wait before trying to connect again
	c.v.SetDefault(varPostgresConnectionRetrySleep, time.Duration(time.Second))

	// Timeout of a transaction in minutes
	c.v.SetDefault(varPostgresTransactionTimeout, time.Duration(5*time.Minute))

	//-----
	// HTTP
	//-----
	c.v.SetDefault(varHTTPAddress, "0.0.0.0:8089")
	c.v.SetDefault(varMetricsHTTPAddress, "0.0.0.0:8089")
	c.v.SetDefault(varHeaderMaxLength, defaultHeaderMaxLength)

	//-----
	// Sentry
	//-----
	// prod-preview or prod
	c.v.SetDefault(varEnvironment, "local")

	//-----
	// Misc
	//-----

	// Enable development related features, e.g. token generation endpoint
	c.v.SetDefault(varDeveloperModeEnabled, false)

	// By default, test data should be cleaned from DB, unless explicitely said otherwise.
	c.v.SetDefault(varCleanTestDataEnabled, true)
	// By default, DB logs are not output in the console
	c.v.SetDefault(varDBLogsEnabled, false)

	c.v.SetDefault(varLogLevel, defaultLogLevel)

}

// GetPostgresHost returns the postgres host as set via default, config file, or environment variable
func (c *Configuration) GetPostgresHost() string {
	return c.v.GetString(varPostgresHost)
}

// GetPostgresPort returns the postgres port as set via default, config file, or environment variable
func (c *Configuration) GetPostgresPort() int64 {
	return c.v.GetInt64(varPostgresPort)
}

// GetPostgresUser returns the postgres user as set via default, config file, or environment variable
func (c *Configuration) GetPostgresUser() string {
	return c.v.GetString(varPostgresUser)
}

// GetPostgresDatabase returns the postgres database as set via default, config file, or environment variable
func (c *Configuration) GetPostgresDatabase() string {
	return c.v.GetString(varPostgresDatabase)
}

// GetPostgresPassword returns the postgres password as set via default, config file, or environment variable
func (c *Configuration) GetPostgresPassword() string {
	return c.v.GetString(varPostgresPassword)
}

// GetPostgresSSLMode returns the postgres sslmode as set via default, config file, or environment variable
func (c *Configuration) GetPostgresSSLMode() string {
	return c.v.GetString(varPostgresSSLMode)
}

// GetPostgresConnectionTimeout returns the postgres connection timeout as set via default, config file, or environment variable
func (c *Configuration) GetPostgresConnectionTimeout() int64 {
	return c.v.GetInt64(varPostgresConnectionTimeout)
}

// GetPostgresConnectionRetrySleep returns the number of seconds (as set via default, config file, or environment variable)
// to wait before trying to connect again
func (c *Configuration) GetPostgresConnectionRetrySleep() time.Duration {
	return c.v.GetDuration(varPostgresConnectionRetrySleep)
}

// GetPostgresTransactionTimeout returns the number of minutes to timeout a transaction
func (c *Configuration) GetPostgresTransactionTimeout() time.Duration {
	return c.v.GetDuration(varPostgresTransactionTimeout)
}

// GetPostgresConnectionMaxIdle returns the number of connections that should be keept alive in the database connection pool at
// any given time. -1 represents no restrictions/default behavior
func (c *Configuration) GetPostgresConnectionMaxIdle() int {
	return c.v.GetInt(varPostgresConnectionMaxIdle)
}

// GetPostgresConnectionMaxOpen returns the max number of open connections that should be open in the database connection pool.
// -1 represents no restrictions/default behavior
func (c *Configuration) GetPostgresConnectionMaxOpen() int {
	return c.v.GetInt(varPostgresConnectionMaxOpen)
}

// GetPostgresConfigString returns a ready to use string for usage in sql.Open()
func (c *Configuration) GetPostgresConfigString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		c.GetPostgresHost(),
		c.GetPostgresPort(),
		c.GetPostgresUser(),
		c.GetPostgresPassword(),
		c.GetPostgresDatabase(),
		c.GetPostgresSSLMode(),
		c.GetPostgresConnectionTimeout(),
	)
}

// GetHTTPAddress returns the HTTP address (as set via default, config file, or environment variable)
// that the auth server binds to (e.g. "0.0.0.0:8089")
func (c *Configuration) GetHTTPAddress() string {
	return c.v.GetString(varHTTPAddress)
}

// GetMetricsHTTPAddress returns the address the /metrics endpoing will be mounted.
// By default GetMetricsHTTPAddress is the same as GetHTTPAddress
func (c *Configuration) GetMetricsHTTPAddress() string {
	return c.v.GetString(varMetricsHTTPAddress)
}

// GetHeaderMaxLength returns the max length of HTTP headers allowed in the system
// For example it can be used to limit the size of bearer tokens returned by the api service
func (c *Configuration) GetHeaderMaxLength() int64 {
	return c.v.GetInt64(varHeaderMaxLength)
}

// IsDeveloperModeEnabled returns if development related features (as set via default, config file, or environment variable),
// e.g. token generation endpoint are enabled
func (c *Configuration) IsDeveloperModeEnabled() bool {
	return c.v.GetBool(varDeveloperModeEnabled)
}

// IsPostgresDeveloperModeEnabled returns if development related features (as set via default, config file, or environment variable),
// e.g. token generation endpoint are enabled
func (c *Configuration) IsPostgresDeveloperModeEnabled() bool {
	// for backward compatibility with fabric8-common
	return false
}

// IsCleanTestDataEnabled returns `true` if the test data should be cleaned after each test. (default: true)
func (c *Configuration) IsCleanTestDataEnabled() bool {
	return c.v.GetBool(varCleanTestDataEnabled)
}

// IsCleanTestDataErrorReportingRequired returns `true` if the test data should be cleaned after each test. (default: true)
func (c *Configuration) IsCleanTestDataErrorReportingRequired() bool {
	return c.v.GetBool(varCleanTestDataErrorReportingEnabled)
}

// IsDBLogsEnabled returns `true` if the DB logs (ie, SQL queries) should be output in the console. (default: false)
func (c *Configuration) IsDBLogsEnabled() bool {
	return c.v.GetBool(varDBLogsEnabled)
}

// GetSentryDSN returns the secret needed to securely communicate with https://errortracking.prod-preview.openshift.io/openshift_io/admin-console
func (c *Configuration) GetSentryDSN() string {
	return c.v.GetString(varSentryDSN)
}

// GetLogLevel returns the logging level (as set via config file or environment variable)
func (c *Configuration) GetLogLevel() string {
	return c.v.GetString(varLogLevel)
}

// IsLogJSON returns if we should log json format (as set via config file or environment variable)
func (c *Configuration) IsLogJSON() bool {
	if c.v.IsSet(varLogJSON) {
		return c.v.GetBool(varLogJSON)
	}
	if c.IsDeveloperModeEnabled() {
		return false
	}
	return true
}

// GetEnvironment returns the current environment application is deployed in
// like 'production', 'prod-preview', 'local', etc as the value of environment variable
// `AUTH_ENVIRONMENT` is set.
func (c *Configuration) GetEnvironment() string {
	return c.v.GetString(varEnvironment)
}

// GetDiagnoseHTTPAddress returns the address of where to start the gops handler.
// By default GetDiagnoseHTTPAddress is 127.0.0.1:0 in devMode, but turned off in prod mode
// unless explicitly configured
func (c *Configuration) GetDiagnoseHTTPAddress() string {
	if c.v.IsSet(varDiagnoseHTTPAddress) {
		return c.v.GetString(varDiagnoseHTTPAddress)
	} else if c.IsDeveloperModeEnabled() {
		return "127.0.0.1:0"
	}
	return ""
}

// GetDevModePrivateKey returns additional public key which should be used by the admin console service in Dev Mode
// Returns an error if the application is not running in dev mode
func (c *Configuration) GetDevModePrivateKey() []byte {
	return []byte(commonconfig.DevModeRsaPrivateKey)
}
