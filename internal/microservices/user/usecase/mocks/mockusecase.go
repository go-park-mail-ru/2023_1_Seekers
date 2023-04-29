// Code generated by MockGen. DO NOT EDIT.
// Source: ../interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUseCaseI is a mock of UseCaseI interface.
type MockUseCaseI struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseIMockRecorder
}

// MockUseCaseIMockRecorder is the mock recorder for MockUseCaseI.
type MockUseCaseIMockRecorder struct {
	mock *MockUseCaseI
}

// NewMockUseCaseI creates a new mock instance.
func NewMockUseCaseI(ctrl *gomock.Controller) *MockUseCaseI {
	mock := &MockUseCaseI{ctrl: ctrl}
	mock.recorder = &MockUseCaseIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCaseI) EXPECT() *MockUseCaseIMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUseCaseI) Create(user *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUseCaseIMockRecorder) Create(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUseCaseI)(nil).Create), user)
}

// Delete mocks base method.
func (m *MockUseCaseI) Delete(ID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUseCaseIMockRecorder) Delete(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUseCaseI)(nil).Delete), ID)
}

// EditAvatar mocks base method.
func (m *MockUseCaseI) EditAvatar(ID uint64, newAvatar *models.Image, isCustom bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditAvatar", ID, newAvatar, isCustom)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditAvatar indicates an expected call of EditAvatar.
func (mr *MockUseCaseIMockRecorder) EditAvatar(ID, newAvatar, isCustom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditAvatar", reflect.TypeOf((*MockUseCaseI)(nil).EditAvatar), ID, newAvatar, isCustom)
}

// EditInfo mocks base method.
func (m *MockUseCaseI) EditInfo(ID uint64, info *models.UserInfo) (*models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditInfo", ID, info)
	ret0, _ := ret[0].(*models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditInfo indicates an expected call of EditInfo.
func (mr *MockUseCaseIMockRecorder) EditInfo(ID, info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditInfo", reflect.TypeOf((*MockUseCaseI)(nil).EditInfo), ID, info)
}

// EditPw mocks base method.
func (m *MockUseCaseI) EditPw(ID uint64, form *models.EditPasswordRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditPw", ID, form)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditPw indicates an expected call of EditPw.
func (mr *MockUseCaseIMockRecorder) EditPw(ID, form interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditPw", reflect.TypeOf((*MockUseCaseI)(nil).EditPw), ID, form)
}

// GetAvatar mocks base method.
func (m *MockUseCaseI) GetAvatar(email string) (*models.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvatar", email)
	ret0, _ := ret[0].(*models.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvatar indicates an expected call of GetAvatar.
func (mr *MockUseCaseIMockRecorder) GetAvatar(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvatar", reflect.TypeOf((*MockUseCaseI)(nil).GetAvatar), email)
}

// GetByEmail mocks base method.
func (m *MockUseCaseI) GetByEmail(email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockUseCaseIMockRecorder) GetByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUseCaseI)(nil).GetByEmail), email)
}

// GetByID mocks base method.
func (m *MockUseCaseI) GetByID(ID uint64) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUseCaseIMockRecorder) GetByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUseCaseI)(nil).GetByID), ID)
}

// GetInfo mocks base method.
func (m *MockUseCaseI) GetInfo(ID uint64) (*models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInfo", ID)
	ret0, _ := ret[0].(*models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInfo indicates an expected call of GetInfo.
func (mr *MockUseCaseIMockRecorder) GetInfo(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfo", reflect.TypeOf((*MockUseCaseI)(nil).GetInfo), ID)
}

// GetInfoByEmail mocks base method.
func (m *MockUseCaseI) GetInfoByEmail(email string) (*models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInfoByEmail", email)
	ret0, _ := ret[0].(*models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInfoByEmail indicates an expected call of GetInfoByEmail.
func (mr *MockUseCaseIMockRecorder) GetInfoByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfoByEmail", reflect.TypeOf((*MockUseCaseI)(nil).GetInfoByEmail), email)
}
