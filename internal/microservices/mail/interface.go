package mail

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

type UseCaseI interface {
	GetFolders(userID uint64) ([]models.Folder, error)
	GetCustomFolders(userID uint64) ([]models.Folder, error)
	GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error)
	GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error)
	CreateDefaultFolders(userID uint64) ([]models.Folder, error)
	CreateFolder(userID uint64, form models.FormFolder) (*models.Folder, error)
	DeleteFolder(userID uint64, folderSlug string) error
	EditFolder(userID uint64, folderSlug string, form models.FormFolder) (*models.Folder, error)
	GetMessage(userID uint64, messageID uint64) (*models.MessageInfo, error)
	DeleteMessage(userID uint64, messageID uint64) error
	ValidateRecipients(recipients []string) ([]string, []string)
	SendMessage(userID uint64, message models.FormMessage) (*models.MessageInfo, error)
	SaveDraft(userID uint64, message models.FormMessage) (*models.MessageInfo, error)
	EditDraft(userID uint64, messageID uint64, message models.FormMessage) (*models.MessageInfo, error)
	SendFailedSendingMessage(recipientEmail string, invalidEmails []string) error
	SendWelcomeMessage(recipientEmail string) error
	MarkMessageAsSeen(userID uint64, messageID uint64) (*models.MessageInfo, error)
	MarkMessageAsUnseen(userID uint64, messageID uint64) (*models.MessageInfo, error)
	MoveMessageToFolder(userID uint64, messageID uint64, folderSlug string) error
}
