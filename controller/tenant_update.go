package controller

import (
	"github.com/fabric8-services/admin-console/app"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/auditlog"
	authsupport "github.com/fabric8-services/fabric8-common/auth"
	"github.com/fabric8-services/fabric8-common/errors"
	"github.com/fabric8-services/fabric8-common/httpsupport"
	"github.com/fabric8-services/fabric8-common/log"
	"github.com/goadesign/goa"
)

// TenantUpdateController implements the TenantUpdate resource.
type TenantUpdateController struct {
	*goa.Controller
	config TenantUpdateControllerConfiguration
	db     application.DB
}

// TenantUpdateControllerConfiguration the configuration for the SearchController
type TenantUpdateControllerConfiguration interface {
	GetTenantServiceURL() string
}

// NewTenantUpdateController creates a TenantUpdate controller.
func NewTenantUpdateController(service *goa.Service, config TenantUpdateControllerConfiguration, db application.DB) *TenantUpdateController {
	return &TenantUpdateController{
		Controller: service.NewController("TenantUpdateController"),
		config:     config,
		db:         db,
	}
}

// Show returns information about the ongoing tenant update
func (c *TenantUpdateController) Show(ctx *app.ShowTenantUpdateContext) error {
	identityID, err := authsupport.LocateIdentity(ctx)
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "invalid or missing authorization token")
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid or missing authorization token"))
	}
	err = application.Transactional(c.db, func(appl application.Application) error {
		eventParams := auditlog.EventParams{}
		if ctx.ClusterURL != nil {
			eventParams["clusterURL"] = *ctx.ClusterURL
		}
		if ctx.EnvType != nil {
			eventParams["envType"] = *ctx.EnvType
		}
		return appl.AuditLogs().Create(ctx, &auditlog.AuditLog{
			EventTypeID: auditlog.ShowTenantUpdate,
			IdentityID:  identityID,
			EventParams: eventParams,
		})
	})
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "unable to record the auditlog while proxying request to tenant")
		return app.JSONErrorResponse(ctx, err)
	}
	return httpsupport.RouteHTTPToPath(ctx, c.config.GetTenantServiceURL(), "/api/update")
}

// Start starts a tenant update
func (c *TenantUpdateController) Start(ctx *app.StartTenantUpdateContext) error {
	identityID, err := authsupport.LocateIdentity(ctx)
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "invalid or missing authorization token")
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid or missing authorization token"))
	}
	err = application.Transactional(c.db, func(appl application.Application) error {
		eventParams := auditlog.EventParams{}
		if ctx.ClusterURL != nil {
			eventParams["clusterURL"] = *ctx.ClusterURL
		}
		if ctx.EnvType != nil {
			eventParams["envType"] = *ctx.EnvType
		}
		return appl.AuditLogs().Create(ctx, &auditlog.AuditLog{
			EventTypeID: auditlog.StartTenantUpdate,
			IdentityID:  identityID,
			EventParams: eventParams,
		})
	})
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "unable to record the auditlog while proxying request to tenant")
		return app.JSONErrorResponse(ctx, err)
	}
	return httpsupport.RouteHTTPToPath(ctx, c.config.GetTenantServiceURL(), "/api/update")
}

// Stop stops the ongoing tenant update
func (c *TenantUpdateController) Stop(ctx *app.StopTenantUpdateContext) error {
	identityID, err := authsupport.LocateIdentity(ctx)
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "invalid or missing authorization token")
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid or missing authorization token"))
	}
	err = application.Transactional(c.db, func(appl application.Application) error {
		return appl.AuditLogs().Create(ctx, &auditlog.AuditLog{
			EventTypeID: auditlog.StopTenantUpdate,
			IdentityID:  identityID,
			EventParams: auditlog.EventParams{},
		})
	})
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "unable to record the auditlog while proxying request to tenant")
		return app.JSONErrorResponse(ctx, err)
	}
	return httpsupport.RouteHTTPToPath(ctx, c.config.GetTenantServiceURL(), "/api/update")
}
