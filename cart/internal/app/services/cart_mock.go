// Code generated by MockGen. DO NOT EDIT.
// Source: route256/cart/internal/app/services (interfaces: Cart)
//
// Generated by this command:
//
//	mockgen -destination=internal/app/services/cart_mock.go -package=services route256/cart/internal/app/services Cart
//
// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"
	models "route256/cart/internal/app/models"

	gomock "go.uber.org/mock/gomock"
)

// MockCart is a mock of Cart interface.
type MockCart struct {
	ctrl     *gomock.Controller
	recorder *MockCartMockRecorder
}

// MockCartMockRecorder is the mock recorder for MockCart.
type MockCartMockRecorder struct {
	mock *MockCart
}

// NewMockCart creates a new mock instance.
func NewMockCart(ctrl *gomock.Controller) *MockCart {
	mock := &MockCart{ctrl: ctrl}
	mock.recorder = &MockCartMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCart) EXPECT() *MockCartMockRecorder {
	return m.recorder
}

// AddItem mocks base method.
func (m *MockCart) AddItem(arg0 context.Context, arg1 models.CartItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddItem indicates an expected call of AddItem.
func (mr *MockCartMockRecorder) AddItem(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddItem", reflect.TypeOf((*MockCart)(nil).AddItem), arg0, arg1)
}

// Checkout mocks base method.
func (m *MockCart) Checkout(arg0 context.Context, arg1 int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checkout", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Checkout indicates an expected call of Checkout.
func (mr *MockCartMockRecorder) Checkout(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checkout", reflect.TypeOf((*MockCart)(nil).Checkout), arg0, arg1)
}

// DeleteItem mocks base method.
func (m *MockCart) DeleteItem(arg0 context.Context, arg1 models.CartItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockCartMockRecorder) DeleteItem(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockCart)(nil).DeleteItem), arg0, arg1)
}

// DeleteItemsByUserId mocks base method.
func (m *MockCart) DeleteItemsByUserId(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItemsByUserId", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItemsByUserId indicates an expected call of DeleteItemsByUserId.
func (mr *MockCartMockRecorder) DeleteItemsByUserId(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItemsByUserId", reflect.TypeOf((*MockCart)(nil).DeleteItemsByUserId), arg0, arg1)
}

// GetItemsByUserId mocks base method.
func (m *MockCart) GetItemsByUserId(arg0 context.Context, arg1 int64) ([]models.CartItemWithInfo, uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItemsByUserId", arg0, arg1)
	ret0, _ := ret[0].([]models.CartItemWithInfo)
	ret1, _ := ret[1].(uint32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetItemsByUserId indicates an expected call of GetItemsByUserId.
func (mr *MockCartMockRecorder) GetItemsByUserId(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemsByUserId", reflect.TypeOf((*MockCart)(nil).GetItemsByUserId), arg0, arg1)
}
