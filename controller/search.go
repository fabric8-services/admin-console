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

// SearchController implements the search resource.
type SearchController struct {
	*goa.Controller
	config  SearchControllerConfiguration
	db      application.DB
	options []httpsupport.HTTPProxyOption
}

// SearchControllerConfiguration the configuration for the SearchController
type SearchControllerConfiguration interface {
	GetAuthServiceURL() string
}

// NewSearchController creates a search controller.
func NewSearchController(service *goa.Service, config SearchControllerConfiguration, db application.DB, options ...httpsupport.HTTPProxyOption) *SearchController {
	return &SearchController{
		Controller: service.NewController("SearchController"),
		config:     config,
		db:         db,
		options:    options,
	}
}

// SearchUsers runs the search_users action.
func (c *SearchController) SearchUsers(ctx *app.SearchUsersSearchContext) error {
	identityID, username, err := authsupport.LocateIdentity(ctx)
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "missing or invalid authorization token")
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("missing or invalid authorization token"))
	}
	record := auditlog.AuditLog{
		EventTypeID: auditlog.UserSearch,
		IdentityID:  identityID,
		Username:    username,
		EventParams: auditlog.EventParams{
			"query": ctx.Q,
		},
	}
	err = application.Transactional(c.db, func(appl application.Application) error {
		return appl.AuditLogs().Create(ctx, &record)
	})
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "unable to record the auditlog while proxying request to auth")
		return app.JSONErrorResponse(ctx, err)
	}

	return httpsupport.RouteHTTP(ctx, c.config.GetAuthServiceURL(), c.options...)
}
