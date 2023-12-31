// Code generated by MockGen. DO NOT EDIT.
// Source: route256/loms/internal/app/services (interfaces: OrdersProvider)
//
// Generated by this command:
//
//	mockgen -destination=internal/app/services/orders_provider_mock.go -package=services route256/loms/internal/app/services OrdersProvider
//
// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"
	models "route256/loms/internal/app/models"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockOrdersProvider is a mock of OrdersProvider interface.
type MockOrdersProvider struct {
	ctrl     *gomock.Controller
	recorder *MockOrdersProviderMockRecorder
}

// MockOrdersProviderMockRecorder is the mock recorder for MockOrdersProvider.
type MockOrdersProviderMockRecorder struct {
	mock *MockOrdersProvider
}

// NewMockOrdersProvider creates a new mock instance.
func NewMockOrdersProvider(ctrl *gomock.Controller) *MockOrdersProvider {
	mock := &MockOrdersProvider{ctrl: ctrl}
	mock.recorder = &MockOrdersProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrdersProvider) EXPECT() *MockOrdersProviderMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockOrdersProvider) Create(arg0 context.Context, arg1 int64, arg2 models.OrderStatus, arg3 []models.OrderItem) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockOrdersProviderMockRecorder) Create(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrdersProvider)(nil).Create), arg0, arg1, arg2, arg3)
}

// GetExpiredOrdersWithStatus mocks base method.
func (m *MockOrdersProvider) GetExpiredOrdersWithStatus(arg0 context.Context, arg1 time.Time, arg2 models.OrderStatus) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExpiredOrdersWithStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExpiredOrdersWithStatus indicates an expected call of GetExpiredOrdersWithStatus.
func (mr *MockOrdersProviderMockRecorder) GetExpiredOrdersWithStatus(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExpiredOrdersWithStatus", reflect.TypeOf((*MockOrdersProvider)(nil).GetExpiredOrdersWithStatus), arg0, arg1, arg2)
}

// GetOrderByOrderId mocks base method.
func (m *MockOrdersProvider) GetOrderByOrderId(arg0 context.Context, arg1 int64) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByOrderId", arg0, arg1)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByOrderId indicates an expected call of GetOrderByOrderId.
func (mr *MockOrdersProviderMockRecorder) GetOrderByOrderId(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByOrderId", reflect.TypeOf((*MockOrdersProvider)(nil).GetOrderByOrderId), arg0, arg1)
}

// SetStatus mocks base method.
func (m *MockOrdersProvider) SetStatus(arg0 context.Context, arg1 models.Order, arg2 models.OrderStatus) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetStatus indicates an expected call of SetStatus.
func (mr *MockOrdersProviderMockRecorder) SetStatus(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatus", reflect.TypeOf((*MockOrdersProvider)(nil).SetStatus), arg0, arg1, arg2)
}
