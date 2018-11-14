package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fabric8-services/admin-console/controller"
	"github.com/goadesign/goa"

	apptest "github.com/fabric8-services/admin-console/app/test"
	"github.com/fabric8-services/admin-console/configuration"
	testcontroller "github.com/fabric8-services/admin-console/test/generated/controller"
	"github.com/fabric8-services/fabric8-common/resource"
	testsuite "github.com/fabric8-services/fabric8-common/test/suite"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestStatusController(t *testing.T) {
	resource.Require(t, resource.Database)
	config, err := configuration.New()
	require.NoError(t, err)
	suite.Run(t, &StatusControllerBlackboxTestSuite{DBTestSuite: testsuite.NewDBTestSuite(config)})
}

type StatusControllerBlackboxTestSuite struct {
	testsuite.DBTestSuite
}

func newStatusController(dbchecker controller.DBChecker, config controller.StatusControllerConfiguration) (*goa.Service, *controller.StatusController) {
	svc := goa.New("status")
	ctrl := controller.NewStatusController(svc, dbchecker, config)
	return svc, ctrl
}

func (s *StatusControllerBlackboxTestSuite) TestShowStatus() {

	dbChecker := testcontroller.NewDBCheckerMock(s.T())
	config := testcontroller.NewStatusControllerConfigurationMock(s.T())
	svc, ctrl := newStatusController(dbChecker, config)
	ctx := context.Background()

	s.T().Run("with DB available", func(t *testing.T) {

		dbChecker.PingFunc = func() error {
			return nil
		}

		t.Run("with dev mode enabled", func(t *testing.T) {
			// given
			config.IsDeveloperModeEnabledFunc = func() bool {
				return true
			}
			config.DefaultConfigurationErrorFunc = func() error {
				return errors.New("developer mode is enabled")
			}
			// when
			_, status := apptest.ShowStatusOK(t, ctx, svc, ctrl)
			// then
			require.NotNil(t, status)
			assert.Contains(t, status.ConfigurationStatus, "developer mode is enabled")
			require.NotNil(t, status.DevMode)
			assert.True(t, *status.DevMode)
		})

		t.Run("with dev mode disabled and no configuration error", func(t *testing.T) {
			// given
			config.IsDeveloperModeEnabledFunc = func() bool {
				return false
			}
			config.DefaultConfigurationErrorFunc = func() error {
				return nil
			}
			// when
			_, status := apptest.ShowStatusOK(t, ctx, svc, ctrl)
			// then
			require.NotNil(t, status)
			assert.Equal(t, "OK", status.ConfigurationStatus)
			assert.Nil(t, status.DevMode)
		})
	})

	s.T().Run("with DB not available", func(t *testing.T) {

		// given
		dbChecker.PingFunc = func() error {
			return errors.New("db unavailable")
		}
		config.IsDeveloperModeEnabledFunc = func() bool {
			return true
		}
		config.DefaultConfigurationErrorFunc = func() error {
			return nil
		}
		// when/then
		apptest.ShowStatusServiceUnavailable(t, ctx, svc, ctrl)
	})

}
