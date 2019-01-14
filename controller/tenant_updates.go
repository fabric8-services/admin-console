package controller

import (
	"github.com/fabric8-services/admin-console/app"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/auditlog"
	"github.com/fabric8-services/fabric8-common/errors"
	"github.com/fabric8-services/fabric8-common/httpsupport"
	"github.com/fabric8-services/fabric8-common/token"
	"github.com/goadesign/goa"
)

// TenantUpdatesController implements the TenantUpdates resource.
type TenantUpdatesController struct {
	*goa.Controller
	config TenantUpdatesControllerConfiguration
	db     application.DB
}

// TenantUpdatesControllerConfiguration the configuration for the SearchController
type TenantUpdatesControllerConfiguration interface {
	GetTenantServiceURL() string
}

// NewTenantUpdatesController creates a TenantUpdates controller.
func NewTenantUpdatesController(service *goa.Service, config TenantUpdatesControllerConfiguration, db application.DB) *TenantUpdatesController {
	return &TenantUpdatesController{
		Controller: service.NewController("TenantUpdatesController"),
		config:     config,
		db:         db,
	}
}

// Show returns information about the ongoing tenant update
func (c *TenantUpdatesController) Show(ctx *app.ShowTenantUpdatesContext) error {
	tokenManager, err := token.ReadManagerFromContext(ctx)
	if err != nil {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("missing token manager in the request context"))
	}
	identityID, err := tokenManager.Locate(ctx)
	if err != nil {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid authorization token (invalid 'sub' claim)"))
	}
	record := auditlog.AuditLog{
		EventTypeID: auditlog.ShowTenantUpdate,
		IdentityID:  identityID,
		EventParams: auditlog.EventParams{},
	}
	err = application.Transactional(c.db, func(appl application.Application) error {
		return appl.AuditLogs().Create(ctx, &record)
	})
	if err != nil {
		return app.JSONErrorResponse(ctx, err)
	}
	return httpsupport.RouteHTTP(ctx, c.config.GetTenantServiceURL())
}

// Start starts a tenant update
func (c *TenantUpdatesController) Start(ctx *app.StartTenantUpdatesContext) error {
	tokenManager, err := token.ReadManagerFromContext(ctx)
	if err != nil {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("missing token manager in the request context"))
	}
	identityID, err := tokenManager.Locate(ctx)
	if err != nil {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid authorization token (invalid 'sub' claim)"))
	}
	eventParams := auditlog.EventParams{}
	if ctx.ClusterURL != nil {
		eventParams["clusterURL"] = *ctx.ClusterURL
	}
	if ctx.EnvType != nil {
		eventParams["envType"] = *ctx.EnvType
	}
	record := auditlog.AuditLog{
		EventTypeID: auditlog.StartTenantUpdate,
		IdentityID:  identityID,
		EventParams: eventParams,
	}
	err = application.Transactional(c.db, func(appl application.Application) error {
		return appl.AuditLogs().Create(ctx, &record)
	})
	if err != nil {
		return app.JSONErrorResponse(ctx, err)
	}
	return httpsupport.RouteHTTP(ctx, c.config.GetTenantServiceURL())
}

// Stop stops the ongoing tenant update
func (c *TenantUpdatesController) Stop(ctx *app.StopTenantUpdatesContext) error {
	tokenManager, err := token.ReadManagerFromContext(ctx)
	if err != nil {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("missing token manager in the request context"))
	}
	identityID, err := tokenManager.Locate(ctx)
	if err != nil {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid authorization token (invalid 'sub' claim)"))
	}
	record := auditlog.AuditLog{
		EventTypeID: auditlog.StopTenantUpdate,
		IdentityID:  identityID,
		EventParams: auditlog.EventParams{},
	}
	err = application.Transactional(c.db, func(appl application.Application) error {
		return appl.AuditLogs().Create(ctx, &record)
	})
	if err != nil {
		return app.JSONErrorResponse(ctx, err)
	}
	return httpsupport.RouteHTTP(ctx, c.config.GetTenantServiceURL())
}
