package reporsitory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/models"
)

type RepositoryI interface {
	SelectIncomingMessagesByUser(userID uint64) ([]models.IncomingMessage, error)
	SelectOutgoingMessagesByUser(userID uint64) ([]models.OutgoingMessage, error)
	SelectFolderByUserNFolder(userID uint64, folderID uint64) (*models.Folder, error)
	SelectFoldersByUser(userID uint64) []models.Folder
	SelectMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.IncomingMessage, error)
}
