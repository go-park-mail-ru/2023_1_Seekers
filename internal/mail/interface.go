package mail

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"net/http"
)

type HandlersI interface {
	GetInboxMessages(w http.ResponseWriter, r *http.Request)
	GetOutboxMessages(w http.ResponseWriter, r *http.Request)
	GetFolderMessages(w http.ResponseWriter, r *http.Request)
	GetFolders(w http.ResponseWriter, r *http.Request)
}

type UseCaseI interface {
	GetIncomingMessages(userID uint64) ([]models.IncomingMessage, error)
	GetOutgoingMessages(userID uint64) ([]models.OutgoingMessage, error)
	GetFolders(userID uint64) []models.Folder
	GetFolderInfo(userID uint64, folderID uint64) (*models.Folder, error)
	GetFolderMessages(userID uint64, folderID uint64) ([]models.IncomingMessage, error)
	CreateHelloMessage(to uint64) error
	CreateMessage(message models.Message, to ...uint64) error
}

type RepoI interface {
	SelectIncomingMessagesByUser(userID uint64) ([]models.IncomingMessage, error)
	SelectOutgoingMessagesByUser(userID uint64) ([]models.OutgoingMessage, error)
	SelectFolderByUserNFolder(userID uint64, folderID uint64) (*models.Folder, error)
	SelectFoldersByUser(userID uint64) []models.Folder
	SelectMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.IncomingMessage, error)
	CreateMessage(message models.Message, to ...uint64) error
}
