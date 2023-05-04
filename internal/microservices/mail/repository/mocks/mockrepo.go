// Code generated by MockGen. DO NOT EDIT.
// Source: ./interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockMailRepoI is a mock of MailRepoI interface.
type MockMailRepoI struct {
	ctrl     *gomock.Controller
	recorder *MockMailRepoIMockRecorder
}

// MockMailRepoIMockRecorder is the mock recorder for MockMailRepoI.
type MockMailRepoIMockRecorder struct {
	mock *MockMailRepoI
}

// NewMockMailRepoI creates a new mock instance.
func NewMockMailRepoI(ctrl *gomock.Controller) *MockMailRepoI {
	mock := &MockMailRepoI{ctrl: ctrl}
	mock.recorder = &MockMailRepoIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMailRepoI) EXPECT() *MockMailRepoIMockRecorder {
	return m.recorder
}

// CheckExistingBox mocks base method.
func (m *MockMailRepoI) CheckExistingBox(userID, messageID, folderID uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckExistingBox", userID, messageID, folderID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckExistingBox indicates an expected call of CheckExistingBox.
func (mr *MockMailRepoIMockRecorder) CheckExistingBox(userID, messageID, folderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckExistingBox", reflect.TypeOf((*MockMailRepoI)(nil).CheckExistingBox), userID, messageID, folderID)
}

// DeleteBox mocks base method.
func (m *MockMailRepoI) DeleteBox(userID, messageID, folderID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBox", userID, messageID, folderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBox indicates an expected call of DeleteBox.
func (mr *MockMailRepoIMockRecorder) DeleteBox(userID, messageID, folderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBox", reflect.TypeOf((*MockMailRepoI)(nil).DeleteBox), userID, messageID, folderID)
}

// DeleteFolder mocks base method.
func (m *MockMailRepoI) DeleteFolder(folderID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFolder", folderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFolder indicates an expected call of DeleteFolder.
func (mr *MockMailRepoIMockRecorder) DeleteFolder(folderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFolder", reflect.TypeOf((*MockMailRepoI)(nil).DeleteFolder), folderID)
}

// DeleteMessageFromMessages mocks base method.
func (m *MockMailRepoI) DeleteMessageFromMessages(messageID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessageFromMessages", messageID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessageFromMessages indicates an expected call of DeleteMessageFromMessages.
func (mr *MockMailRepoIMockRecorder) DeleteMessageFromMessages(messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessageFromMessages", reflect.TypeOf((*MockMailRepoI)(nil).DeleteMessageFromMessages), messageID)
}

// InsertFolder mocks base method.
func (m *MockMailRepoI) InsertFolder(folder *models.Folder) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertFolder", folder)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertFolder indicates an expected call of InsertFolder.
func (mr *MockMailRepoIMockRecorder) InsertFolder(folder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertFolder", reflect.TypeOf((*MockMailRepoI)(nil).InsertFolder), folder)
}

// InsertMessage mocks base method.
func (m *MockMailRepoI) InsertMessage(fromUserID uint64, message *models.MessageInfo, user2folder []models.User2Folder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMessage", fromUserID, message, user2folder)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertMessage indicates an expected call of InsertMessage.
func (mr *MockMailRepoIMockRecorder) InsertMessage(fromUserID, message, user2folder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMessage", reflect.TypeOf((*MockMailRepoI)(nil).InsertMessage), fromUserID, message, user2folder)
}

// SelectCustomFoldersByUser mocks base method.
func (m *MockMailRepoI) SelectCustomFoldersByUser(userID uint64, defaultLocalNames []string) ([]models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectCustomFoldersByUser", userID, defaultLocalNames)
	ret0, _ := ret[0].([]models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectCustomFoldersByUser indicates an expected call of SelectCustomFoldersByUser.
func (mr *MockMailRepoIMockRecorder) SelectCustomFoldersByUser(userID, defaultLocalNames interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectCustomFoldersByUser", reflect.TypeOf((*MockMailRepoI)(nil).SelectCustomFoldersByUser), userID, defaultLocalNames)
}

// SelectFolderByUserNFolderName mocks base method.
func (m *MockMailRepoI) SelectFolderByUserNFolderName(userID uint64, folderName string) (*models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectFolderByUserNFolderName", userID, folderName)
	ret0, _ := ret[0].(*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectFolderByUserNFolderName indicates an expected call of SelectFolderByUserNFolderName.
func (mr *MockMailRepoIMockRecorder) SelectFolderByUserNFolderName(userID, folderName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectFolderByUserNFolderName", reflect.TypeOf((*MockMailRepoI)(nil).SelectFolderByUserNFolderName), userID, folderName)
}

// SelectFolderByUserNFolderSlug mocks base method.
func (m *MockMailRepoI) SelectFolderByUserNFolderSlug(userID uint64, folderSlug string) (*models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectFolderByUserNFolderSlug", userID, folderSlug)
	ret0, _ := ret[0].(*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectFolderByUserNFolderSlug indicates an expected call of SelectFolderByUserNFolderSlug.
func (mr *MockMailRepoIMockRecorder) SelectFolderByUserNFolderSlug(userID, folderSlug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectFolderByUserNFolderSlug", reflect.TypeOf((*MockMailRepoI)(nil).SelectFolderByUserNFolderSlug), userID, folderSlug)
}

// SelectFolderByUserNMessage mocks base method.
func (m *MockMailRepoI) SelectFolderByUserNMessage(userID, messageID uint64) (*models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectFolderByUserNMessage", userID, messageID)
	ret0, _ := ret[0].(*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectFolderByUserNMessage indicates an expected call of SelectFolderByUserNMessage.
func (mr *MockMailRepoIMockRecorder) SelectFolderByUserNMessage(userID, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectFolderByUserNMessage", reflect.TypeOf((*MockMailRepoI)(nil).SelectFolderByUserNMessage), userID, messageID)
}

// SelectFolderMessagesByUserNFolderID mocks base method.
func (m *MockMailRepoI) SelectFolderMessagesByUserNFolderID(userID, folderID uint64, isDraft bool) ([]models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectFolderMessagesByUserNFolderID", userID, folderID, isDraft)
	ret0, _ := ret[0].([]models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectFolderMessagesByUserNFolderID indicates an expected call of SelectFolderMessagesByUserNFolderID.
func (mr *MockMailRepoIMockRecorder) SelectFolderMessagesByUserNFolderID(userID, folderID, isDraft interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectFolderMessagesByUserNFolderID", reflect.TypeOf((*MockMailRepoI)(nil).SelectFolderMessagesByUserNFolderID), userID, folderID, isDraft)
}

// SelectFoldersByUser mocks base method.
func (m *MockMailRepoI) SelectFoldersByUser(userID uint64) ([]models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectFoldersByUser", userID)
	ret0, _ := ret[0].([]models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectFoldersByUser indicates an expected call of SelectFoldersByUser.
func (mr *MockMailRepoIMockRecorder) SelectFoldersByUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectFoldersByUser", reflect.TypeOf((*MockMailRepoI)(nil).SelectFoldersByUser), userID)
}

// SelectMessageByUserNMessage mocks base method.
func (m *MockMailRepoI) SelectMessageByUserNMessage(userID, messageID uint64) (*models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectMessageByUserNMessage", userID, messageID)
	ret0, _ := ret[0].(*models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectMessageByUserNMessage indicates an expected call of SelectMessageByUserNMessage.
func (mr *MockMailRepoIMockRecorder) SelectMessageByUserNMessage(userID, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectMessageByUserNMessage", reflect.TypeOf((*MockMailRepoI)(nil).SelectMessageByUserNMessage), userID, messageID)
}

// SelectRecipientsByMessage mocks base method.
func (m *MockMailRepoI) SelectRecipientsByMessage(messageID, fromUserID uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectRecipientsByMessage", messageID, fromUserID)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectRecipientsByMessage indicates an expected call of SelectRecipientsByMessage.
func (mr *MockMailRepoIMockRecorder) SelectRecipientsByMessage(messageID, fromUserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectRecipientsByMessage", reflect.TypeOf((*MockMailRepoI)(nil).SelectRecipientsByMessage), messageID, fromUserID)
}

// UpdateFolder mocks base method.
func (m *MockMailRepoI) UpdateFolder(folder models.Folder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFolder", folder)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFolder indicates an expected call of UpdateFolder.
func (mr *MockMailRepoIMockRecorder) UpdateFolder(folder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFolder", reflect.TypeOf((*MockMailRepoI)(nil).UpdateFolder), folder)
}

// UpdateMessage mocks base method.
func (m *MockMailRepoI) UpdateMessage(message *models.MessageInfo, toInsert, toDelete []models.User2Folder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessage", message, toInsert, toDelete)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMessage indicates an expected call of UpdateMessage.
func (mr *MockMailRepoIMockRecorder) UpdateMessage(message, toInsert, toDelete interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessage", reflect.TypeOf((*MockMailRepoI)(nil).UpdateMessage), message, toInsert, toDelete)
}

// UpdateMessageFolder mocks base method.
func (m *MockMailRepoI) UpdateMessageFolder(userID, messageID, oldFolderID, newFolderID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessageFolder", userID, messageID, oldFolderID, newFolderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMessageFolder indicates an expected call of UpdateMessageFolder.
func (mr *MockMailRepoIMockRecorder) UpdateMessageFolder(userID, messageID, oldFolderID, newFolderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessageFolder", reflect.TypeOf((*MockMailRepoI)(nil).UpdateMessageFolder), userID, messageID, oldFolderID, newFolderID)
}

// UpdateMessageState mocks base method.
func (m *MockMailRepoI) UpdateMessageState(userID, messageID, folderID uint64, stateName string, stateValue bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessageState", userID, messageID, folderID, stateName, stateValue)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMessageState indicates an expected call of UpdateMessageState.
func (mr *MockMailRepoIMockRecorder) UpdateMessageState(userID, messageID, folderID, stateName, stateValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessageState", reflect.TypeOf((*MockMailRepoI)(nil).UpdateMessageState), userID, messageID, folderID, stateName, stateValue)
}
