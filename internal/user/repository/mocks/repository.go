	// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2023_1_Seekers/internal/user (interfaces: RepoI)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRepoI is a mock of RepoI interface.
type MockRepoI struct {
	ctrl     *gomock.Controller
	recorder *MockRepoIMockRecorder
}

// MockRepoIMockRecorder is the mock recorder for MockRepoI.
type MockRepoIMockRecorder struct {
	mock *MockRepoI
}

// NewMockRepoI creates a new mock instance.
func NewMockRepoI(ctrl *gomock.Controller) *MockRepoI {
	mock := &MockRepoI{ctrl: ctrl}
	mock.recorder = &MockRepoIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepoI) EXPECT() *MockRepoIMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepoI) Create(arg0 *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepoIMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepoI)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockRepoI) Delete(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepoIMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepoI)(nil).Delete), arg0)
}

// EditInfo mocks base method.
func (m *MockRepoI) EditInfo(arg0 uint64, arg1 models.UserInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditInfo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditInfo indicates an expected call of EditInfo.
func (mr *MockRepoIMockRecorder) EditInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditInfo", reflect.TypeOf((*MockRepoI)(nil).EditInfo), arg0, arg1)
}

// EditPw mocks base method.
func (m *MockRepoI) EditPw(arg0 uint64, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditPw", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditPw indicates an expected call of EditPw.
func (mr *MockRepoIMockRecorder) EditPw(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditPw", reflect.TypeOf((*MockRepoI)(nil).EditPw), arg0, arg1)
}

// GetByEmail mocks base method.
func (m *MockRepoI) GetByEmail(arg0 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockRepoIMockRecorder) GetByEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockRepoI)(nil).GetByEmail), arg0)
}

// GetByID mocks base method.
func (m *MockRepoI) GetByID(arg0 uint64) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepoIMockRecorder) GetByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepoI)(nil).GetByID), arg0)
}

// GetInfoByEmail mocks base method.
func (m *MockRepoI) GetInfoByEmail(arg0 string) (*models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInfoByEmail", arg0)
	ret0, _ := ret[0].(*models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInfoByEmail indicates an expected call of GetInfoByEmail.
func (mr *MockRepoIMockRecorder) GetInfoByEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfoByEmail", reflect.TypeOf((*MockRepoI)(nil).GetInfoByEmail), arg0)
}

// GetInfoByID mocks base method.
func (m *MockRepoI) GetInfoByID(arg0 uint64) (*models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInfoByID", arg0)
	ret0, _ := ret[0].(*models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInfoByID indicates an expected call of GetInfoByID.
func (mr *MockRepoIMockRecorder) GetInfoByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfoByID", reflect.TypeOf((*MockRepoI)(nil).GetInfoByID), arg0)
}

// SetAvatar mocks base method.
func (m *MockRepoI) SetAvatar(arg0 uint64, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAvatar", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAvatar indicates an expected call of SetAvatar.
func (mr *MockRepoIMockRecorder) SetAvatar(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAvatar", reflect.TypeOf((*MockRepoI)(nil).SetAvatar), arg0, arg1)
}
