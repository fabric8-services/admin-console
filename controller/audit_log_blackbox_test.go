package controller_test

import (
	"context"
	"encoding/json"
	"fmt"
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
	app    *application.GormApplication
	config *configuration.Configuration
}

func TestAuditLogs(t *testing.T) {
	resource.Require(t, resource.Database)
	config := configuration.New()
	suite.Run(t, &AuditLogsControllerBlackboxTestSuite{
		DBTestSuite: testsuite.NewDBTestSuite(config),
		config:      config,
	})
}

func (s *AuditLogsControllerBlackboxTestSuite) SetupSuite() {
	s.DBTestSuite.SetupSuite()
	s.app = application.NewGormApplication(s.DB)
}

func (s *AuditLogsControllerBlackboxTestSuite) TestCreateAuditLog() {

	// given
	svc := goa.New("auditlogs")
	ctrl := controller.NewAuditLogsController(svc, s.config, s.app)
	ctx, err := testauth.EmbedServiceAccountTokenInContext(context.Background(), &testauth.Identity{
		Username: auth.Auth,
		ID:       uuid.NewV4(),
	})
	require.NoError(s.T(), err)
	s.Run("success", func() {

		s.Run("with event params", func() {
			// when
			username := fmt.Sprintf("user-%v", uuid.NewV4())
			eventParams := auditlog.EventParams{
				"notification_deactivation": time.Now().Format("2006-01-02:15:04:05"),
				"scheduled_deactivation":    time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02:15:04:05"),
			}
			apptest.CreateAuditLogNoContent(s.T(), ctx, svc, ctrl, username, &app.CreateAuditLogPayload{
				Data: &app.CreateAuditLogData{
					Type: "audit_logs",
					Attributes: &app.CreateAuditLogDataAttributes{
						EventType:   auditlog.UserDeactivationEvent,
						EventParams: eventParams,
					},
				},
			})
			// then check that the data was collected
			records, total, err := auditlog.NewRepository(s.DB).ListByUsername(context.Background(), username, 0, 5)
			require.NoError(s.T(), err)
			require.Equal(s.T(), 1, total)
			record := records[0]
			assert.Equal(s.T(), auditlog.UserDeactivation, record.EventTypeID)
			assert.Equal(s.T(), eventParams, record.EventParams)
		})

		s.Run("without event params", func() {
			// when
			username := fmt.Sprintf("user-%v", uuid.NewV4())
			apptest.CreateAuditLogNoContent(s.T(), ctx, svc, ctrl, username, &app.CreateAuditLogPayload{
				Data: &app.CreateAuditLogData{
					Type: "audit_logs",
					Attributes: &app.CreateAuditLogDataAttributes{
						EventType: auditlog.UserDeactivationEvent,
					},
				},
			})
			// then check that the data was collected
			records, total, err := auditlog.NewRepository(s.DB).ListByUsername(context.Background(), username, 0, 5)
			require.NoError(s.T(), err)
			require.Equal(s.T(), 1, total)
			record := records[0]
			assert.Equal(s.T(), auditlog.UserDeactivation, record.EventTypeID)
			assert.Nil(s.T(), record.EventParams)
		})
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

func (s *AuditLogsControllerBlackboxTestSuite) TestListAuditLogs() {

	// given
	svc := goa.New("auditlogs")
	ctrl := controller.NewAuditLogsController(svc, s.config, s.app)

	s.Run("success", func() {

		var ctx context.Context
		var requestingUser, targetUser string

		s.SetupSubtest = func() {
			// target user logs
			targetUser = fmt.Sprintf("user-foo-%v", uuid.NewV4())
			r := auditlog.NewRepository(s.DB)
			err := r.Create(context.Background(), &auditlog.AuditLog{
				Username:    targetUser,
				EventTypeID: auditlog.UserDeactivationNotification,
			})
			require.NoError(s.T(), err)
			err = r.Create(context.Background(), &auditlog.AuditLog{
				Username:    targetUser,
				EventTypeID: auditlog.UserDeactivation,
			})
			require.NoError(s.T(), err)
			// requesting user context with token
			requestingUser = fmt.Sprintf("requesting_user-%v", uuid.NewV4())
			ctx, _, err = testauth.EmbedTokenInContext("identity", requestingUser, testauth.WithEmailClaim("user@redhat.com"), testauth.WithEmailVerifiedClaim(true))
			require.NoError(s.T(), err)
		}
		s.TearDownSubtest = func() {
			s.CleanTest()
		}

		s.Run("first page of results", func() {
			// when
			_, result := apptest.ListForUserAuditLogOK(s.T(), ctx, svc, ctrl, targetUser, 0, 1)
			// then
			require.NotNil(s.T(), result)
			json, _ := json.MarshalIndent(result, "", "  ")
			s.T().Logf("result:\n%v\n", string(json))
			// verify data
			require.NotNil(s.T(), result.Data)
			require.Len(s.T(), result.Data, 1)
			assert.Equal(s.T(), auditlog.UserDeactivationNotificationEvent, result.Data[0].Attributes.EventType)
			assert.Nil(s.T(), result.Data[0].Attributes.EventParams)
			// verify links
			require.NotNil(s.T(), result.Links.First)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=0&page[size]=1", targetUser), *result.Links.First)
			require.Nil(s.T(), result.Links.Prev)
			require.NotNil(s.T(), result.Links.Next)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=1&page[size]=1", targetUser), *result.Links.Next)
			require.NotNil(s.T(), result.Links.Last)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=1&page[size]=1", targetUser), *result.Links.Last)
			// verify meta
			require.NotNil(s.T(), result.Meta)
			assert.Equal(s.T(), result.Meta.TotalCount, 2)
			// also, verify that an event was logged on behalf of the requesting user
			s.assertRequesterLogs(requestingUser, targetUser)
		})

		s.Run("last page of results", func() {
			// when
			_, result := apptest.ListForUserAuditLogOK(s.T(), ctx, svc, ctrl, targetUser, 1, 1)
			// then
			require.NotNil(s.T(), result)
			json, _ := json.MarshalIndent(result, "", "  ")
			s.T().Logf("result:\n%v\n", string(json))
			require.NotNil(s.T(), result.Data)
			require.Len(s.T(), result.Data, 1)
			require.NotNil(s.T(), result.Links)
			require.NotNil(s.T(), result.Meta)
			assert.Equal(s.T(), result.Meta.TotalCount, 2)
			// verify links
			require.NotNil(s.T(), result.Links.First)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=0&page[size]=1", targetUser), *result.Links.First)
			require.NotNil(s.T(), result.Links.Prev)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=0&page[size]=1", targetUser), *result.Links.Prev)
			require.Nil(s.T(), result.Links.Next)
			require.NotNil(s.T(), result.Links.Last)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=1&page[size]=1", targetUser), *result.Links.Last)
			// also, verify that an event was logged on behalf of the requesting user
			s.assertRequesterLogs(requestingUser, targetUser)
		})

		s.Run("all results", func() {
			// when
			_, result := apptest.ListForUserAuditLogOK(s.T(), ctx, svc, ctrl, targetUser, 0, 10)
			// then
			require.NotNil(s.T(), result)
			json, _ := json.MarshalIndent(result, "", "  ")
			s.T().Logf("result:\n%v\n", string(json))
			require.NotNil(s.T(), result.Data)
			require.Len(s.T(), result.Data, 2)
			require.NotNil(s.T(), result.Links)
			require.NotNil(s.T(), result.Meta)
			assert.Equal(s.T(), result.Meta.TotalCount, 2)
			// verify links
			require.NotNil(s.T(), result.Links.First)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=0&page[size]=10", targetUser), *result.Links.First)
			require.Nil(s.T(), result.Links.Prev)
			require.Nil(s.T(), result.Links.Next)
			require.NotNil(s.T(), result.Links.Last)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=0&page[size]=10", targetUser), *result.Links.Last)
			// also, verify that an event was logged on behalf of the requesting user
			s.assertRequesterLogs(requestingUser, targetUser)
		})

		s.Run("out of range", func() {
			// when
			_, result := apptest.ListForUserAuditLogOK(s.T(), ctx, svc, ctrl, targetUser, 100, 100)
			// then
			require.NotNil(s.T(), result)
			json, _ := json.MarshalIndent(result, "", "  ")
			s.T().Logf("result:\n%v\n", string(json))
			require.NotNil(s.T(), result.Data)
			require.Empty(s.T(), result.Data)
			require.NotNil(s.T(), result.Meta)
			assert.Equal(s.T(), result.Meta.TotalCount, 2)
			// verify links
			require.NotNil(s.T(), result.Links.First)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=0&page[size]=10", targetUser), *result.Links.First)
			require.Nil(s.T(), result.Links.Prev)
			require.Nil(s.T(), result.Links.Next)
			require.NotNil(s.T(), result.Links.Last)
			assert.Equal(s.T(), fmt.Sprintf("http:///api/auditlogs/users/%s?page[start]=0&page[size]=10", targetUser), *result.Links.Last)
			// also, verify that an event was logged on behalf of the requesting user
			s.assertRequesterLogs(requestingUser, targetUser)
		})

		s.Run("user has no audit log", func() {
			// when
			_, result := apptest.ListForUserAuditLogOK(s.T(), ctx, svc, ctrl, "user-bar", 1, 100)
			// then
			require.NotNil(s.T(), result.Data)
			require.Empty(s.T(), result.Data)
			require.NotNil(s.T(), result.Meta)
			assert.Equal(s.T(), result.Meta.TotalCount, 0)
			// also, verify that an event was logged on behalf of the requesting user
			s.assertRequesterLogs(requestingUser, "user-bar")
		})

	})

	s.Run("failure", func() {

		requestingUser := "user-foo"
		targetUser := "user-bar"
		s.Run("unauthorized - missing token", func() {
			// given
			ctx := context.Background()
			// when/then
			apptest.ListForUserAuditLogUnauthorized(s.T(), ctx, svc, ctrl, targetUser, 1, 0)
		})

		s.Run("forbidden - external user", func() {
			// given
			ctx, _, err := testauth.EmbedTokenInContext("identity", requestingUser, testauth.WithEmailClaim("user@foo.com"), testauth.WithEmailVerifiedClaim(true))
			require.NoError(s.T(), err)
			// when/then
			apptest.ListForUserAuditLogForbidden(s.T(), ctx, svc, ctrl, targetUser, 1, 0)
		})

		s.Run("forbidden - internal user with email not verified", func() {
			// given
			ctx, _, err := testauth.EmbedTokenInContext("identity", requestingUser, testauth.WithEmailClaim("user@redhat.com"), testauth.WithEmailVerifiedClaim(false))
			require.NoError(s.T(), err)
			// when/then
			apptest.ListForUserAuditLogForbidden(s.T(), ctx, svc, ctrl, targetUser, 1, 0)
		})

	})
}

func (s *AuditLogsControllerBlackboxTestSuite) assertRequesterLogs(requestingUser, eventUser string) {
	r := auditlog.NewRepository(s.DB)
	logs, total, err := r.ListByUsername(context.Background(), requestingUser, 0, 100)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, total)
	require.Len(s.T(), logs, 1)
	assert.Equal(s.T(), auditlog.ListAuditLogs, logs[0].EventTypeID)
	assert.Equal(s.T(), eventUser, logs[0].EventParams["user"])
}
