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

type RecordRepositoryBlackBoxTest struct {
	testsuite.DBTestSuite
	repo auditlog.RecordRepository
}

func TestRecordRepository(t *testing.T) {
	resource.Require(t, resource.Database)
	config, err := configuration.New()
	require.NoError(t, err)
	suite.Run(t, &RecordRepositoryBlackBoxTest{DBTestSuite: testsuite.NewDBTestSuite(config)})
}

func (s *RecordRepositoryBlackBoxTest) SetupTest() {
	s.DBTestSuite.SetupTest()
	s.repo = auditlog.NewRecordRepository(s.DB)
}

func (s *RecordRepositoryBlackBoxTest) TestCreateRecord() {

	s.T().Run("ok", func(t *testing.T) {
		// given
		before := time.Now()
		record := auditlog.Record{
			EventTypeID: auditlog.UserSearch,
			IdentityID:  uuid.NewV4(),
			EventParams: auditlog.EventParams{},
		}
		// when
		err := s.repo.Create(context.Background(), &record)
		// then
		require.NoError(t, err)
		assert.NotEqual(t, uuid.NullUUID{}, record.ID)
		assert.True(t, record.CreatedAt.After(before)) // "is after before". hahahaha....
	})

	s.T().Run("failure", func(t *testing.T) {

		t.Run("missing event type", func(t *testing.T) {
			// given
			record := auditlog.Record{
				IdentityID:  uuid.NewV4(),
				EventParams: auditlog.EventParams{},
			}
			// when
			err := s.repo.Create(context.Background(), &record)
			// then
			require.Error(t, err)
		})

		t.Run("missing identity id", func(t *testing.T) {
			// given
			record := auditlog.Record{
				EventTypeID: auditlog.UserSearch,
				EventParams: auditlog.EventParams{},
			}
			// when
			err := s.repo.Create(context.Background(), &record)
			// then
			require.Error(t, err)
		})
	})
}

func (s *RecordRepositoryBlackBoxTest) TestLoadByID() {

	s.T().Run("ok", func(t *testing.T) {
		// given
		record := auditlog.Record{
			EventTypeID: auditlog.UserSearch,
			IdentityID:  uuid.NewV4(),
			EventParams: auditlog.EventParams{
				"idx": 1,
			},
		}
		err := s.repo.Create(context.Background(), &record)
		require.NoError(t, err)
		// when
		result, err := s.repo.LoadByID(context.Background(), record.ID)
		// then
		require.NoError(t, err)
		// comparing 'CreatedAt' may cause troubles b/c of nanosecond roundings, so let's just verify that the result ID is the one expected
		assert.Equal(t, record.ID, result.ID)

	})
	s.T().Run("not found", func(t *testing.T) {
		// when
		_, err := s.repo.LoadByID(context.Background(), uuid.NewV4())
		// then
		require.Error(t, err)
		assert.IsType(t, errors.NotFoundError{}, err)
	})
}

func (s *RecordRepositoryBlackBoxTest) TestList() {
	// given 2 users with 12 records each
	identity1 := uuid.NewV4()
	identity2 := uuid.NewV4()
	for _, identity := range []uuid.UUID{identity1, identity2} {
		for i := 0; i < 12; i++ {
			record := auditlog.Record{
				EventTypeID: auditlog.UserSearch,
				IdentityID:  identity,
				EventParams: auditlog.EventParams{
					"idx": i,
				},
			}
			err := s.repo.Create(context.Background(), &record)
			require.NoError(s.T(), err)
		}
	}

	s.T().Run("ok", func(t *testing.T) {

		t.Run("1st page of 5", func(t *testing.T) {
			// when
			records, count, err := s.repo.List(context.Background(), identity1, 0, 5)
			// then
			require.NoError(t, err)
			assert.Equal(t, 12, count)
			require.Len(t, records, 5) // full page
			for idx, record := range records {
				assert.Equal(t, identity1, record.IdentityID)
				require.NotNil(t, record.EventParams["idx"])
				assert.Equal(t, float64(idx), record.EventParams["idx"])
			}
		})

		t.Run("2nd page of 5", func(t *testing.T) {
			// when
			records, count, err := s.repo.List(context.Background(), identity1, 5, 5)
			// then
			require.NoError(t, err)
			assert.Equal(t, 12, count)
			require.Len(t, records, 5) // full page
			for idx, record := range records {
				assert.Equal(t, identity1, record.IdentityID)
				require.NotNil(t, record.EventParams["idx"])
				assert.Equal(t, float64(idx+5), record.EventParams["idx"])
			}
		})

		t.Run("last page of 2", func(t *testing.T) {
			// when
			records, count, err := s.repo.List(context.Background(), identity1, 10, 5)
			// then
			require.NoError(t, err)
			assert.Equal(t, 12, count)
			require.Len(t, records, 2) // last records, not a full page
			for idx, record := range records {
				assert.Equal(t, identity1, record.IdentityID)
				require.NotNil(t, record.EventParams["idx"])
				assert.Equal(t, float64(idx+10), record.EventParams["idx"])
			}
		})

		t.Run("out of range", func(t *testing.T) {
			// when
			records, count, err := s.repo.List(context.Background(), identity1, 15, 5)
			// then
			require.NoError(t, err)
			assert.Equal(t, 12, count)
			assert.Len(t, records, 0)
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
