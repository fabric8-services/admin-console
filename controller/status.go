package controller

import (
	"fmt"

	"github.com/fabric8-services/admin-console/app"
	"github.com/fabric8-services/fabric8-common/log"
	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
)

// StatusController implements the status resource.
type statusConfiguration interface {
	IsDeveloperModeEnabled() bool
	DefaultConfigurationError() error
}

// DBChecker is to be used to check if the DB is reachable
type DBChecker interface {
	Ping() error
}

// StatusController implements the status resource.
type StatusController struct {
	*goa.Controller
	dbChecker DBChecker
	config    statusConfiguration
}

// NewStatusController creates a status controller.
func NewStatusController(service *goa.Service, dbChecker DBChecker, config statusConfiguration) *StatusController {
	return &StatusController{
		Controller: service.NewController("StatusController"),
		dbChecker:  dbChecker,
		config:     config,
	}
}

// Show runs the show action.
func (c *StatusController) Show(ctx *app.ShowStatusContext) error {
	res := &app.Status{
		Commit:    app.Commit,
		BuildTime: app.BuildTime,
		StartTime: app.StartTime,
	}

	devMode := c.config.IsDeveloperModeEnabled()
	if devMode {
		res.DevMode = &devMode
	}

	dbErr := c.dbChecker.Ping()
	if dbErr != nil {
		log.Error(ctx, map[string]interface{}{
			"db_error": dbErr.Error(),
		}, "database configuration error")
		res.DatabaseStatus = fmt.Sprintf("Error: %s", dbErr.Error())
	} else {
		res.DatabaseStatus = "OK"
	}

	configErr := c.config.DefaultConfigurationError()
	if configErr != nil {
		log.Error(ctx, map[string]interface{}{
			"config_error": configErr.Error(),
		}, "configuration error")
		res.ConfigurationStatus = fmt.Sprintf("Error: %s", configErr.Error())
	} else {
		res.ConfigurationStatus = "OK"
	}

	if dbErr != nil || (configErr != nil && !devMode) {
		return ctx.ServiceUnavailable(res)
	}
	return ctx.OK(res)
}

// GormDBChecker implements DB checker
type GormDBChecker struct {
	db *gorm.DB
}

// NewGormDBChecker constructs a new GormDBChecker
func NewGormDBChecker(db *gorm.DB) DBChecker {
	return &GormDBChecker{
		db: db,
	}
}

// Ping performs a basic query to check that the Database connection is valid
func (c *GormDBChecker) Ping() error {
	_, err := c.db.DB().Exec("select 1")
	return err
}
