// Code generated by MockGen. DO NOT EDIT.
// Source: ./interface.go

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

// CreateAnonymousEmail mocks base method.
func (m *MockUseCaseI) CreateAnonymousEmail(userID uint64) (*models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAnonymousEmail", userID)
	ret0, _ := ret[0].(*models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAnonymousEmail indicates an expected call of CreateAnonymousEmail.
func (mr *MockUseCaseIMockRecorder) CreateAnonymousEmail(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAnonymousEmail", reflect.TypeOf((*MockUseCaseI)(nil).CreateAnonymousEmail), userID)
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

// CreateFolder mocks base method.
func (m *MockUseCaseI) CreateFolder(userID uint64, form models.FormFolder) (*models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFolder", userID, form)
	ret0, _ := ret[0].(*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFolder indicates an expected call of CreateFolder.
func (mr *MockUseCaseIMockRecorder) CreateFolder(userID, form interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFolder", reflect.TypeOf((*MockUseCaseI)(nil).CreateFolder), userID, form)
}

// DeleteFolder mocks base method.
func (m *MockUseCaseI) DeleteFolder(userID uint64, folderSlug string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFolder", userID, folderSlug)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFolder indicates an expected call of DeleteFolder.
func (mr *MockUseCaseIMockRecorder) DeleteFolder(userID, folderSlug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFolder", reflect.TypeOf((*MockUseCaseI)(nil).DeleteFolder), userID, folderSlug)
}

// DeleteMessage mocks base method.
func (m *MockUseCaseI) DeleteMessage(userID, messageID uint64, folderSlug string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessage", userID, messageID, folderSlug)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessage indicates an expected call of DeleteMessage.
func (mr *MockUseCaseIMockRecorder) DeleteMessage(userID, messageID, folderSlug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MockUseCaseI)(nil).DeleteMessage), userID, messageID, folderSlug)
}

// EditDraft mocks base method.
func (m *MockUseCaseI) EditDraft(userID, messageID uint64, message models.FormEditMessage) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditDraft", userID, messageID, message)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditDraft indicates an expected call of EditDraft.
func (mr *MockUseCaseIMockRecorder) EditDraft(userID, messageID, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditDraft", reflect.TypeOf((*MockUseCaseI)(nil).EditDraft), userID, messageID, message)
}

// EditFolder mocks base method.
func (m *MockUseCaseI) EditFolder(userID uint64, folderSlug string, form models.FormFolder) (*models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditFolder", userID, folderSlug, form)
	ret0, _ := ret[0].(*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditFolder indicates an expected call of EditFolder.
func (mr *MockUseCaseIMockRecorder) EditFolder(userID, folderSlug, form interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditFolder", reflect.TypeOf((*MockUseCaseI)(nil).EditFolder), userID, folderSlug, form)
}

// GetAttach mocks base method.
func (m *MockUseCaseI) GetAttach(attachID, userID uint64) (*models.AttachmentInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttach", attachID, userID)
	ret0, _ := ret[0].(*models.AttachmentInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttach indicates an expected call of GetAttach.
func (mr *MockUseCaseIMockRecorder) GetAttach(attachID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttach", reflect.TypeOf((*MockUseCaseI)(nil).GetAttach), attachID, userID)
}

// GetAttachInfo mocks base method.
func (m *MockUseCaseI) GetAttachInfo(attachID, userID uint64) (*models.AttachmentInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttachInfo", attachID, userID)
	ret0, _ := ret[0].(*models.AttachmentInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttachInfo indicates an expected call of GetAttachInfo.
func (mr *MockUseCaseIMockRecorder) GetAttachInfo(attachID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttachInfo", reflect.TypeOf((*MockUseCaseI)(nil).GetAttachInfo), attachID, userID)
}

// GetCustomFolders mocks base method.
func (m *MockUseCaseI) GetCustomFolders(userID uint64) ([]models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomFolders", userID)
	ret0, _ := ret[0].([]models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCustomFolders indicates an expected call of GetCustomFolders.
func (mr *MockUseCaseIMockRecorder) GetCustomFolders(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomFolders", reflect.TypeOf((*MockUseCaseI)(nil).GetCustomFolders), userID)
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
func (m *MockUseCaseI) MarkMessageAsSeen(userID, messageID uint64, folderSlug string) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkMessageAsSeen", userID, messageID, folderSlug)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkMessageAsSeen indicates an expected call of MarkMessageAsSeen.
func (mr *MockUseCaseIMockRecorder) MarkMessageAsSeen(userID, messageID, folderSlug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkMessageAsSeen", reflect.TypeOf((*MockUseCaseI)(nil).MarkMessageAsSeen), userID, messageID, folderSlug)
}

// MarkMessageAsUnseen mocks base method.
func (m *MockUseCaseI) MarkMessageAsUnseen(userID, messageID uint64, folderSlug string) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkMessageAsUnseen", userID, messageID, folderSlug)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkMessageAsUnseen indicates an expected call of MarkMessageAsUnseen.
func (mr *MockUseCaseIMockRecorder) MarkMessageAsUnseen(userID, messageID, folderSlug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkMessageAsUnseen", reflect.TypeOf((*MockUseCaseI)(nil).MarkMessageAsUnseen), userID, messageID, folderSlug)
}

// MoveMessageToFolder mocks base method.
func (m *MockUseCaseI) MoveMessageToFolder(userID, messageID uint64, fromFolder, toFolder string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MoveMessageToFolder", userID, messageID, fromFolder, toFolder)
	ret0, _ := ret[0].(error)
	return ret0
}

// MoveMessageToFolder indicates an expected call of MoveMessageToFolder.
func (mr *MockUseCaseIMockRecorder) MoveMessageToFolder(userID, messageID, fromFolder, toFolder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MoveMessageToFolder", reflect.TypeOf((*MockUseCaseI)(nil).MoveMessageToFolder), userID, messageID, fromFolder, toFolder)
}

// SaveDraft mocks base method.
func (m *MockUseCaseI) SaveDraft(userID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDraft", userID, message)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveDraft indicates an expected call of SaveDraft.
func (mr *MockUseCaseIMockRecorder) SaveDraft(userID, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDraft", reflect.TypeOf((*MockUseCaseI)(nil).SaveDraft), userID, message)
}

// SearchMessages mocks base method.
func (m *MockUseCaseI) SearchMessages(userID uint64, fromUser, toUser, folder, filter string) ([]models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchMessages", userID, fromUser, toUser, folder, filter)
	ret0, _ := ret[0].([]models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchMessages indicates an expected call of SearchMessages.
func (mr *MockUseCaseIMockRecorder) SearchMessages(userID, fromUser, toUser, folder, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchMessages", reflect.TypeOf((*MockUseCaseI)(nil).SearchMessages), userID, fromUser, toUser, folder, filter)
}

// SearchRecipients mocks base method.
func (m *MockUseCaseI) SearchRecipients(userID uint64) ([]models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRecipients", userID)
	ret0, _ := ret[0].([]models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRecipients indicates an expected call of SearchRecipients.
func (mr *MockUseCaseIMockRecorder) SearchRecipients(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRecipients", reflect.TypeOf((*MockUseCaseI)(nil).SearchRecipients), userID)
}

// SendFailedSendingMessage mocks base method.
func (m *MockUseCaseI) SendFailedSendingMessage(recipientEmail string, invalidEmails []string) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendFailedSendingMessage", recipientEmail, invalidEmails)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendFailedSendingMessage indicates an expected call of SendFailedSendingMessage.
func (mr *MockUseCaseIMockRecorder) SendFailedSendingMessage(recipientEmail, invalidEmails interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendFailedSendingMessage", reflect.TypeOf((*MockUseCaseI)(nil).SendFailedSendingMessage), recipientEmail, invalidEmails)
}

// SendMessage mocks base method.
func (m *MockUseCaseI) SendMessage(userId uint64, message models.FormMessage) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", userId, message)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockUseCaseIMockRecorder) SendMessage(userId, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockUseCaseI)(nil).SendMessage), userId, message)
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
