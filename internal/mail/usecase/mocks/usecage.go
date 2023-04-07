// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2023_1_Seekers/internal/mail (interfaces: UseCaseI)

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

// CreateDefaultFolders mocks base method.
func (m *MockUseCaseI) CreateDefaultFolders(arg0 uint64) ([]models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDefaultFolders", arg0)
	ret0, _ := ret[0].([]models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDefaultFolders indicates an expected call of CreateDefaultFolders.
func (mr *MockUseCaseIMockRecorder) CreateDefaultFolders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDefaultFolders", reflect.TypeOf((*MockUseCaseI)(nil).CreateDefaultFolders), arg0)
}

// GetFolderInfo mocks base method.
func (m *MockUseCaseI) GetFolderInfo(arg0 uint64, arg1 string) (*models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFolderInfo", arg0, arg1)
	ret0, _ := ret[0].(*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFolderInfo indicates an expected call of GetFolderInfo.
func (mr *MockUseCaseIMockRecorder) GetFolderInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFolderInfo", reflect.TypeOf((*MockUseCaseI)(nil).GetFolderInfo), arg0, arg1)
}

// GetFolderMessages mocks base method.
func (m *MockUseCaseI) GetFolderMessages(arg0 uint64, arg1 string) ([]models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFolderMessages", arg0, arg1)
	ret0, _ := ret[0].([]models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFolderMessages indicates an expected call of GetFolderMessages.
func (mr *MockUseCaseIMockRecorder) GetFolderMessages(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFolderMessages", reflect.TypeOf((*MockUseCaseI)(nil).GetFolderMessages), arg0, arg1)
}

// GetFolders mocks base method.
func (m *MockUseCaseI) GetFolders(arg0 uint64) ([]models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFolders", arg0)
	ret0, _ := ret[0].([]models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFolders indicates an expected call of GetFolders.
func (mr *MockUseCaseIMockRecorder) GetFolders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFolders", reflect.TypeOf((*MockUseCaseI)(nil).GetFolders), arg0)
}

// GetMessage mocks base method.
func (m *MockUseCaseI) GetMessage(arg0, arg1 uint64) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessage", arg0, arg1)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessage indicates an expected call of GetMessage.
func (mr *MockUseCaseIMockRecorder) GetMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessage", reflect.TypeOf((*MockUseCaseI)(nil).GetMessage), arg0, arg1)
}

// MarkMessageAsSeen mocks base method.
func (m *MockUseCaseI) MarkMessageAsSeen(arg0, arg1 uint64) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkMessageAsSeen", arg0, arg1)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkMessageAsSeen indicates an expected call of MarkMessageAsSeen.
func (mr *MockUseCaseIMockRecorder) MarkMessageAsSeen(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkMessageAsSeen", reflect.TypeOf((*MockUseCaseI)(nil).MarkMessageAsSeen), arg0, arg1)
}

// MarkMessageAsUnseen mocks base method.
func (m *MockUseCaseI) MarkMessageAsUnseen(arg0, arg1 uint64) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkMessageAsUnseen", arg0, arg1)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkMessageAsUnseen indicates an expected call of MarkMessageAsUnseen.
func (mr *MockUseCaseIMockRecorder) MarkMessageAsUnseen(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkMessageAsUnseen", reflect.TypeOf((*MockUseCaseI)(nil).MarkMessageAsUnseen), arg0, arg1)
}

// SendFailedSendingMessage mocks base method.
func (m *MockUseCaseI) SendFailedSendingMessage(arg0 string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendFailedSendingMessage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendFailedSendingMessage indicates an expected call of SendFailedSendingMessage.
func (mr *MockUseCaseIMockRecorder) SendFailedSendingMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendFailedSendingMessage", reflect.TypeOf((*MockUseCaseI)(nil).SendFailedSendingMessage), arg0, arg1)
}

// SendMessage mocks base method.
func (m *MockUseCaseI) SendMessage(arg0 uint64, arg1 models.FormMessage) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", arg0, arg1)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockUseCaseIMockRecorder) SendMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockUseCaseI)(nil).SendMessage), arg0, arg1)
}

// SendWelcomeMessage mocks base method.
func (m *MockUseCaseI) SendWelcomeMessage(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendWelcomeMessage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendWelcomeMessage indicates an expected call of SendWelcomeMessage.
func (mr *MockUseCaseIMockRecorder) SendWelcomeMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendWelcomeMessage", reflect.TypeOf((*MockUseCaseI)(nil).SendWelcomeMessage), arg0)
}

// ValidateRecipients mocks base method.
func (m *MockUseCaseI) ValidateRecipients(arg0 []string) ([]string, []string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateRecipients", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// ValidateRecipients indicates an expected call of ValidateRecipients.
func (mr *MockUseCaseIMockRecorder) ValidateRecipients(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateRecipients", reflect.TypeOf((*MockUseCaseI)(nil).ValidateRecipients), arg0)
}
