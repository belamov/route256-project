// Code generated by MockGen. DO NOT EDIT.
// Source: route256/loms/internal/app/services (interfaces: Loms)
//
// Generated by this command:
//
//	mockgen -destination=internal/app/mocks/loms.go -package=services route256/loms/internal/app/services Loms
//
// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"
	models "route256/loms/internal/app/models"

	gomock "go.uber.org/mock/gomock"
)

// MockLoms is a mock of Loms interface.
type MockLoms struct {
	ctrl     *gomock.Controller
	recorder *MockLomsMockRecorder
}

// MockLomsMockRecorder is the mock recorder for MockLoms.
type MockLomsMockRecorder struct {
	mock *MockLoms
}

// NewMockLoms creates a new mock instance.
func NewMockLoms(ctrl *gomock.Controller) *MockLoms {
	mock := &MockLoms{ctrl: ctrl}
	mock.recorder = &MockLomsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoms) EXPECT() *MockLomsMockRecorder {
	return m.recorder
}

// OrderCancel mocks base method.
func (m *MockLoms) OrderCancel(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrderCancel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// OrderCancel indicates an expected call of OrderCancel.
func (mr *MockLomsMockRecorder) OrderCancel(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderCancel", reflect.TypeOf((*MockLoms)(nil).OrderCancel), arg0, arg1)
}

// OrderCreate mocks base method.
func (m *MockLoms) OrderCreate(arg0 context.Context, arg1 int64, arg2 []models.OrderItem) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrderCreate", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OrderCreate indicates an expected call of OrderCreate.
func (mr *MockLomsMockRecorder) OrderCreate(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderCreate", reflect.TypeOf((*MockLoms)(nil).OrderCreate), arg0, arg1, arg2)
}

// OrderInfo mocks base method.
func (m *MockLoms) OrderInfo(arg0 context.Context, arg1 int64) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrderInfo", arg0, arg1)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OrderInfo indicates an expected call of OrderInfo.
func (mr *MockLomsMockRecorder) OrderInfo(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderInfo", reflect.TypeOf((*MockLoms)(nil).OrderInfo), arg0, arg1)
}

// OrderPay mocks base method.
func (m *MockLoms) OrderPay(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrderPay", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// OrderPay indicates an expected call of OrderPay.
func (mr *MockLomsMockRecorder) OrderPay(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderPay", reflect.TypeOf((*MockLoms)(nil).OrderPay), arg0, arg1)
}

// StockInfo mocks base method.
func (m *MockLoms) StockInfo(arg0 context.Context, arg1 uint32) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StockInfo", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StockInfo indicates an expected call of StockInfo.
func (mr *MockLomsMockRecorder) StockInfo(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StockInfo", reflect.TypeOf((*MockLoms)(nil).StockInfo), arg0, arg1)
}
