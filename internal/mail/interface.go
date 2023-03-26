package mail

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"net/http"
)

type HandlersI interface {
	GetFolderMessages(w http.ResponseWriter, r *http.Request)
	GetFolders(w http.ResponseWriter, r *http.Request)
	GetMessage(w http.ResponseWriter, r *http.Request)
	SendMessage(w http.ResponseWriter, r *http.Request)
	ReadMessage(w http.ResponseWriter, r *http.Request)
	UnreadMessage(w http.ResponseWriter, r *http.Request)
}

type UseCaseI interface {
	GetFolders(userID uint64) ([]models.Folder, error)
	GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error)
	GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error)
	GetMessage(userID uint64, messageID uint64) (*models.MessageInfo, error)
	ValidateRecipients(recipients []string) ([]string, []string)
	SendMessage(userID uint64, message models.FormMessage) (*models.MessageInfo, error)
	SendFailedSendingMessage(recipientEmail string, invalidEmails []string) error
	SendWelcomeMessage(recipientEmail string) error
	MarkMessageAsSeen(userID uint64, messageID uint64) (*models.MessageInfo, error)
	MarkMessageAsUnseen(userID uint64, messageID uint64) (*models.MessageInfo, error)
}

type RepoI interface {
	SelectFolderByUserNFolder(userID uint64, folderSlug string) (*models.Folder, error)
	SelectFoldersByUser(userID uint64) ([]models.Folder, error)
	SelectFolderMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.MessageInfo, error)
	SelectRecipientsByMessage(messageID uint64, fromUserID uint64) ([]uint64, error)
	SelectMessageByUserNMessage(userID uint64, messageID uint64) (*models.MessageInfo, error)
	InsertMessageToMessages(message *models.MessageInfo) (uint64, error)
	InsertMessageToBoxes(userID uint64, folderID uint64, message *models.MessageInfo) error
	UpdateMessageState(userID uint64, messageID uint64, stateName string, stateValue bool) error
}
