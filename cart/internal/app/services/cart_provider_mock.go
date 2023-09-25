// Code generated by MockGen. DO NOT EDIT.
// Source: route256/cart/internal/app/services (interfaces: CartProvider)
//
// Generated by this command:
//
//	mockgen -destination=internal/app/services/cart_provider_mock.go -package=services route256/cart/internal/app/services CartProvider
//
// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"
	models "route256/cart/internal/app/models"

	gomock "go.uber.org/mock/gomock"
)

// MockCartProvider is a mock of CartProvider interface.
type MockCartProvider struct {
	ctrl     *gomock.Controller
	recorder *MockCartProviderMockRecorder
}

// MockCartProviderMockRecorder is the mock recorder for MockCartProvider.
type MockCartProviderMockRecorder struct {
	mock *MockCartProvider
}

// NewMockCartProvider creates a new mock instance.
func NewMockCartProvider(ctrl *gomock.Controller) *MockCartProvider {
	mock := &MockCartProvider{ctrl: ctrl}
	mock.recorder = &MockCartProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCartProvider) EXPECT() *MockCartProviderMockRecorder {
	return m.recorder
}

// DeleteItem mocks base method.
func (m *MockCartProvider) DeleteItem(arg0 context.Context, arg1 models.CartItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockCartProviderMockRecorder) DeleteItem(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockCartProvider)(nil).DeleteItem), arg0, arg1)
}

// DeleteItemsByUserId mocks base method.
func (m *MockCartProvider) DeleteItemsByUserId(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItemsByUserId", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItemsByUserId indicates an expected call of DeleteItemsByUserId.
func (mr *MockCartProviderMockRecorder) DeleteItemsByUserId(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItemsByUserId", reflect.TypeOf((*MockCartProvider)(nil).DeleteItemsByUserId), arg0, arg1)
}

// GetItemsByUserId mocks base method.
func (m *MockCartProvider) GetItemsByUserId(arg0 context.Context, arg1 int64) ([]models.CartItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItemsByUserId", arg0, arg1)
	ret0, _ := ret[0].([]models.CartItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItemsByUserId indicates an expected call of GetItemsByUserId.
func (mr *MockCartProviderMockRecorder) GetItemsByUserId(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemsByUserId", reflect.TypeOf((*MockCartProvider)(nil).GetItemsByUserId), arg0, arg1)
}

// SaveItem mocks base method.
func (m *MockCartProvider) SaveItem(arg0 context.Context, arg1 models.CartItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveItem indicates an expected call of SaveItem.
func (mr *MockCartProviderMockRecorder) SaveItem(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveItem", reflect.TypeOf((*MockCartProvider)(nil).SaveItem), arg0, arg1)
}
