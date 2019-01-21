package controller_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	apptest "github.com/fabric8-services/admin-console/app/test"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/auditlog"
	"github.com/fabric8-services/admin-console/configuration"
	"github.com/fabric8-services/admin-console/controller"
	testconfig "github.com/fabric8-services/admin-console/test/generated/controller"
	"github.com/fabric8-services/fabric8-common/resource"
	testauth "github.com/fabric8-services/fabric8-common/test/auth"
	testsuite "github.com/fabric8-services/fabric8-common/test/suite"

	"github.com/goadesign/goa"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"
)

func newTenantUpdateController(config controller.TenantUpdateControllerConfiguration, db application.DB) (*goa.Service, *controller.TenantUpdateController) {
	svc := goa.New("search")
	ctrl := controller.NewTenantUpdateController(svc,
		config,
		db,
	)
	return svc, ctrl
}

type TenantUpdateControllerBlackboxTestSuite struct {
	testsuite.DBTestSuite
	app *application.GormApplication
}

func (s *TenantUpdateControllerBlackboxTestSuite) SetupSuite() {
	s.DBTestSuite.SetupSuite()
	s.app = application.NewGormApplication(s.DB)
}

func TestTenantUpdateController(t *testing.T) {
	resource.Require(t, resource.Database)
	config := configuration.New()
	suite.Run(t, &TenantUpdateControllerBlackboxTestSuite{DBTestSuite: testsuite.NewDBTestSuite(config)})
}

