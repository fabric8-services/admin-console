package controller

import (
	"github.com/fabric8-services/admin-console/app"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/auditlog"
	"github.com/fabric8-services/fabric8-common/errors"
	"github.com/fabric8-services/fabric8-common/httpsupport"
	"github.com/fabric8-services/fabric8-common/log"
	"github.com/fabric8-services/fabric8-common/token"

	"github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
	"github.com/satori/go.uuid"
)

// SearchController implements the search resource.
type SearchController struct {
	*goa.Controller
	config      SearchControllerConfiguration
	tokenParser token.Parser
	db          application.DB
	options     []httpsupport.HTTPProxyOption
}

// SearchControllerConfiguration the configuration for the SearchController
type SearchControllerConfiguration interface {
	GetAuthServiceURL() string
}

// NewSearchController creates a search controller.
func NewSearchController(service *goa.Service, config SearchControllerConfiguration, db application.DB, tokenParser token.Parser, options ...httpsupport.HTTPProxyOption) *SearchController {
	return &SearchController{
		Controller:  service.NewController("SearchController"),
		config:      config,
		db:          db,
		tokenParser: tokenParser,
		options:     options,
	}
}

// SearchUsers runs the search_users action.
func (c *SearchController) SearchUsers(ctx *app.SearchUsersSearchContext) error {
	jwtToken := goajwt.ContextJWT(ctx)
	if jwtToken == nil {
		log.Error(ctx, map[string]interface{}{}, "No JWT found in the request.")
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("missing authorization token"))
	}
	tk, err := c.tokenParser.Parse(ctx, jwtToken.Raw)
	if err != nil {
		log.Error(ctx, map[string]interface{}{"error": err.Error()}, "error while parsing the user's token")
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid authorization token (unable to parse)"))
	}
	var identityID uuid.UUID
	if claims, ok := tk.Claims.(jwt.MapClaims); ok {
		var err error
		sub := claims["sub"].(string)
		identityID, err = uuid.FromString(sub)
		if err != nil {
			return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid authorization token (invalid 'sub' claim)"))
		}
	} else {
		return app.JSONErrorResponse(ctx, errors.NewUnauthorizedError("invalid authorization token (invalid type of claims)"))
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
