package configuration

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ManagerConfiguration" can be found in github.com/fabric8-services/admin-console/vendor/github.com/fabric8-services/fabric8-common/token
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//ManagerConfigurationMock implements github.com/fabric8-services/admin-console/vendor/github.com/fabric8-services/fabric8-common/token.ManagerConfiguration
type ManagerConfigurationMock struct {
	t minimock.Tester

	GetAuthServiceURLFunc       func() (r string)
	GetAuthServiceURLCounter    uint64
	GetAuthServiceURLPreCounter uint64
	GetAuthServiceURLMock       mManagerConfigurationMockGetAuthServiceURL

	GetDevModePrivateKeyFunc       func() (r []byte)
	GetDevModePrivateKeyCounter    uint64
	GetDevModePrivateKeyPreCounter uint64
	GetDevModePrivateKeyMock       mManagerConfigurationMockGetDevModePrivateKey
}

//NewManagerConfigurationMock returns a mock for github.com/fabric8-services/admin-console/vendor/github.com/fabric8-services/fabric8-common/token.ManagerConfiguration
func NewManagerConfigurationMock(t minimock.Tester) *ManagerConfigurationMock {
	m := &ManagerConfigurationMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetAuthServiceURLMock = mManagerConfigurationMockGetAuthServiceURL{mock: m}
	m.GetDevModePrivateKeyMock = mManagerConfigurationMockGetDevModePrivateKey{mock: m}

	return m
}

type mManagerConfigurationMockGetAuthServiceURL struct {
	mock *ManagerConfigurationMock
}

//Return sets up a mock for ManagerConfiguration.GetAuthServiceURL to return Return's arguments
func (m *mManagerConfigurationMockGetAuthServiceURL) Return(r string) *ManagerConfigurationMock {
	m.mock.GetAuthServiceURLFunc = func() string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ManagerConfiguration.GetAuthServiceURL method
func (m *mManagerConfigurationMockGetAuthServiceURL) Set(f func() (r string)) *ManagerConfigurationMock {
	m.mock.GetAuthServiceURLFunc = f

	return m.mock
}

//GetAuthServiceURL implements github.com/fabric8-services/admin-console/vendor/github.com/fabric8-services/fabric8-common/token.ManagerConfiguration interface
func (m *ManagerConfigurationMock) GetAuthServiceURL() (r string) {
	atomic.AddUint64(&m.GetAuthServiceURLPreCounter, 1)
	defer atomic.AddUint64(&m.GetAuthServiceURLCounter, 1)

	if m.GetAuthServiceURLFunc == nil {
		m.t.Fatal("Unexpected call to ManagerConfigurationMock.GetAuthServiceURL")
		return
	}

	return m.GetAuthServiceURLFunc()
}

//GetAuthServiceURLMinimockCounter returns a count of ManagerConfigurationMock.GetAuthServiceURLFunc invocations
func (m *ManagerConfigurationMock) GetAuthServiceURLMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAuthServiceURLCounter)
}

//GetAuthServiceURLMinimockPreCounter returns the value of ManagerConfigurationMock.GetAuthServiceURL invocations
func (m *ManagerConfigurationMock) GetAuthServiceURLMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAuthServiceURLPreCounter)
}

type mManagerConfigurationMockGetDevModePrivateKey struct {
	mock *ManagerConfigurationMock
}

//Return sets up a mock for ManagerConfiguration.GetDevModePrivateKey to return Return's arguments
func (m *mManagerConfigurationMockGetDevModePrivateKey) Return(r []byte) *ManagerConfigurationMock {
	m.mock.GetDevModePrivateKeyFunc = func() []byte {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ManagerConfiguration.GetDevModePrivateKey method
func (m *mManagerConfigurationMockGetDevModePrivateKey) Set(f func() (r []byte)) *ManagerConfigurationMock {
	m.mock.GetDevModePrivateKeyFunc = f

	return m.mock
}

//GetDevModePrivateKey implements github.com/fabric8-services/admin-console/vendor/github.com/fabric8-services/fabric8-common/token.ManagerConfiguration interface
func (m *ManagerConfigurationMock) GetDevModePrivateKey() (r []byte) {
	atomic.AddUint64(&m.GetDevModePrivateKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetDevModePrivateKeyCounter, 1)

	if m.GetDevModePrivateKeyFunc == nil {
		m.t.Fatal("Unexpected call to ManagerConfigurationMock.GetDevModePrivateKey")
		return
	}

	return m.GetDevModePrivateKeyFunc()
}

//GetDevModePrivateKeyMinimockCounter returns a count of ManagerConfigurationMock.GetDevModePrivateKeyFunc invocations
func (m *ManagerConfigurationMock) GetDevModePrivateKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDevModePrivateKeyCounter)
}

//GetDevModePrivateKeyMinimockPreCounter returns the value of ManagerConfigurationMock.GetDevModePrivateKey invocations
func (m *ManagerConfigurationMock) GetDevModePrivateKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDevModePrivateKeyPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ManagerConfigurationMock) ValidateCallCounters() {

	if m.GetAuthServiceURLFunc != nil && atomic.LoadUint64(&m.GetAuthServiceURLCounter) == 0 {
		m.t.Fatal("Expected call to ManagerConfigurationMock.GetAuthServiceURL")
	}

	if m.GetDevModePrivateKeyFunc != nil && atomic.LoadUint64(&m.GetDevModePrivateKeyCounter) == 0 {
		m.t.Fatal("Expected call to ManagerConfigurationMock.GetDevModePrivateKey")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ManagerConfigurationMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ManagerConfigurationMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ManagerConfigurationMock) MinimockFinish() {

	if m.GetAuthServiceURLFunc != nil && atomic.LoadUint64(&m.GetAuthServiceURLCounter) == 0 {
		m.t.Fatal("Expected call to ManagerConfigurationMock.GetAuthServiceURL")
	}

	if m.GetDevModePrivateKeyFunc != nil && atomic.LoadUint64(&m.GetDevModePrivateKeyCounter) == 0 {
		m.t.Fatal("Expected call to ManagerConfigurationMock.GetDevModePrivateKey")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ManagerConfigurationMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ManagerConfigurationMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetAuthServiceURLFunc == nil || atomic.LoadUint64(&m.GetAuthServiceURLCounter) > 0)
		ok = ok && (m.GetDevModePrivateKeyFunc == nil || atomic.LoadUint64(&m.GetDevModePrivateKeyCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetAuthServiceURLFunc != nil && atomic.LoadUint64(&m.GetAuthServiceURLCounter) == 0 {
				m.t.Error("Expected call to ManagerConfigurationMock.GetAuthServiceURL")
			}

			if m.GetDevModePrivateKeyFunc != nil && atomic.LoadUint64(&m.GetDevModePrivateKeyCounter) == 0 {
				m.t.Error("Expected call to ManagerConfigurationMock.GetDevModePrivateKey")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *ManagerConfigurationMock) AllMocksCalled() bool {

	if m.GetAuthServiceURLFunc != nil && atomic.LoadUint64(&m.GetAuthServiceURLCounter) == 0 {
		return false
	}

	if m.GetDevModePrivateKeyFunc != nil && atomic.LoadUint64(&m.GetDevModePrivateKeyCounter) == 0 {
		return false
	}

	return true
}
