// Code generated by MockGen. DO NOT EDIT.
// Source: alert_manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	alertmanager "github.com/stackrox/rox/central/detection/alertmanager"
	storage "github.com/stackrox/rox/generated/storage"
	set "github.com/stackrox/rox/pkg/set"
	reflect "reflect"
)

// MockAlertManager is a mock of AlertManager interface
type MockAlertManager struct {
	ctrl     *gomock.Controller
	recorder *MockAlertManagerMockRecorder
}

// MockAlertManagerMockRecorder is the mock recorder for MockAlertManager
type MockAlertManagerMockRecorder struct {
	mock *MockAlertManager
}

// NewMockAlertManager creates a new mock instance
func NewMockAlertManager(ctrl *gomock.Controller) *MockAlertManager {
	mock := &MockAlertManager{ctrl: ctrl}
	mock.recorder = &MockAlertManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAlertManager) EXPECT() *MockAlertManagerMockRecorder {
	return m.recorder
}

// AlertAndNotify mocks base method
func (m *MockAlertManager) AlertAndNotify(ctx context.Context, alerts []*storage.Alert, oldAlertFilters ...alertmanager.AlertFilterOption) (set.StringSet, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, alerts}
	for _, a := range oldAlertFilters {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AlertAndNotify", varargs...)
	ret0, _ := ret[0].(set.StringSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AlertAndNotify indicates an expected call of AlertAndNotify
func (mr *MockAlertManagerMockRecorder) AlertAndNotify(ctx, alerts interface{}, oldAlertFilters ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, alerts}, oldAlertFilters...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AlertAndNotify", reflect.TypeOf((*MockAlertManager)(nil).AlertAndNotify), varargs...)
}