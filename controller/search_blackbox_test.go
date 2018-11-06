package controller_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fabric8-services/admin-console/auditlog"
	"github.com/fabric8-services/admin-console/configuration"
	"github.com/stretchr/testify/suite"

	apptest "github.com/fabric8-services/admin-console/app/test"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/controller"
	testconfig "github.com/fabric8-services/admin-console/test/generated/configuration"
	testjwt "github.com/fabric8-services/admin-console/test/jwt"
	"github.com/fabric8-services/fabric8-common/httpsupport"
	"github.com/fabric8-services/fabric8-common/resource"
	testrecorder "github.com/fabric8-services/fabric8-common/test/recorder"
	testsuite "github.com/fabric8-services/fabric8-common/test/suite"
	"github.com/fabric8-services/fabric8-common/token"

	"github.com/goadesign/goa"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func newSearchController(config controller.SearchControllerConfiguration, db application.DB, tokenParser token.Parser, options ...httpsupport.HTTPProxyOption) (*goa.Service, *controller.SearchController) {
	svc := goa.New("feature")
	ctrl := controller.NewSearchController(svc,
		config,
		db,
		tokenParser,
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
	config, err := configuration.New()
	require.NoError(t, err)
	suite.Run(t, &SearchControllerBlackboxTestSuite{DBTestSuite: testsuite.NewDBTestSuite(config)})
}

func (s *SearchControllerBlackboxTestSuite) TestSearchUsers() {

	// given
	kid := "test-key"
	config := testconfig.NewManagerConfigurationMock(s.T())
	config.GetAuthServiceURLFunc = func() string {
		return "https://test-auth"
	}
	config.GetDevModePrivateKeyFunc = func() []byte {
		return []byte(testjwt.PrivateKey1)
	}
	r, err := testrecorder.New("search_blackbox_test")
	require.NoError(s.T(), err)
	defer r.Stop()
	tokenParser, err := token.NewManager(config, httpsupport.WithRoundTripper(r))
	require.NoError(s.T(), err)
	svc, ctrl := newSearchController(config, s.app, tokenParser, httpsupport.WithProxyTransport(r))

	s.T().Run("ok", func(t *testing.T) {
		// given
		identityID := uuid.NewV4()
		require.NoError(t, err)
		ctx, err := testjwt.NewJWTContext(identityID.String(), kid, testjwt.PrivateKey1)
		require.NoError(t, err)
		// when
		apptest.SearchUsersSearchOK(t, ctx, svc, ctrl, nil, nil, "foo")
		// then check that an audit record was created
		recordRepo := auditlog.NewRepository(s.DB)
		records, total, err := recordRepo.ListByIdentityID(context.Background(), identityID, 0, 5)
		require.NoError(t, err)
		require.Equal(t, 1, total)
		record := records[0]
		assert.Equal(t, auditlog.UserSearch, record.EventTypeID)
		assert.Equal(t, auditlog.EventParams{
			"query": "foo",
		}, record.EventParams)
	})

	s.T().Run("failures", func(t *testing.T) {

		t.Run("missing JWT", func(t *testing.T) {
			// when/then
			apptest.SearchUsersSearchUnauthorized(t, context.Background(), svc, ctrl, nil, nil, "foo")
		})

		t.Run("invalid JWT - invalid signing key", func(t *testing.T) {
			// given
			identityID := uuid.NewV4()
			// using PrivateKey2 but sibling PublicKey2 has not been loaded in TokenManager
			ctx, err := testjwt.NewJWTContext(identityID.String(), kid, testjwt.PrivateKey2)
			require.NoError(t, err)
			// when/then
			apptest.SearchUsersSearchUnauthorized(t, ctx, svc, ctrl, nil, nil, "foo")
		})

		t.Run("invalid JWT - invalid sub", func(t *testing.T) {
			// given
			identityID := "not.a.uuid"
			ctx, err := testjwt.NewJWTContext(identityID, kid, testjwt.PrivateKey1)
			require.NoError(t, err)
			// when/then
			apptest.SearchUsersSearchUnauthorized(t, ctx, svc, ctrl, nil, nil, "foo")
		})
	})

}

func (s *SearchControllerBlackboxTestSuite) TestSearchUsersFailures() {

}
