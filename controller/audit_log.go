package controller

import (
	"github.com/fabric8-services/admin-console/app"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/auditlog"
	"github.com/fabric8-services/fabric8-common/auth"
	"github.com/fabric8-services/fabric8-common/errors"
	"github.com/fabric8-services/fabric8-common/log"
	"github.com/goadesign/goa"
)

// AuditLogsController implements the auditlogs resource.
type AuditLogsController struct {
	*goa.Controller
	db application.DB
}

// NewAuditLogsController creates a auditlogs controller.
func NewAuditLogsController(service *goa.Service, db application.DB) *AuditLogsController {
	return &AuditLogsController{
		Controller: service.NewController("AuditLogsController"),
		db:         db,
	}
}

// Create runs the create action.
func (c *AuditLogsController) Create(ctx *app.CreateAuditLogContext) error {
	// check the token and make sure it belongs to `auth`
	if !auth.IsSpecificServiceAccount(ctx, auth.Auth) {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid or missing authorization token"))
	}
	// lookup the event type ID
	eventTypeID, found := auditlog.EventTypes[ctx.Payload.Data.Attributes.EventType]
	if !found {
		return app.JSONErrorResponse(ctx, errors.NewBadParameterError("event_type", ctx.Payload.Data.Attributes.EventType))
	}
	log.Info(ctx, map[string]interface{}{
		"username": ctx.Username,
	}, "creating audit log for user")
	err := application.Transactional(c.db, func(appl application.Application) error {
		return appl.AuditLogs().Create(ctx, &auditlog.AuditLog{
			Username:    ctx.Username,
			EventTypeID: eventTypeID,
			EventParams: ctx.Payload.Data.Attributes.EventParams,
		})
	})
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err":      err,
			"username": ctx.Username,
		}, "unable to record the auditlog for user")
		return app.JSONErrorResponse(ctx, err)
	}

	return ctx.NoContent()
}
