package postgres

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
	"os"
)

type mailRepository struct {
	db *gorm.DB
}

type Box struct {
	UserID    uint64
	MessageID uint64
	FolderID  uint64
	Seen      bool
	Favorite  bool
	Deleted   bool
}

func (Box) TableName() string {
	return os.Getenv(config.DBSchemaNameEnv) + ".boxes"
}

type Message struct {
	MessageID        uint64 `gorm:"primaryKey"`
	FromUserID       uint64
	Title            string
	Text             string
	CreatedAt        string
	ReplyToMessageID *uint64
}

func (Message) TableName() string {
	return os.Getenv(config.DBSchemaNameEnv) + ".messages"
}

func New(db *gorm.DB) mail.RepoI {
	return &mailRepository{
		db: db,
	}
}

func (m mailRepository) SelectFolderByUserNFolder(userID uint64, folderSlug string) (*models.Folder, error) {
	var folder models.Folder

	tx := m.db.Where("user_id = ? AND local_name = ?", userID, folderSlug).First(&folder)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.WithMessage(errors.ErrFolderNotFound, err.Error())
		}
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return &folder, nil
}

func (m mailRepository) SelectFoldersByUser(userID uint64) ([]models.Folder, error) {
	var folders []models.Folder

	tx := m.db.Where("user_id = ?", userID).Find(&folders)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return folders, nil
}

func (m mailRepository) SelectFolderMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.MessageInfo, error) {
	var messages []models.MessageInfo

	tx := m.db.Model(Box{}).Select("*").Joins("JOIN "+Message{}.TableName()+" using(message_id)").
		Where("user_id = ? AND folder_id = ?", userID, folderID).Scan(&messages)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.WithMessage(errors.ErrMessageNotFound, err.Error())
		}
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return messages, nil
}

func (m mailRepository) SelectRecipientsByMessage(messageID uint64, fromUserID uint64) ([]uint64, error) {
	var recipientsIDs []uint64

	tx := m.db.Model(Box{}).Select("user_id").Where("message_id = ? AND user_id != ?", messageID, fromUserID).
		Scan(&recipientsIDs)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return recipientsIDs, nil
}

func (m mailRepository) SelectMessageByUserNMessage(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	var message *models.MessageInfo

	tx := m.db.Model(Box{}).Select("*").Joins("JOIN "+Message{}.TableName()+" using(message_id)").
		Where("user_id = ? AND message_id = ?", userID, messageID).Scan(&message)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.WithMessage(errors.ErrMessageNotFound, err.Error())
		}
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return message, nil
}

func (m mailRepository) InsertMessageToMessages(message *models.MessageInfo) (uint64, error) {
	convMsg := convertToMessageDB(message)

	tx := m.db.Select("from_user_id", "title", "text", "created_at", "reply_to_message_id").Create(&convMsg)
	if err := tx.Error; err != nil {
		return 0, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return convMsg.MessageID, nil
}

func convertToMessageDB(message *models.MessageInfo) Message {
	return Message{
		FromUserID:       message.FromUser.UserID,
		Title:            message.Title,
		Text:             message.Text,
		CreatedAt:        message.CreatedAt,
		ReplyToMessageID: message.ReplyToMessageID,
	}
}

func (m mailRepository) InsertMessageToBoxes(userID uint64, folderID uint64, message *models.MessageInfo) error {
	convMsg := convertToBoxDB(userID, folderID, message)
	tx := m.db.Create(&convMsg)

	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func convertToBoxDB(userID uint64, folderID uint64, message *models.MessageInfo) Box {
	return Box{
		UserID:    userID,
		MessageID: message.MessageID,
		FolderID:  folderID,
		Seen:      message.Seen,
		Favorite:  message.Favorite,
		Deleted:   message.Deleted,
	}
}

func (m mailRepository) InsertFolder(folder *models.Folder) (uint64, error) {
	tx := m.db.Create(&folder)

	if err := tx.Error; err != nil {
		return 0, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return folder.FolderID, nil
}

func (m mailRepository) UpdateMessageState(userID uint64, messageID uint64, stateName string, stateValue bool) error {
	tx := m.db.Model(Box{}).Where("user_id = ? AND message_id = ?", userID, messageID).Update(stateName, stateValue)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.WithMessage(errors.ErrMessageNotFound, err.Error())
		}
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}
	return nil
}