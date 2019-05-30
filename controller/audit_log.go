package controller

import (
	"fmt"
	"strings"

	"github.com/fabric8-services/admin-console/app"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/auditlog"
	"github.com/fabric8-services/admin-console/configuration"
	"github.com/fabric8-services/fabric8-common/auth"
	"github.com/fabric8-services/fabric8-common/errors"
	"github.com/fabric8-services/fabric8-common/httpsupport"
	"github.com/fabric8-services/fabric8-common/log"

	"github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
)

// AuditLogsController implements the auditlogs resource.
type AuditLogsController struct {
	*goa.Controller
	db     application.DB
	config *configuration.Configuration
}

// NewAuditLogsController creates a auditlogs controller.
func NewAuditLogsController(service *goa.Service, config *configuration.Configuration, db application.DB) *AuditLogsController {
	return &AuditLogsController{
		Controller: service.NewController("AuditLogsController"),
		config:     config,
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
	eventTypeID, found := auditlog.EventTypesByName[ctx.Payload.Data.Attributes.EventType]
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

// ListForUser lists the audit logs for a given user
func (c *AuditLogsController) ListForUser(ctx *app.ListForUserAuditLogContext) error {
	// check the token and make sure it belongs to a Red Hat employee
	// retrieve the user's identity ID from the token
	username, emailAddress, emailVerified, err := parseToken(ctx)
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "unable to parse request token")
		return app.JSONErrorResponse(ctx, err)
	}
	if !(strings.HasSuffix(emailAddress, "@redhat.com") && emailVerified) {
		log.Error(ctx, map[string]interface{}{
			"target_username":     ctx.Username,
			"requesting_username": username,
			"email_address":       emailAddress,
			"email_verified":      emailVerified,
		}, "user is not allowed to list audit logs")
		return app.JSONErrorResponse(ctx, errors.NewForbiddenError("forbidden"))
	}
	// log an audit log for the current user for her action
	record := auditlog.AuditLog{
		EventTypeID: auditlog.ListAuditLogs,
		Username:    username,
		EventParams: auditlog.EventParams{
			"user": ctx.Username,
		},
	}
	err = application.Transactional(c.db, func(appl application.Application) error {
		return appl.AuditLogs().Create(ctx, &record)
	})
	// search for audit logs for the request (target) user
	var logs []auditlog.AuditLog
	var total int
	err = application.Transactional(c.db, func(appl application.Application) error {
		var err error
		logs, total, err = appl.AuditLogs().ListByUsername(ctx, ctx.Username, ctx.PageNumber, ctx.PageSize)
		return err
	})
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err":      err,
			"username": ctx.Username,
		}, "unable to list auditlogs for user")
		return app.JSONErrorResponse(ctx, err)
	}
	return ctx.OK(convertAuditLogs(ctx, logs, total, c.config))
}

func parseToken(ctx *app.ListForUserAuditLogContext) (string, string, bool, error) {
	token := goajwt.ContextJWT(ctx)
	if token == nil {
		return "", "", false, errors.NewUnauthorizedError("bad or missing token")
	}
	claims := token.Claims.(jwt.MapClaims)
	username, ok := claims["preferred_username"].(string)
	if !ok {
		log.Error(ctx, map[string]interface{}{}, "token is missing 'preferred_username' claim")
		return "", "", false, errors.NewUnauthorizedError("bad or missing token")
	}
	emailAddress, ok := claims["email"].(string)
	if !ok {
		log.Error(ctx, map[string]interface{}{}, "token is missing 'email' claim")
		return "", "", false, errors.NewUnauthorizedError("bad or missing token")
	}
	emailVerified, ok := claims["email_verified"].(bool)
	if !ok {
		log.Error(ctx, map[string]interface{}{}, "token is missing 'email_verified' claim")
		return "", "", false, errors.NewUnauthorizedError("bad or missing token")
	}
	return username, emailAddress, emailVerified, nil
}

// convertAuditLogs converts the audit logs to their resource-API counterpart
func convertAuditLogs(ctx *app.ListForUserAuditLogContext, logs []auditlog.AuditLog, total int, config httpsupport.Configuration) *app.AuditLogList {
	data := []*app.AuditLogData{}
	for _, log := range logs {
		fmt.Printf("event type: %v -> %s\n", log.EventTypeID, auditlog.EventTypesByID[log.EventTypeID])
		fmt.Printf("event params: %v\n", log.EventParams)
		data = append(data, &app.AuditLogData{
			Type: "audit_logs",
			Attributes: &app.AuditLogDataAttributes{
				Date:        log.CreatedAt.Format("2006-01-02:15:03:04"),
				EventType:   auditlog.EventTypesByID[log.EventTypeID],
				EventParams: log.EventParams,
			},
		})
	}
	response := &app.AuditLogList{
		Data:  data,
		Links: &app.PagingLinks{},
		Meta: &app.UserListMeta{
			TotalCount: total,
		},
	}
	pageNumber, pageSize := computePagingLimits(ctx.PageNumber, ctx.PageSize)
	path := httpsupport.AbsoluteURL(ctx.RequestData, ctx.Request.URL.Path, config)
	setPagingLinks(response.Links, path, len(logs), pageNumber, pageSize, total)
	return response
}
