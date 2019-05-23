package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/fabric8-services/admin-console/app"
	apptest "github.com/fabric8-services/admin-console/app/test"
	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/admin-console/auditlog"
	"github.com/fabric8-services/admin-console/configuration"
	"github.com/fabric8-services/admin-console/controller"
	"github.com/fabric8-services/fabric8-common/auth"
	"github.com/fabric8-services/fabric8-common/resource"
	testauth "github.com/fabric8-services/fabric8-common/test/auth"
	testsuite "github.com/fabric8-services/fabric8-common/test/suite"
	"github.com/goadesign/goa"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuditLogsControllerBlackboxTestSuite struct {
	testsuite.DBTestSuite
	app *application.GormApplication
}

func TestAuditLogs(t *testing.T) {
	resource.Require(t, resource.Database)
	config := configuration.New()
	suite.Run(t, &AuditLogsControllerBlackboxTestSuite{DBTestSuite: testsuite.NewDBTestSuite(config)})
}

func (s *AuditLogsControllerBlackboxTestSuite) SetupSuite() {
	s.DBTestSuite.SetupSuite()
	s.app = application.NewGormApplication(s.DB)
}

func (s *AuditLogsControllerBlackboxTestSuite) TestCreateAuditLog() {

	// given
	svc := goa.New("auditlogs")
	ctrl := controller.NewAuditLogsController(svc, s.app)
	ctx, err := testauth.EmbedServiceAccountTokenInContext(context.Background(), &testauth.Identity{
		Username: auth.Auth,
		ID:       uuid.NewV4(),
	})
	require.NoError(s.T(), err)
	s.Run("success", func() {
		// when
		eventParams := auditlog.EventParams{
			"notification_deactivation": time.Now().Format("2006-01-02:15:04:05"),
			"scheduled_deactivation":    time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02:15:04:05"),
		}
		apptest.CreateAuditLogNoContent(s.T(), ctx, svc, ctrl, "username", &app.CreateAuditLogPayload{
			Data: &app.CreateAuditLogData{
				Type: "audit_logs",
				Attributes: &app.CreateAuditLogDataAttributes{
					EventType:   auditlog.UserDeactivationEvent,
					EventParams: eventParams,
				},
			},
		})
		// then check that the data was collected
		records, total, err := auditlog.NewRepository(s.DB).ListByUsername(context.Background(), "username", 0, 5)
		require.NoError(s.T(), err)
		require.Equal(s.T(), 1, total)
		record := records[0]
		assert.Equal(s.T(), auditlog.UserDeactivation, record.EventTypeID)
		assert.Equal(s.T(), eventParams, record.EventParams)
	})

	s.Run("failures", func() {
		s.Run("bad request", func() {
			s.Run("invalid event type", func() {
				// when/then
				apptest.CreateAuditLogBadRequest(s.T(), ctx, svc, ctrl, "username", &app.CreateAuditLogPayload{
					Data: &app.CreateAuditLogData{
						Type: "audit_logs",
						Attributes: &app.CreateAuditLogDataAttributes{
							EventType:   "invalid",
							EventParams: map[string]interface{}{},
						},
					},
				})
			})
		})

		s.Run("unauthorized", func() {
			s.Run("missing token", func() {
				// when/then
				apptest.CreateAuditLogUnauthorized(s.T(), context.Background(), svc, ctrl, "username", &app.CreateAuditLogPayload{
					Data: &app.CreateAuditLogData{
						Type: "audit_logs",
						Attributes: &app.CreateAuditLogDataAttributes{
							EventType: "user_deactivation_notification",
							EventParams: map[string]interface{}{
								"notification_deactivation": time.Now().Format("2006-01-02:15:04:05"),
								"scheduled_deactivation":    time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02:15:04:05"),
							},
						},
					},
				})
			})

			s.Run("invalid token", func() {
				// given
				ctx, _, err := testauth.EmbedUserTokenInContext(context.Background(), testauth.NewIdentity())
				require.NoError(s.T(), err)
				// when/then
				apptest.CreateAuditLogUnauthorized(s.T(), ctx, svc, ctrl, "username", &app.CreateAuditLogPayload{
					Data: &app.CreateAuditLogData{
						Type: "audit_logs",
						Attributes: &app.CreateAuditLogDataAttributes{
							EventType: "user_deactivation_notification",
							EventParams: map[string]interface{}{
								"notification_deactivation": time.Now().Format("2006-01-02:15:04:05"),
								"scheduled_deactivation":    time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02:15:04:05"),
							},
						},
					},
				})
			})
		})
	})
}
