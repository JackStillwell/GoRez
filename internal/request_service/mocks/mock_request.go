// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces/request.go

// Package mock is a generated GoMock package.
package mock

import (
	models "github.com/JackStillwell/GoRez/internal/request_service/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	http "net/http"
	reflect "reflect"
)

// MockHTTPGet is a mock of HTTPGet interface
type MockHTTPGet struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPGetMockRecorder
}

// MockHTTPGetMockRecorder is the mock recorder for MockHTTPGet
type MockHTTPGetMockRecorder struct {
	mock *MockHTTPGet
}

// NewMockHTTPGet creates a new mock instance
func NewMockHTTPGet(ctrl *gomock.Controller) *MockHTTPGet {
	mock := &MockHTTPGet{ctrl: ctrl}
	mock.recorder = &MockHTTPGetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHTTPGet) EXPECT() *MockHTTPGetMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockHTTPGet) Get(url string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", url)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockHTTPGetMockRecorder) Get(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockHTTPGet)(nil).Get), url)
}

// MockRequestService is a mock of RequestService interface
type MockRequestService struct {
	ctrl     *gomock.Controller
	recorder *MockRequestServiceMockRecorder
}

// MockRequestServiceMockRecorder is the mock recorder for MockRequestService
type MockRequestServiceMockRecorder struct {
	mock *MockRequestService
}

// NewMockRequestService creates a new mock instance
func NewMockRequestService(ctrl *gomock.Controller) *MockRequestService {
	mock := &MockRequestService{ctrl: ctrl}
	mock.recorder = &MockRequestServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequestService) EXPECT() *MockRequestServiceMockRecorder {
	return m.recorder
}

// MakeRequest mocks base method
func (m *MockRequestService) MakeRequest(arg0 *models.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MakeRequest", arg0)
}

// MakeRequest indicates an expected call of MakeRequest
func (mr *MockRequestServiceMockRecorder) MakeRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeRequest", reflect.TypeOf((*MockRequestService)(nil).MakeRequest), arg0)
}

// GetResponse mocks base method
func (m *MockRequestService) GetResponse(arg0 *uuid.UUID, arg1 chan *models.RequestResponse) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetResponse", arg0, arg1)
}

// GetResponse indicates an expected call of GetResponse
func (mr *MockRequestServiceMockRecorder) GetResponse(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResponse", reflect.TypeOf((*MockRequestService)(nil).GetResponse), arg0, arg1)
}

// Close mocks base method
func (m *MockRequestService) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockRequestServiceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRequestService)(nil).Close))
}

// Request mocks base method
func (m *MockRequestService) Request(arg0 *models.Request) *models.RequestResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Request", arg0)
	ret0, _ := ret[0].(*models.RequestResponse)
	return ret0
}

// Request indicates an expected call of Request
func (mr *MockRequestServiceMockRecorder) Request(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Request", reflect.TypeOf((*MockRequestService)(nil).Request), arg0)
}

// MockRequestManager is a mock of RequestManager interface
type MockRequestManager struct {
	ctrl     *gomock.Controller
	recorder *MockRequestManagerMockRecorder
}

// MockRequestManagerMockRecorder is the mock recorder for MockRequestManager
type MockRequestManagerMockRecorder struct {
	mock *MockRequestManager
}

// NewMockRequestManager creates a new mock instance
func NewMockRequestManager(ctrl *gomock.Controller) *MockRequestManager {
	mock := &MockRequestManager{ctrl: ctrl}
	mock.recorder = &MockRequestManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequestManager) EXPECT() *MockRequestManagerMockRecorder {
	return m.recorder
}

// MakeRequest mocks base method
func (m *MockRequestManager) MakeRequest(arg0 *models.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MakeRequest", arg0)
}

// MakeRequest indicates an expected call of MakeRequest
func (mr *MockRequestManagerMockRecorder) MakeRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeRequest", reflect.TypeOf((*MockRequestManager)(nil).MakeRequest), arg0)
}

// GetResponse mocks base method
func (m *MockRequestManager) GetResponse(arg0 *uuid.UUID, arg1 chan *models.RequestResponse) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetResponse", arg0, arg1)
}

// GetResponse indicates an expected call of GetResponse
func (mr *MockRequestManagerMockRecorder) GetResponse(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResponse", reflect.TypeOf((*MockRequestManager)(nil).GetResponse), arg0, arg1)
}

// Close mocks base method
func (m *MockRequestManager) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockRequestManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRequestManager)(nil).Close))
}

// MockRequester is a mock of Requester interface
type MockRequester struct {
	ctrl     *gomock.Controller
	recorder *MockRequesterMockRecorder
}

// MockRequesterMockRecorder is the mock recorder for MockRequester
type MockRequesterMockRecorder struct {
	mock *MockRequester
}

// NewMockRequester creates a new mock instance
func NewMockRequester(ctrl *gomock.Controller) *MockRequester {
	mock := &MockRequester{ctrl: ctrl}
	mock.recorder = &MockRequesterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequester) EXPECT() *MockRequesterMockRecorder {
	return m.recorder
}

// Request mocks base method
func (m *MockRequester) Request(arg0 *models.Request) *models.RequestResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Request", arg0)
	ret0, _ := ret[0].(*models.RequestResponse)
	return ret0
}

// Request indicates an expected call of Request
func (mr *MockRequesterMockRecorder) Request(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Request", reflect.TypeOf((*MockRequester)(nil).Request), arg0)
}
