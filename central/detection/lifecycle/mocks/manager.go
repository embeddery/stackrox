// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/detection/lifecycle (interfaces: Manager)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	common "github.com/stackrox/rox/central/sensor/service/common"
	storage "github.com/stackrox/rox/generated/storage"
	reflect "reflect"
)

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// DeploymentRemoved mocks base method
func (m *MockManager) DeploymentRemoved(arg0 *storage.Deployment) error {
	ret := m.ctrl.Call(m, "DeploymentRemoved", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeploymentRemoved indicates an expected call of DeploymentRemoved
func (mr *MockManagerMockRecorder) DeploymentRemoved(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeploymentRemoved", reflect.TypeOf((*MockManager)(nil).DeploymentRemoved), arg0)
}

// DeploymentUpdated mocks base method
func (m *MockManager) DeploymentUpdated(arg0 *storage.Deployment) (string, string, storage.EnforcementAction, error) {
	ret := m.ctrl.Call(m, "DeploymentUpdated", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(storage.EnforcementAction)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// DeploymentUpdated indicates an expected call of DeploymentUpdated
func (mr *MockManagerMockRecorder) DeploymentUpdated(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeploymentUpdated", reflect.TypeOf((*MockManager)(nil).DeploymentUpdated), arg0)
}

// IndicatorAdded mocks base method
func (m *MockManager) IndicatorAdded(arg0 *storage.ProcessIndicator, arg1 common.MessageInjector) error {
	ret := m.ctrl.Call(m, "IndicatorAdded", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// IndicatorAdded indicates an expected call of IndicatorAdded
func (mr *MockManagerMockRecorder) IndicatorAdded(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndicatorAdded", reflect.TypeOf((*MockManager)(nil).IndicatorAdded), arg0, arg1)
}

// RecompilePolicy mocks base method
func (m *MockManager) RecompilePolicy(arg0 *storage.Policy) error {
	ret := m.ctrl.Call(m, "RecompilePolicy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecompilePolicy indicates an expected call of RecompilePolicy
func (mr *MockManagerMockRecorder) RecompilePolicy(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecompilePolicy", reflect.TypeOf((*MockManager)(nil).RecompilePolicy), arg0)
}

// RemovePolicy mocks base method
func (m *MockManager) RemovePolicy(arg0 string) error {
	ret := m.ctrl.Call(m, "RemovePolicy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePolicy indicates an expected call of RemovePolicy
func (mr *MockManagerMockRecorder) RemovePolicy(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePolicy", reflect.TypeOf((*MockManager)(nil).RemovePolicy), arg0)
}

// UpsertPolicy mocks base method
func (m *MockManager) UpsertPolicy(arg0 *storage.Policy) error {
	ret := m.ctrl.Call(m, "UpsertPolicy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertPolicy indicates an expected call of UpsertPolicy
func (mr *MockManagerMockRecorder) UpsertPolicy(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertPolicy", reflect.TypeOf((*MockManager)(nil).UpsertPolicy), arg0)
}
