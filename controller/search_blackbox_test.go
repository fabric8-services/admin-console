package controller_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	apptest "github.com/fabric8-services/admin-console/app/test"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/auditlog"
	"github.com/fabric8-services/admin-console/configuration"
	"github.com/fabric8-services/admin-console/controller"
	testconfig "github.com/fabric8-services/admin-console/test/generated/configuration"
	commonconfig "github.com/fabric8-services/fabric8-common/configuration"
	"github.com/fabric8-services/fabric8-common/httpsupport"
	"github.com/fabric8-services/fabric8-common/resource"
	testauth "github.com/fabric8-services/fabric8-common/test/auth"
	testsuite "github.com/fabric8-services/fabric8-common/test/suite"

	"github.com/goadesign/goa"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"
)

func newSearchController(config controller.SearchControllerConfiguration, db application.DB, options ...httpsupport.HTTPProxyOption) (*goa.Service, *controller.SearchController) {
	svc := goa.New("search")
	ctrl := controller.NewSearchController(svc,
		config,
		db,
		options...,
	)
	return svc, ctrl
}

type SearchControllerBlackboxTestSuite struct {
	testsuite.DBTestSuite
	app *application.GormApplication
}

func (s *SearchControllerBlackboxTestSuite) SetupSuite() {
	s.DBTestSuite.SetupSuite()
	s.app = application.NewGormApplication(s.DB)
}

func TestSearchController(t *testing.T) {
	resource.Require(t, resource.Database)
	config := configuration.New()
	suite.Run(t, &SearchControllerBlackboxTestSuite{DBTestSuite: testsuite.NewDBTestSuite(config)})
}

func (s *SearchControllerBlackboxTestSuite) TestSearchUsers() {

	// given
	config := testconfig.NewManagerConfigurationMock(s.T())
	config.GetAuthServiceURLFunc = func() string {
		return "https://test-auth"
	}
	config.GetDevModePrivateKeyFunc = func() []byte {
		return []byte(commonconfig.DevModeRsaPrivateKey)
	}
	svc, ctrl := newSearchController(config, s.app)
	defer gock.OffAll()

	s.T().Run("ok", func(t *testing.T) {
		// given
		ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
		require.NoError(t, err)
		tk := goajwt.ContextJWT(ctx)
		require.NotNil(t, tk)
		authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
		gock.Observe(gock.DumpRequest)
		gock.New("https://test-auth").
			Get("/api/search/users").
			MatchHeader("Authorization", authzHeader).
			MatchParam("q", "foo").
			Reply(http.StatusOK).
			BodyString(`{"data":"whatever"}`)

		// when
		apptest.SearchUsersSearchOK(t, ctx, svc, ctrl, nil, nil, "foo", &authzHeader)
		// then check that an audit record was created
		assertAuditLog(t, s.DB, *identity, auditlog.UserSearch, auditlog.EventParams{"query": "foo"})
	})

	s.T().Run("failures", func(t *testing.T) {

		t.Run("missing JWT", func(t *testing.T) {
			// given
			gock.New("http://test-tenant").
				Get("/api/search/users?q=foo").
				Reply(http.StatusUnauthorized)
			ctx := context.Background() // context is missing a JWT
			// when/then
			apptest.SearchUsersSearchUnauthorized(t, ctx, svc, ctrl, nil, nil, "foo", nil)
		})
	})
}
