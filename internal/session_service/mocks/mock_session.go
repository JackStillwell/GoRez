// Code generated by MockGen. DO NOT EDIT.
// Source: session.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/JackStillwell/GoRez/internal/session_service/models"
	gomock "github.com/golang/mock/gomock"
)

// MockSessionService is a mock of SessionService interface.
type MockSessionService struct {
	ctrl     *gomock.Controller
	recorder *MockSessionServiceMockRecorder
}

// MockSessionServiceMockRecorder is the mock recorder for MockSessionService.
type MockSessionServiceMockRecorder struct {
	mock *MockSessionService
}

// NewMockSessionService creates a new mock instance.
func NewMockSessionService(ctrl *gomock.Controller) *MockSessionService {
	mock := &MockSessionService{ctrl: ctrl}
	mock.recorder = &MockSessionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionService) EXPECT() *MockSessionServiceMockRecorder {
	return m.recorder
}

// BadSession mocks base method.
func (m *MockSessionService) BadSession(arg0 []*models.Session) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BadSession", arg0)
}

// BadSession indicates an expected call of BadSession.
func (mr *MockSessionServiceMockRecorder) BadSession(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BadSession", reflect.TypeOf((*MockSessionService)(nil).BadSession), arg0)
}

// GetAvailableSessions mocks base method.
func (m *MockSessionService) GetAvailableSessions() []*models.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableSessions")
	ret0, _ := ret[0].([]*models.Session)
	return ret0
}

// GetAvailableSessions indicates an expected call of GetAvailableSessions.
func (mr *MockSessionServiceMockRecorder) GetAvailableSessions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableSessions", reflect.TypeOf((*MockSessionService)(nil).GetAvailableSessions))
}

// ReleaseSession mocks base method.
func (m *MockSessionService) ReleaseSession(arg0 []*models.Session) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReleaseSession", arg0)
}

// ReleaseSession indicates an expected call of ReleaseSession.
func (mr *MockSessionServiceMockRecorder) ReleaseSession(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseSession", reflect.TypeOf((*MockSessionService)(nil).ReleaseSession), arg0)
}

// ReserveSession mocks base method.
func (m *MockSessionService) ReserveSession(arg0 int, arg1 chan *models.Session) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReserveSession", arg0, arg1)
}

// ReserveSession indicates an expected call of ReserveSession.
func (mr *MockSessionServiceMockRecorder) ReserveSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReserveSession", reflect.TypeOf((*MockSessionService)(nil).ReserveSession), arg0, arg1)
}
