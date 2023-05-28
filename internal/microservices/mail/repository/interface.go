package repository

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

//go:generate mockgen -destination=./mocks/mockrepo.go -source=./interface.go -package=mocks

type MailRepoI interface {
	SelectFolderByUserNFolderSlug(userID uint64, folderSlug string) (*models.Folder, error)
	SelectFolderByUserNFolderName(userID uint64, folderName string) (*models.Folder, error)
	SelectFoldersByUser(userID uint64) ([]models.Folder, error)
	SelectCustomFoldersByUser(userID uint64, defaultLocalNames []string) ([]models.Folder, error)
	//SelectFolderByUserNMessage(userID uint64, messageID uint64) (*models.Folder, error)
	CheckExistingBox(userIDs []uint64, messageID uint64, folderID uint64) (bool, error)
	SelectFolderMessagesByUserNFolderID(userIDs []uint64, folderID uint64, isDraft bool) ([]models.MessageInfo, error)
	SearchMessages(userIDs []uint64, folderID uint64, fromUser, toUser, filterText string, isDraft bool) ([]models.MessageInfo, error)
	SearchRecipients(userIDs []uint64) ([]models.UserInfo, error)
	DeleteFolder(folderID uint64) error
	DeleteMessageFromMessages(messageID uint64) error
	DeleteBox(userIDs []uint64, messageID uint64, folderID uint64) error
	DeleteBoxByUserNMessage(userID uint64, messageID uint64) error
	UpdateFolder(folder models.Folder) error
	SelectRecipientsByMessage(messageID uint64, fromUserID uint64) ([]uint64, error)
	SelectMessageByUserNMessage(userIDs []uint64, messageID uint64) (*models.MessageInfo, error)
	InsertMessage(fromUserID uint64, message *models.MessageInfo, user2folder []models.User2Folder) error
	UpdateMessage(message *models.MessageInfo, toInsert []models.User2Folder, toDelete []models.User2Folder) error
	InsertFolder(folder *models.Folder) (uint64, error)
	UpdateMessageState(userIDs []uint64, messageID uint64, folderID uint64, stateName string, stateValue bool) error
	UpdateMessageFolder(userIDs []uint64, messageID uint64, oldFolderID uint64, newFolderID uint64) error
	GetAttach(attachID, userID uint64) (*models.AttachmentInfo, error)
	GetMessageAttachments(messageID uint64) ([]models.AttachmentInfo, error)
	//UpdateMessageSender(messageID uint64, fromUserID uint64) error
	InsertFakeAccount(userID uint64, fakeUserID uint64) error
	SelectFakeIDs(userID uint64) ([]uint64, error)
	IsOwnerFakeAccount(userID uint64, fakeID uint64) error
	DeleteFakeAccount(userID uint64, fakeID uint64) error
	SelectOwnerFakeAccount(fakeID uint64) (uint64, error)
	SelectMessagesByFakeAccount(fakeID uint64, isDraft bool) ([]models.MessageInfo, error)
}
