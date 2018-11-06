package application_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/fabric8-services/admin-console/configuration"

	"github.com/fabric8-services/admin-console/application"
	"github.com/fabric8-services/fabric8-common/resource"
	testsuite "github.com/fabric8-services/fabric8-common/test/suite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransactionTestSuite struct {
	testsuite.DBTestSuite
	gormApplication *application.GormApplication
}

func TestRunTransaction(t *testing.T) {
	resource.Require(t, resource.Database)
	config, err := configuration.New()
	require.NoError(t, err)
	suite.Run(t, &TransactionTestSuite{DBTestSuite: testsuite.NewDBTestSuite(config)})
}

func (s *TransactionTestSuite) SetupTest() {
	s.DBTestSuite.SetupTest()
	s.gormApplication = application.NewGormApplication(s.DB)
}

func (s *TransactionTestSuite) TransactionTestSuiteInTime() {
	// given
	computeTime := 10 * time.Second
	// then
	err := application.Transactional(s.gormApplication, func(appl application.Application) error {
		time.Sleep(computeTime)
		return nil
	})
	// then
	require.NoError(s.T(), err)
}

func (s *TransactionTestSuite) TransactionTestSuiteOut() {
	// given
	computeTime := 6 * time.Minute
	application.SetDatabaseTransactionTimeout(5 * time.Second)
	// then
	err := application.Transactional(s.gormApplication, func(appl application.Application) error {
		time.Sleep(computeTime)
		return nil
	})
	// then
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "database transaction timeout")
}

func (s *TransactionTestSuite) TransactionTestSuitePanicAndRecoverWithStack() {
	// then
	err := application.Transactional(s.gormApplication, func(appl application.Application) error {
		bar := func(a, b interface{}) {
			// This comparison while legal at compile time will cause a runtime
			// error like this: "comparing uncomparable type
			// map[string]interface {}". The transaction will panic and recover
			// but you will probably never find out where the error came from if
			// the stack is not captured in the transaction recovery. This test
			// ensures that the stack is captured.
			if a == b {
				fmt.Printf("never executed")
			}
		}
		foo := func() {
			a := map[string]interface{}{}
			b := map[string]interface{}{}
			bar(a, b)
		}
		foo()
		return nil
	})
	// then
	require.Error(s.T(), err)
	// ensure there's a proper stack trace that contains the name of this test
	require.Contains(s.T(), err.Error(), "(*TransactionTestSuite).TransactionTestSuitePanicAndRecoverWithStack.func1(")
}
