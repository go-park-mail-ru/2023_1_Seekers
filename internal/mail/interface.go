package mail

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"net/http"
)

type HandlersI interface {
	GetInboxMessages(w http.ResponseWriter, r *http.Request)
	GetOutboxMessages(w http.ResponseWriter, r *http.Request)
	GetFolderMessages(w http.ResponseWriter, r *http.Request)
}

type UseCaseI interface {
	GetIncomingMessages(userID uint64) ([]models.IncomingMessage, error)
	GetOutgoingMessages(userID uint64) ([]models.OutgoingMessage, error)
	GetFolders(userID uint64) []models.Folder
	GetFolderMessages(userID uint64, folderID uint64) ([]models.IncomingMessage, error)
}

type RepoI interface {
	SelectIncomingMessagesByUser(userID uint64) ([]models.IncomingMessage, error)
	SelectOutgoingMessagesByUser(userID uint64) ([]models.OutgoingMessage, error)
	SelectFolderByUserNFolder(userID uint64, folderID uint64) (*models.Folder, error)
	SelectFoldersByUser(userID uint64) []models.Folder
	SelectMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.IncomingMessage, error)
}
