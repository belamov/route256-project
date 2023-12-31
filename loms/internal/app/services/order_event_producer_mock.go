// Code generated by MockGen. DO NOT EDIT.
// Source: route256/loms/internal/app/services (interfaces: OrderEventsProducer)
//
// Generated by this command:
//
//	mockgen -destination=internal/app/services/order_event_producer_mock.go -package=services route256/loms/internal/app/services OrderEventsProducer
//
// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"
	models "route256/loms/internal/app/models"

	gomock "go.uber.org/mock/gomock"
)

// MockOrderEventsProducer is a mock of OrderEventsProducer interface.
type MockOrderEventsProducer struct {
	ctrl     *gomock.Controller
	recorder *MockOrderEventsProducerMockRecorder
}

// MockOrderEventsProducerMockRecorder is the mock recorder for MockOrderEventsProducer.
type MockOrderEventsProducerMockRecorder struct {
	mock *MockOrderEventsProducer
}

// NewMockOrderEventsProducer creates a new mock instance.
func NewMockOrderEventsProducer(ctrl *gomock.Controller) *MockOrderEventsProducer {
	mock := &MockOrderEventsProducer{ctrl: ctrl}
	mock.recorder = &MockOrderEventsProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderEventsProducer) EXPECT() *MockOrderEventsProducerMockRecorder {
	return m.recorder
}

// OrderStatusChangedEventEmit mocks base method.
func (m *MockOrderEventsProducer) OrderStatusChangedEventEmit(arg0 context.Context, arg1 models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrderStatusChangedEventEmit", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// OrderStatusChangedEventEmit indicates an expected call of OrderStatusChangedEventEmit.
func (mr *MockOrderEventsProducerMockRecorder) OrderStatusChangedEventEmit(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderStatusChangedEventEmit", reflect.TypeOf((*MockOrderEventsProducer)(nil).OrderStatusChangedEventEmit), arg0, arg1)
}
