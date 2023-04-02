// Code generated by MockGen. DO NOT EDIT.
// Source: session.go

// Package mock is a generated GoMock package.
package mock

import (
	models "github.com/JackStillwell/GoRez/internal/session/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetAvailableSessions mocks base method
func (m *MockService) GetAvailableSessions() []*models.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableSessions")
	ret0, _ := ret[0].([]*models.Session)
	return ret0
}

// GetAvailableSessions indicates an expected call of GetAvailableSessions
func (mr *MockServiceMockRecorder) GetAvailableSessions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableSessions", reflect.TypeOf((*MockService)(nil).GetAvailableSessions))
}

// ReserveSession mocks base method
func (m *MockService) ReserveSession(arg0 int, arg1 chan *models.Session) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReserveSession", arg0, arg1)
}

// ReserveSession indicates an expected call of ReserveSession
func (mr *MockServiceMockRecorder) ReserveSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReserveSession", reflect.TypeOf((*MockService)(nil).ReserveSession), arg0, arg1)
}

// ReleaseSession mocks base method
func (m *MockService) ReleaseSession(arg0 []*models.Session) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReleaseSession", arg0)
}

// ReleaseSession indicates an expected call of ReleaseSession
func (mr *MockServiceMockRecorder) ReleaseSession(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseSession", reflect.TypeOf((*MockService)(nil).ReleaseSession), arg0)
}

// BadSession mocks base method
func (m *MockService) BadSession(arg0 []*models.Session) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BadSession", arg0)
}

// BadSession indicates an expected call of BadSession
func (mr *MockServiceMockRecorder) BadSession(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BadSession", reflect.TypeOf((*MockService)(nil).BadSession), arg0)
}