// Package mock is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/osalomon89/test-crud-api/internal/core/domain"
)

// MockItemRepository is a mock of ItemRepository interface
type MockItemRepository struct {
	ctrl     *gomock.Controller
	recorder *MockItemRepositoryMockRecorder
}

// MockItemRepositoryMockRecorder is the mock recorder for MockItemRepository
type MockItemRepositoryMockRecorder struct {
	mock *MockItemRepository
}

// NewMockItemRepository creates a new mock instance
func NewMockItemRepository(ctrl *gomock.Controller) *MockItemRepository {
	mock := &MockItemRepository{ctrl: ctrl}
	mock.recorder = &MockItemRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockItemRepository) EXPECT() *MockItemRepositoryMockRecorder {
	return m.recorder
}

// GetItemByID mocks base method
func (m *MockItemRepository) GetItemByID(arg0 context.Context, arg1 uint) (*domain.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItemByID", arg0, arg1)
	ret0, _ := ret[0].(*domain.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItemByID indicates an expected call of GetItemByID
func (mr *MockItemRepositoryMockRecorder) GetItemByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemByID", reflect.TypeOf((*MockItemRepository)(nil).GetItemByID), arg0, arg1)
}

// SaveItem mocks base method
func (m *MockItemRepository) SaveItem(arg0 context.Context, arg1 *domain.Item) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveItem indicates an expected call of SaveItem
func (mr *MockItemRepositoryMockRecorder) SaveItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveItem", reflect.TypeOf((*MockItemRepository)(nil).SaveItem), arg0, arg1)
}
