// Code generated by MockGen. DO NOT EDIT.
// Source: ./interfaces.go

// Package mock_gorezinternal is a generated GoMock package.
package mock_gorezinternal

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockHTTPGetter is a mock of HTTPGetter interface
type MockHTTPGetter struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPGetterMockRecorder
}

// MockHTTPGetterMockRecorder is the mock recorder for MockHTTPGetter
type MockHTTPGetterMockRecorder struct {
	mock *MockHTTPGetter
}

// NewMockHTTPGetter creates a new mock instance
func NewMockHTTPGetter(ctrl *gomock.Controller) *MockHTTPGetter {
	mock := &MockHTTPGetter{ctrl: ctrl}
	mock.recorder = &MockHTTPGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHTTPGetter) EXPECT() *MockHTTPGetterMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockHTTPGetter) Get(url string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", url)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockHTTPGetterMockRecorder) Get(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockHTTPGetter)(nil).Get), url)
}
