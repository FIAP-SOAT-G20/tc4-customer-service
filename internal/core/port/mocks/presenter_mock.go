// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/port/presenter_port.go
//
// Generated by this command:
//
//	mockgen -source=internal/core/port/presenter_port.go -destination=internal/core/port/mocks/presenter_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	dto "github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockPresenter is a mock of Presenter interface.
type MockPresenter struct {
	ctrl     *gomock.Controller
	recorder *MockPresenterMockRecorder
	isgomock struct{}
}

// MockPresenterMockRecorder is the mock recorder for MockPresenter.
type MockPresenterMockRecorder struct {
	mock *MockPresenter
}

// NewMockPresenter creates a new mock instance.
func NewMockPresenter(ctrl *gomock.Controller) *MockPresenter {
	mock := &MockPresenter{ctrl: ctrl}
	mock.recorder = &MockPresenterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPresenter) EXPECT() *MockPresenterMockRecorder {
	return m.recorder
}

// Present mocks base method.
func (m *MockPresenter) Present(arg0 dto.PresenterInput) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Present", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Present indicates an expected call of Present.
func (mr *MockPresenterMockRecorder) Present(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Present", reflect.TypeOf((*MockPresenter)(nil).Present), arg0)
}
