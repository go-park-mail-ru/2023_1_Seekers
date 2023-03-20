package postgres

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	return config.DBSchemaName + ".box"
}

type Message struct {
	MessageID  uint64
	FromUserID uint64
	Size       int
	Title      string
	Text       string
	CreatedAt  string
	ReplyTo    uint64
}

func (Message) TableName() string {
	return config.DBSchemaName + ".messages"
}

func New(db *gorm.DB) mail.RepoI {
	return &mailRepository{
		db: db,
	}
}

func (m mailRepository) SelectFolderByUserNFolder(userID uint64, folderSlug string) (*models.Folder, error) {
	var folder models.Folder

	tx := m.db.Where("user_id = ? AND local_name = ?", userID, folderSlug).First(&folder)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &folder, nil
}

func (m mailRepository) SelectFoldersByUser(userID uint64) ([]models.Folder, error) {
	var folders []models.Folder

	tx := m.db.Where("user_id = ?", userID).Find(&folders)

	if tx.Error != nil {
		return folders, tx.Error
	}

	return folders, nil
}

func (m mailRepository) SelectFolderMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.MessageInfo, error) {
	var messages []models.MessageInfo
	tx := m.db.Model(Box{}).Select("*").Joins("JOIN "+Message{}.TableName()+" using(message_id)").
		Where("user_id = ? AND folder_id = ?", userID, folderID).Scan(&messages)

	if tx.Error != nil {
		return messages, tx.Error
	}

	return messages, nil
}

func (m mailRepository) SelectRecipientsByMessage(messageID uint64, fromUserID uint64) ([]uint64, error) {
	var recipientsIDs []uint64

	tx := m.db.Model(Box{}).Select("user_id").Where("message_id = ? AND user_id != ?", messageID, fromUserID).Scan(&recipientsIDs)

	if tx.Error != nil {
		return recipientsIDs, tx.Error
	}

	return recipientsIDs, nil
}

func (m mailRepository) CreateMessage(message models.MessageInfo, to ...uint64) error {
	//TODO implement me
	log.Warn("Now not creating hello meassage , NOT IMPLEMENTED")
	return nil
}
