package repository

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

//go:generate mockgen -destination=./mocks/mockrepo.go -source=./interface.go -package=mocks

type MailRepoI interface {
	SelectFolderByUserNFolder(userID uint64, folderSlug string) (*models.Folder, error)
	SelectFoldersByUser(userID uint64) ([]models.Folder, error)
	SelectFolderMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.MessageInfo, error)
	SelectRecipientsByMessage(messageID uint64, fromUserID uint64) ([]uint64, error)
	SelectMessageByUserNMessage(userID uint64, messageID uint64) (*models.MessageInfo, error)
	InsertMessage(fromUserID uint64, message *models.MessageInfo, user2folder []models.User2Folder) error
	InsertFolder(folder *models.Folder) (uint64, error)
	UpdateMessageState(userID uint64, messageID uint64, stateName string, stateValue bool) error
}
