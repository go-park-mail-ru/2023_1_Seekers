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

// SelectFolderByUserNFolder mocks base method.
func (m *MockMailRepoI) SelectFolderByUserNFolder(userID uint64, folderSlug string) (*models.Folder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectFolderByUserNFolder", userID, folderSlug)
	ret0, _ := ret[0].(*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectFolderByUserNFolder indicates an expected call of SelectFolderByUserNFolder.
func (mr *MockMailRepoIMockRecorder) SelectFolderByUserNFolder(userID, folderSlug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectFolderByUserNFolder", reflect.TypeOf((*MockMailRepoI)(nil).SelectFolderByUserNFolder), userID, folderSlug)
}

// SelectFolderMessagesByUserNFolder mocks base method.
func (m *MockMailRepoI) SelectFolderMessagesByUserNFolder(userID, folderID uint64) ([]models.MessageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectFolderMessagesByUserNFolder", userID, folderID)
	ret0, _ := ret[0].([]models.MessageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectFolderMessagesByUserNFolder indicates an expected call of SelectFolderMessagesByUserNFolder.
func (mr *MockMailRepoIMockRecorder) SelectFolderMessagesByUserNFolder(userID, folderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectFolderMessagesByUserNFolder", reflect.TypeOf((*MockMailRepoI)(nil).SelectFolderMessagesByUserNFolder), userID, folderID)
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

// UpdateMessageState mocks base method.
func (m *MockMailRepoI) UpdateMessageState(userID, messageID uint64, stateName string, stateValue bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessageState", userID, messageID, stateName, stateValue)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMessageState indicates an expected call of UpdateMessageState.
func (mr *MockMailRepoIMockRecorder) UpdateMessageState(userID, messageID, stateName, stateValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessageState", reflect.TypeOf((*MockMailRepoI)(nil).UpdateMessageState), userID, messageID, stateName, stateValue)
}
