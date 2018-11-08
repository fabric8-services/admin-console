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
	tokenManager, err := token.ReadManagerFromContext(ctx)
	if err != nil {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("missing token manager in the request context"))
	}
	identityID, err := tokenManager.Locate(ctx)
	if err != nil {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid authorization token (invalid 'sub' claim)"))
	}
	record := auditlog.AuditLog{
		EventTypeID: auditlog.UserSearch,
		IdentityID:  identityID,
		EventParams: auditlog.EventParams{
			"query": ctx.Q,
		},
	}
	err = application.Transactional(c.db, func(appl application.Application) error {
		return appl.AuditLogs().Create(ctx, &record)
	})
	if err != nil {
		return app.JSONErrorResponse(ctx, err)
	}

	return httpsupport.RouteHTTP(ctx, c.config.GetAuthServiceURL(), c.options...)
}
