// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zitadel/zitadel/internal/eventstore (interfaces: Querier,Pusher)
//
// Generated by this command:
//
//	mockgen -package mock -destination ./repository.mock.go github.com/zitadel/zitadel/internal/eventstore Querier,Pusher
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	decimal "github.com/shopspring/decimal"
	eventstore "github.com/zitadel/zitadel/v2/internal/eventstore"
	gomock "go.uber.org/mock/gomock"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// FilterToReducer mocks base method.
func (m *MockQuerier) FilterToReducer(arg0 context.Context, arg1 *eventstore.SearchQueryBuilder, arg2 eventstore.Reducer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterToReducer", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// FilterToReducer indicates an expected call of FilterToReducer.
func (mr *MockQuerierMockRecorder) FilterToReducer(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterToReducer", reflect.TypeOf((*MockQuerier)(nil).FilterToReducer), arg0, arg1, arg2)
}

// Health mocks base method.
func (m *MockQuerier) Health(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Health indicates an expected call of Health.
func (mr *MockQuerierMockRecorder) Health(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockQuerier)(nil).Health), arg0)
}

// InstanceIDs mocks base method.
func (m *MockQuerier) InstanceIDs(arg0 context.Context, arg1 *eventstore.SearchQueryBuilder) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceIDs", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceIDs indicates an expected call of InstanceIDs.
func (mr *MockQuerierMockRecorder) InstanceIDs(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceIDs", reflect.TypeOf((*MockQuerier)(nil).InstanceIDs), arg0, arg1)
}

// LatestPosition mocks base method.
func (m *MockQuerier) LatestPosition(arg0 context.Context, arg1 *eventstore.SearchQueryBuilder) (decimal.Decimal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LatestPosition", arg0, arg1)
	ret0, _ := ret[0].(decimal.Decimal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LatestPosition indicates an expected call of LatestPosition.
func (mr *MockQuerierMockRecorder) LatestPosition(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestPosition", reflect.TypeOf((*MockQuerier)(nil).LatestPosition), arg0, arg1)
}

// MockPusher is a mock of Pusher interface.
type MockPusher struct {
	ctrl     *gomock.Controller
	recorder *MockPusherMockRecorder
}

// MockPusherMockRecorder is the mock recorder for MockPusher.
type MockPusherMockRecorder struct {
	mock *MockPusher
}

// NewMockPusher creates a new mock instance.
func NewMockPusher(ctrl *gomock.Controller) *MockPusher {
	mock := &MockPusher{ctrl: ctrl}
	mock.recorder = &MockPusherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPusher) EXPECT() *MockPusherMockRecorder {
	return m.recorder
}

// Health mocks base method.
func (m *MockPusher) Health(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Health indicates an expected call of Health.
func (mr *MockPusherMockRecorder) Health(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockPusher)(nil).Health), arg0)
}

// Push mocks base method.
func (m *MockPusher) Push(arg0 context.Context, arg1 ...eventstore.Command) ([]eventstore.Event, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Push", varargs...)
	ret0, _ := ret[0].([]eventstore.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Push indicates an expected call of Push.
func (mr *MockPusherMockRecorder) Push(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockPusher)(nil).Push), varargs...)
}
