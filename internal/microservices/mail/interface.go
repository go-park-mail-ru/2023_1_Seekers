package mail

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

type UseCaseI interface {
	GetFolders(userID uint64) ([]models.Folder, error)
	GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error)
	GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error)
	CreateDefaultFolders(userID uint64) ([]models.Folder, error)
	GetMessage(userID uint64, messageID uint64) (*models.MessageInfo, error)
	ValidateRecipients(recipients []string) ([]string, []string)
	SendMessage(userID uint64, message models.FormMessage) (*models.MessageInfo, error)
	SendFailedSendingMessage(recipientEmail string, invalidEmails []string) error
	SendWelcomeMessage(recipientEmail string) error
	MarkMessageAsSeen(userID uint64, messageID uint64) (*models.MessageInfo, error)
	MarkMessageAsUnseen(userID uint64, messageID uint64) (*models.MessageInfo, error)
}
