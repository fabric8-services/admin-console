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

// TenantsUpdateController implements the TenantsUpdate resource.
type TenantsUpdateController struct {
	*goa.Controller
	config TenantsUpdateControllerConfiguration
	db     application.DB
}

// TenantsUpdateControllerConfiguration the configuration for the SearchController
type TenantsUpdateControllerConfiguration interface {
	GetTenantServiceURL() string
}

// NewTenantsUpdateController creates a TenantsUpdate controller.
func NewTenantsUpdateController(service *goa.Service, config TenantsUpdateControllerConfiguration, db application.DB) *TenantsUpdateController {
	return &TenantsUpdateController{
		Controller: service.NewController("TenantsUpdateController"),
		config:     config,
		db:         db,
	}
}

// Show returns information about the ongoing tenant update
func (c *TenantsUpdateController) Show(ctx *app.ShowTenantsUpdateContext) error {
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
func (c *TenantsUpdateController) Start(ctx *app.StartTenantsUpdateContext) error {
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
func (c *TenantsUpdateController) Stop(ctx *app.StopTenantsUpdateContext) error {
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
