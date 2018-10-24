package auditlog_test

import (
	"context"
	"testing"
	"time"

	"github.com/fabric8-services/fabric8-common/errors"

	"github.com/fabric8-services/admin-console/auditlog"
	"github.com/fabric8-services/admin-console/configuration"
	"github.com/fabric8-services/fabric8-common/resource"
	testsuite "github.com/fabric8-services/fabric8-common/test/suite"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepositoryBlackBoxTest struct {
	testsuite.DBTestSuite
	repo auditlog.Repository
}

func TestRecordRepository(t *testing.T) {
	resource.Require(t, resource.Database)
	config, err := configuration.New()
	require.NoError(t, err)
	suite.Run(t, &RepositoryBlackBoxTest{DBTestSuite: testsuite.NewDBTestSuite(config)})
}

func (s *RepositoryBlackBoxTest) SetupTest() {
	s.DBTestSuite.SetupTest()
	s.repo = auditlog.NewRepository(s.DB)
}

func (s *RepositoryBlackBoxTest) TestCreateRecord() {

	s.T().Run("ok", func(t *testing.T) {
		// given
		before := time.Now()
		auditLog := auditlog.AuditLog{
			EventTypeID: auditlog.UserSearch,
			IdentityID:  uuid.NewV4(),
			EventParams: auditlog.EventParams{},
		}
		// when
		err := s.repo.Create(context.Background(), &auditLog)
		// then
		require.NoError(t, err)
		assert.NotEqual(t, uuid.NullUUID{}, auditLog.ID)
		assert.True(t, auditLog.CreatedAt.After(before)) // "is after before". hahahaha....
	})

	s.T().Run("failure", func(t *testing.T) {

		t.Run("missing event type", func(t *testing.T) {
			// given
			auditLog := auditlog.AuditLog{
				IdentityID:  uuid.NewV4(),
				EventParams: auditlog.EventParams{},
			}
			// when
			err := s.repo.Create(context.Background(), &auditLog)
			// then
			require.Error(t, err)
			assert.IsType(t, errors.BadParameterError{}, err)
			assert.Contains(t, err.Error(), "identity_id")
		})

		t.Run("missing identity id", func(t *testing.T) {
			// given
			auditLog := auditlog.AuditLog{
				EventTypeID: auditlog.UserSearch,
				EventParams: auditlog.EventParams{},
			}
			// when
			err := s.repo.Create(context.Background(), &auditLog)
			// then
			require.Error(t, err)
			assert.IsType(t, errors.BadParameterError{}, err)
			assert.Contains(t, err.Error(), "identity_id")
		})
	})
}

func (s *RepositoryBlackBoxTest) TestLoadByID() {

	s.T().Run("ok", func(t *testing.T) {
		// given
		auditLog := auditlog.AuditLog{
			EventTypeID: auditlog.UserSearch,
			IdentityID:  uuid.NewV4(),
			EventParams: auditlog.EventParams{
				"idx": 1,
			},
		}
		err := s.repo.Create(context.Background(), &auditLog)
		require.NoError(t, err)
		// when
		result, err := s.repo.LoadByID(context.Background(), auditLog.ID)
		// then
		require.NoError(t, err)
		// comparing 'CreatedAt' may cause troubles b/c of nanosecond roundings, so let's just verify that the result ID is the one expected
		assert.Equal(t, auditLog.ID, result.ID)

	})
	s.T().Run("not found", func(t *testing.T) {
		// when
		_, err := s.repo.LoadByID(context.Background(), uuid.NewV4())
		// then
		require.Error(t, err)
		assert.IsType(t, errors.NotFoundError{}, err)
	})
}

func (s *RepositoryBlackBoxTest) TestList() {
	// given 2 users with 12 auditLogs each
	identity1 := uuid.NewV4()
	identity2 := uuid.NewV4()
	for _, identity := range []uuid.UUID{identity1, identity2} {
		for i := 0; i < 12; i++ {
			auditLog := auditlog.AuditLog{
				EventTypeID: auditlog.UserSearch,
				IdentityID:  identity,
				EventParams: auditlog.EventParams{
					"idx": i,
				},
			}
			err := s.repo.Create(context.Background(), &auditLog)
			require.NoError(s.T(), err)
		}
	}

	s.T().Run("ok", func(t *testing.T) {

		t.Run("1st page of 5", func(t *testing.T) {
			// when
			auditLogs, count, err := s.repo.List(context.Background(), identity1, 0, 5)
			// then
			require.NoError(t, err)
			assert.Equal(t, 12, count)
			require.Len(t, auditLogs, 5) // full page
			for idx, auditLog := range auditLogs {
				assert.Equal(t, identity1, auditLog.IdentityID)
				require.NotNil(t, auditLog.EventParams["idx"])
				assert.Equal(t, float64(idx), auditLog.EventParams["idx"])
			}
		})

		t.Run("2nd page of 5", func(t *testing.T) {
			// when
			auditLogs, count, err := s.repo.List(context.Background(), identity1, 5, 5)
			// then
			require.NoError(t, err)
			assert.Equal(t, 12, count)
			require.Len(t, auditLogs, 5) // full page
			for idx, auditLog := range auditLogs {
				assert.Equal(t, identity1, auditLog.IdentityID)
				require.NotNil(t, auditLog.EventParams["idx"])
				assert.Equal(t, float64(idx+5), auditLog.EventParams["idx"])
			}
		})

		t.Run("last page of 2", func(t *testing.T) {
			// when
			auditLogs, count, err := s.repo.List(context.Background(), identity1, 10, 5)
			// then
			require.NoError(t, err)
			assert.Equal(t, 12, count)
			require.Len(t, auditLogs, 2) // last auditLogs, not a full page
			for idx, auditLog := range auditLogs {
				assert.Equal(t, identity1, auditLog.IdentityID)
				require.NotNil(t, auditLog.EventParams["idx"])
				assert.Equal(t, float64(idx+10), auditLog.EventParams["idx"])
			}
		})

		t.Run("out of range", func(t *testing.T) {
			// when
			auditLogs, count, err := s.repo.List(context.Background(), identity1, 15, 5)
			// then
			require.NoError(t, err)
			assert.Equal(t, 12, count)
			assert.Len(t, auditLogs, 0)
		})
	})

	s.T().Run("failures", func(t *testing.T) {

		t.Run("invalid start", func(t *testing.T) {
			// when
			_, _, err := s.repo.List(context.Background(), identity1, -1, 5)
			// then
			require.Error(t, err)
			require.IsType(t, errors.BadParameterError{}, err)
		})

		t.Run("invalid limit", func(t *testing.T) {
			// when
			_, _, err := s.repo.List(context.Background(), identity1, 0, -5)
			// then
			require.Error(t, err)
			require.IsType(t, errors.BadParameterError{}, err)
		})
	})

}