func (s *TenantUpdateControllerBlackboxTestSuite) TestShowTenantUpdate() {
	// given
	config := testconfig.NewTenantUpdateControllerConfigurationMock(s.T())
	config.GetTenantServiceURLFunc = func() string {
		return "http://test-tenant"
	}
	svc, ctrl := newTenantUpdateController(config, s.app)
	defer gock.OffAll()

	s.T().Run("ok", func(t *testing.T) {
		t.Run("all clusters all envs", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Get("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusOK).BodyString(`{"data":"whatever"}`)
			// when
			apptest.ShowTenantUpdateOK(t, ctx, svc, ctrl, &authzHeader, nil, nil)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.ShowTenantUpdate, auditlog.EventParams{})
		})
		t.Run("single clusters single env", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)

			cluster := "cluster1"
			envType := "stage"
			gock.New("http://test-tenant").
				Get("/api/update").
				MatchHeader("Authorization", authzHeader).
				MatchParam("cluster_url", cluster).
				MatchParam("env_type", envType).
				Reply(http.StatusOK).BodyString(`{"data":"whatever"}`)
			// when
			apptest.ShowTenantUpdateOK(t, ctx, svc, ctrl, &authzHeader, nil, nil)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.ShowTenantUpdate, auditlog.EventParams{
				"clusterURL": cluster,
				"envType":    envType,
			})
		})
	})

	s.T().Run("failures", func(t *testing.T) {

		t.Run("missing JWT", func(t *testing.T) {
			// given
			gock.New("http://test-tenant").
				Get("/api/update").
				Reply(http.StatusUnauthorized)
			ctx := context.Background() // context is missing a JWT
			// when/then
			apptest.ShowTenantUpdateUnauthorized(t, ctx, svc, ctrl, nil, nil, nil)
		})

		t.Run("unauthorized", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Get("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusUnauthorized)
			// when
			apptest.ShowTenantUpdateUnauthorized(t, ctx, svc, ctrl, &authzHeader, nil, nil)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.ShowTenantUpdate, auditlog.EventParams{})
		})

		t.Run("internal server error", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Get("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusInternalServerError)
			// when
			apptest.ShowTenantUpdateInternalServerError(t, ctx, svc, ctrl, &authzHeader, nil, nil)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.ShowTenantUpdate, auditlog.EventParams{})
		})
	})
}
func (s *TenantUpdateControllerBlackboxTestSuite) TestStartTenantUpdate() {
	// given
	config := testconfig.NewTenantUpdateControllerConfigurationMock(s.T())
	config.GetTenantServiceURLFunc = func() string {
		return "http://test-tenant"
	}
	svc, ctrl := newTenantUpdateController(config, s.app)
	defer gock.OffAll()

	s.T().Run("ok", func(t *testing.T) {

		t.Run("all clusters all envs", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Post("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusAccepted).BodyString(`{"data":"whatever"}`)
			// when
			apptest.StartTenantUpdateAccepted(t, ctx, svc, ctrl, nil, nil, &authzHeader)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.StartTenantUpdate, auditlog.EventParams{})
		})

		t.Run("single clusters single env", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			cluster := "cluster1"
			envType := "stage"
			gock.New("http://test-tenant").
				Post("/api/update").
				MatchHeader("Authorization", authzHeader).
				MatchParam("cluster_url", cluster).
				MatchParam("env_type", envType).
				Reply(http.StatusAccepted).BodyString(`{"data":"whatever"}`)
			// when
			apptest.StartTenantUpdateAccepted(t, ctx, svc, ctrl, &cluster, &envType, &authzHeader)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.StartTenantUpdate, auditlog.EventParams{
				"clusterURL": cluster,
				"envType":    envType,
			})
		})
	})

	s.T().Run("failures", func(t *testing.T) {

		t.Run("missing JWT", func(t *testing.T) {
			// given
			gock.New("http://test-tenant").
				Post("/api/update").
				Reply(http.StatusUnauthorized)
			ctx := context.Background() // context is missing a JWT
			// when/then
			apptest.StartTenantUpdateUnauthorized(t, ctx, svc, ctrl, nil, nil, nil)
		})

		t.Run("unauthorized", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Post("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusUnauthorized)
			// when
			apptest.StartTenantUpdateUnauthorized(t, ctx, svc, ctrl, nil, nil, &authzHeader)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.StartTenantUpdate, auditlog.EventParams{})
		})

		t.Run("conflict", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Post("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusConflict)
			// when
			apptest.StartTenantUpdateConflict(t, ctx, svc, ctrl, nil, nil, &authzHeader)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.StartTenantUpdate, auditlog.EventParams{})
		})
		t.Run("bad request", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Post("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusBadRequest)
			// when
			apptest.StartTenantUpdateBadRequest(t, ctx, svc, ctrl, nil, nil, &authzHeader)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.StartTenantUpdate, auditlog.EventParams{})
		})

		t.Run("internal server error", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Post("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusInternalServerError)
			// when
			apptest.StartTenantUpdateInternalServerError(t, ctx, svc, ctrl, nil, nil, &authzHeader)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.StartTenantUpdate, auditlog.EventParams{})
		})
	})
}
func (s *TenantUpdateControllerBlackboxTestSuite) TestStopTenantUpdate() {
	// given
	config := testconfig.NewTenantUpdateControllerConfigurationMock(s.T())
	config.GetTenantServiceURLFunc = func() string {
		return "http://test-tenant"
	}
	svc, ctrl := newTenantUpdateController(config, s.app)
	defer gock.OffAll()

	s.T().Run("accepted", func(t *testing.T) {
		// given
		ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
		require.NoError(t, err)
		tk := goajwt.ContextJWT(ctx)
		require.NotNil(t, tk)
		authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
		gock.New("http://test-tenant").
			Delete("/api/update").
			MatchHeader("Authorization", authzHeader).
			Reply(http.StatusAccepted).BodyString(`{"data":"whatever"}`)
		// when
		apptest.StopTenantUpdateAccepted(t, ctx, svc, ctrl, &authzHeader)
		// then check that an audit record was created
		assertAuditLog(t, s.DB, *identity, auditlog.StopTenantUpdate, auditlog.EventParams{})
	})

	s.T().Run("failures", func(t *testing.T) {

		t.Run("missing JWT", func(t *testing.T) {
			// given
			gock.New("http://test-tenant").
				Delete("/api/update").
				Reply(http.StatusUnauthorized)
			ctx := context.Background() // context is missing a JWT
			// when/then
			apptest.ShowTenantUpdateUnauthorized(t, ctx, svc, ctrl, nil)
		})

		t.Run("unauthorized", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Delete("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusUnauthorized)
			// when
			apptest.StopTenantUpdateUnauthorized(t, ctx, svc, ctrl, &authzHeader)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.StopTenantUpdate, auditlog.EventParams{})
		})

		t.Run("internal server error", func(t *testing.T) {
			// given
			ctx, identity, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
			require.NoError(t, err)
			tk := goajwt.ContextJWT(ctx)
			require.NotNil(t, tk)
			authzHeader := fmt.Sprintf("Bearer %s", tk.Raw)
			gock.New("http://test-tenant").
				Delete("/api/update").
				MatchHeader("Authorization", authzHeader).
				Reply(http.StatusInternalServerError)
			// when
			apptest.StopTenantUpdateInternalServerError(t, ctx, svc, ctrl, &authzHeader)
			// then check that an audit record was created
			assertAuditLog(t, s.DB, *identity, auditlog.StopTenantUpdate, auditlog.EventParams{})
		})
	})
}

func assertAuditLog(t *testing.T, db *gorm.DB, identity testauth.Identity, expectedEventType uuid.UUID, expectedQueryParams auditlog.EventParams) {
	recordRepo := auditlog.NewRepository(db)
	records, total, err := recordRepo.ListByIdentityID(context.Background(), identity.ID, 0, 5)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	record := records[0]
	assert.Equal(t, expectedEventType, record.EventTypeID)
	assert.Equal(t, expectedQueryParams, record.EventParams)
}
