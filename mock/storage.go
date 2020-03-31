// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/guilherme-santos/messageboard (interfaces: Storage)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	messageboard "github.com/guilherme-santos/messageboard"
	reflect "reflect"
)

// Storage is a mock of Storage interface
type Storage struct {
	ctrl     *gomock.Controller
	recorder *StorageMockRecorder
}

// StorageMockRecorder is the mock recorder for Storage
type StorageMockRecorder struct {
	mock *Storage
}

// NewStorage creates a new mock instance
func NewStorage(ctrl *gomock.Controller) *Storage {
	mock := &Storage{ctrl: ctrl}
	mock.recorder = &StorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Storage) EXPECT() *StorageMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *Storage) Create(arg0 context.Context, arg1 *messageboard.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *StorageMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Storage)(nil).Create), arg0, arg1)
}

// Get mocks base method
func (m *Storage) Get(arg0 context.Context, arg1 string) (*messageboard.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*messageboard.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *StorageMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*Storage)(nil).Get), arg0, arg1)
}

// List mocks base method
func (m *Storage) List(arg0 context.Context, arg1 *messageboard.ListOptions) (*messageboard.MessageList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*messageboard.MessageList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *StorageMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*Storage)(nil).List), arg0, arg1)
}

// Update mocks base method
func (m *Storage) Update(arg0 context.Context, arg1 *messageboard.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *StorageMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*Storage)(nil).Update), arg0, arg1)
}