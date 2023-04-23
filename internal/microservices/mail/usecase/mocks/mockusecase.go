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

// CreateDefaultFolders mocks base method.
func (m *MockUseCaseI) CreateDefaultFolders(userID uint64) ([]models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDefaultFolders", userID)
	ret0, _ := ret[0].([]models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDefaultFolders indicates an expected call of CreateDefaultFolders.
func (mr *MockUseCaseIMockRecorder) CreateDefaultFolders(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDefaultFolders", reflect.TypeOf((*MockUseCaseI)(nil).CreateDefaultFolders), userID)
}

// GetFolderInfo mocks base method.
func (m *MockUseCaseI) GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFolderInfo", userID, folderSlug)
	ret0, _ := ret[0].(*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFolderInfo indicates an expected call of GetFolderInfo.
func (mr *MockUseCaseIMockRecorder) GetFolderInfo(userID, folderSlug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFolderInfo", reflect.TypeOf((*MockUseCaseI)(nil).GetFolderInfo), userID, folderSlug)
}

// GetFolderMessages mocks base method.
func (m *MockUseCaseI) GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFolderMessages", userID, folderSlug)
	ret0, _ := ret[0].([]models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFolderMessages indicates an expected call of GetFolderMessages.
func (mr *MockUseCaseIMockRecorder) GetFolderMessages(userID, folderSlug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFolderMessages", reflect.TypeOf((*MockUseCaseI)(nil).GetFolderMessages), userID, folderSlug)
}

// GetFolders mocks base method.
func (m *MockUseCaseI) GetFolders(userID uint64) ([]models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFolders", userID)
	ret0, _ := ret[0].([]models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFolders indicates an expected call of GetFolders.
func (mr *MockUseCaseIMockRecorder) GetFolders(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFolders", reflect.TypeOf((*MockUseCaseI)(nil).GetFolders), userID)
}

// GetMessage mocks base method.
func (m *MockUseCaseI) GetMessage(userID, messageID uint64) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessage", userID, messageID)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessage indicates an expected call of GetMessage.
func (mr *MockUseCaseIMockRecorder) GetMessage(userID, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessage", reflect.TypeOf((*MockUseCaseI)(nil).GetMessage), userID, messageID)
}

// MarkMessageAsSeen mocks base method.
func (m *MockUseCaseI) MarkMessageAsSeen(userID, messageID uint64) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkMessageAsSeen", userID, messageID)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkMessageAsSeen indicates an expected call of MarkMessageAsSeen.
func (mr *MockUseCaseIMockRecorder) MarkMessageAsSeen(userID, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkMessageAsSeen", reflect.TypeOf((*MockUseCaseI)(nil).MarkMessageAsSeen), userID, messageID)
}

// MarkMessageAsUnseen mocks base method.
func (m *MockUseCaseI) MarkMessageAsUnseen(userID, messageID uint64) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkMessageAsUnseen", userID, messageID)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkMessageAsUnseen indicates an expected call of MarkMessageAsUnseen.
func (mr *MockUseCaseIMockRecorder) MarkMessageAsUnseen(userID, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkMessageAsUnseen", reflect.TypeOf((*MockUseCaseI)(nil).MarkMessageAsUnseen), userID, messageID)
}

// SendFailedSendingMessage mocks base method.
func (m *MockUseCaseI) SendFailedSendingMessage(recipientEmail string, invalidEmails []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendFailedSendingMessage", recipientEmail, invalidEmails)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendFailedSendingMessage indicates an expected call of SendFailedSendingMessage.
func (mr *MockUseCaseIMockRecorder) SendFailedSendingMessage(recipientEmail, invalidEmails interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendFailedSendingMessage", reflect.TypeOf((*MockUseCaseI)(nil).SendFailedSendingMessage), recipientEmail, invalidEmails)
}

// SendMessage mocks base method.
func (m *MockUseCaseI) SendMessage(userID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", userID, message)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockUseCaseIMockRecorder) SendMessage(userID, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockUseCaseI)(nil).SendMessage), userID, message)
}

// SendWelcomeMessage mocks base method.
func (m *MockUseCaseI) SendWelcomeMessage(recipientEmail string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendWelcomeMessage", recipientEmail)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendWelcomeMessage indicates an expected call of SendWelcomeMessage.
func (mr *MockUseCaseIMockRecorder) SendWelcomeMessage(recipientEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendWelcomeMessage", reflect.TypeOf((*MockUseCaseI)(nil).SendWelcomeMessage), recipientEmail)
}

// ValidateRecipients mocks base method.
func (m *MockUseCaseI) ValidateRecipients(recipients []string) ([]string, []string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateRecipients", recipients)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// ValidateRecipients indicates an expected call of ValidateRecipients.
func (mr *MockUseCaseIMockRecorder) ValidateRecipients(recipients interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateRecipients", reflect.TypeOf((*MockUseCaseI)(nil).ValidateRecipients), recipients)
}