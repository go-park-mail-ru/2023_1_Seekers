package mail

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"net/http"
)

type HandlersI interface {
	GetFolderMessages(w http.ResponseWriter, r *http.Request)
	GetFolders(w http.ResponseWriter, r *http.Request)
}

type UseCaseI interface {
	GetFolders(userID uint64) ([]models.Folder, error)
	GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error)
	GetFolderMessages(userID uint64, folderID string) ([]models.MessageInfo, error)
	CreateHelloMessage(to uint64) error
	CreateMessage(message models.MessageInfo, to ...uint64) error
}

type RepoI interface {
	SelectFolderByUserNFolder(userID uint64, folderID string) (*models.Folder, error)
	SelectFoldersByUser(userID uint64) ([]models.Folder, error)
	SelectFolderMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.MessageInfo, error)
	SelectRecipientsByMessage(messageID uint64, fromUserID uint64) ([]uint64, error)
	CreateMessage(message models.MessageInfo, to ...uint64) error
}
