// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// GetRole mocks base method
func (m *MockStore) GetRole(name string) (*storage.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", name)
	ret0, _ := ret[0].(*storage.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole
func (mr *MockStoreMockRecorder) GetRole(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockStore)(nil).GetRole), name)
}

// GetAllRoles mocks base method
func (m *MockStore) GetAllRoles() ([]*storage.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRoles")
	ret0, _ := ret[0].([]*storage.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllRoles indicates an expected call of GetAllRoles
func (mr *MockStoreMockRecorder) GetAllRoles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRoles", reflect.TypeOf((*MockStore)(nil).GetAllRoles))
}

// AddRole mocks base method
func (m *MockStore) AddRole(arg0 *storage.Role) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRole", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRole indicates an expected call of AddRole
func (mr *MockStoreMockRecorder) AddRole(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRole", reflect.TypeOf((*MockStore)(nil).AddRole), arg0)
}

// UpdateRole mocks base method
func (m *MockStore) UpdateRole(arg0 *storage.Role) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRole", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRole indicates an expected call of UpdateRole
func (mr *MockStoreMockRecorder) UpdateRole(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRole", reflect.TypeOf((*MockStore)(nil).UpdateRole), arg0)
}

// RemoveRole mocks base method
func (m *MockStore) RemoveRole(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveRole", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveRole indicates an expected call of RemoveRole
func (mr *MockStoreMockRecorder) RemoveRole(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRole", reflect.TypeOf((*MockStore)(nil).RemoveRole), name)
}