// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/providers/gitlab_client.go

// Package providers is a generated GoMock package.
package providers

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gitlab "github.com/xanzy/go-gitlab"
)

// MockClientInterface is a mock of ClientInterface interface.
type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientInterfaceMockRecorder
}

// MockClientInterfaceMockRecorder is the mock recorder for MockClientInterface.
type MockClientInterfaceMockRecorder struct {
	mock *MockClientInterface
}

// NewMockClientInterface creates a new mock instance.
func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &MockClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientInterface) EXPECT() *MockClientInterfaceMockRecorder {
	return m.recorder
}

// ListAllGroups mocks base method.
func (m *MockClientInterface) ListAllGroups() ([]*gitlab.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllGroups")
	ret0, _ := ret[0].([]*gitlab.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllGroups indicates an expected call of ListAllGroups.
func (mr *MockClientInterfaceMockRecorder) ListAllGroups() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllGroups", reflect.TypeOf((*MockClientInterface)(nil).ListAllGroups))
}

// ListAllUsers mocks base method.
func (m *MockClientInterface) ListAllUsers() ([]*gitlab.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllUsers")
	ret0, _ := ret[0].([]*gitlab.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllUsers indicates an expected call of ListAllUsers.
func (mr *MockClientInterfaceMockRecorder) ListAllUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllUsers", reflect.TypeOf((*MockClientInterface)(nil).ListAllUsers))
}

// NewClient mocks base method.
func (m *MockClientInterface) NewClient(token, baseURL string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "NewClient", token, baseURL)
}

// NewClient indicates an expected call of NewClient.
func (mr *MockClientInterfaceMockRecorder) NewClient(token, baseURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClient", reflect.TypeOf((*MockClientInterface)(nil).NewClient), token, baseURL)
}
