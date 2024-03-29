package postgres

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/repository"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type mailRepository struct {
	cfg *config.Config
	db  *gorm.DB
}

type Box struct {
	UserID    uint64
	MessageID uint64
	FolderID  uint64
	Seen      bool
	Favorite  bool
	Deleted   bool
	IsDraft   bool
}

func (Box) TableName(schemaName string) string {
	return schemaName + ".boxes"
}

type Message struct {
	MessageID        uint64 `gorm:"primaryKey"`
	FromUserID       uint64
	Title            string
	Text             string
	CreatedAt        string
	ReplyToMessageID *uint64
}

func (Message) TableName(schemaName string) string {
	return schemaName + ".messages"
}

func New(c *config.Config, db *gorm.DB) repository.MailRepoI {
	return &mailRepository{
		cfg: c,
		db:  db,
	}
}

func (m mailRepository) SelectFolderByUserNFolderSlug(userID uint64, folderSlug string) (*models.Folder, error) {
	var folder models.Folder

	tx := m.db.Where("user_id = ? AND local_name = ?", userID, folderSlug).First(&folder)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.WithMessage(errors.ErrFolderNotFound, "select folder by local_name")
		}

		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return &folder, nil
}

func (m mailRepository) SelectFolderByUserNFolderName(userID uint64, folderName string) (*models.Folder, error) {
	var folder models.Folder

	tx := m.db.Where("user_id = ? AND name = ?", userID, folderName).First(&folder)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.WithMessage(errors.ErrFolderNotFound, "select folder by name")
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

func (m mailRepository) SearchRecipients(userId uint64) ([]models.UserInfo, error) {
	var result []models.UserInfo
	tx := m.db.Raw("SELECT * FROM get_recipes( $1 );", userId).Scan(&result)

	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return result, nil
}

func (m mailRepository) SelectFolderByUserNMessage(userID uint64, messageID uint64) (*models.Folder, error) {
	var folder models.Folder

	tx := m.db.Model(Box{}).Select("folders.*").Joins("JOIN mail.folders using(folder_id)").
		Where("boxes.user_id = ? AND message_id = ?", userID, messageID).First(&folder)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return &folder, nil
}

func (m mailRepository) CheckExistingBox(userID uint64, messageID uint64, folderID uint64) (bool, error) {
	var box Box

	tx := m.db.Where("user_id = ? AND message_id = ? AND folder_id = ?", userID, messageID, folderID).First(&box)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return true, nil
}

func (m mailRepository) SelectFolderMessagesByUserNFolderID(userID uint64, folderID uint64, isDraft bool) ([]models.MessageInfo, error) {
	var messages []models.MessageInfo

	tx := m.db.Model(Box{}).Select("*").Joins("JOIN "+Message{}.TableName(m.cfg.DB.DBSchemaName)+" using(message_id)").
		Where("user_id = ? AND folder_id = ? AND is_draft = ?", userID, folderID, isDraft).Order("created_at DESC").Scan(&messages)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return messages, nil
}

func (m mailRepository) SearchMessages(userId uint64, fromUser, toUser, folder, filter string) ([]models.MessageInfo, error) {
	var messages []models.MessageInfo
	var messagesIds []uint64

	tx := m.db.Raw("SELECT * FROM get_messages($1, $2, $3, $4, $5);", userId, fromUser, toUser, folder, filter).Scan(&messagesIds)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	for _, mid := range messagesIds {
		mInfo, err := m.SelectMessageByUserNMessage(userId, mid)
		if err != nil {
			return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
		}
		messages = append(messages, *mInfo)
	}

	return messages, nil
}

func (m mailRepository) DeleteFolder(folderID uint64) error {
	tx := m.db.Where("folder_id = ?", folderID).Delete(&models.Folder{})
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (m mailRepository) DeleteBox(userID uint64, messageID uint64, folderID uint64) error {
	tx := m.db.Where("user_id = ? AND message_id = ? AND folder_id = ?", userID, messageID, folderID).Delete(&Box{})
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (m mailRepository) DeleteMessageFromMessages(messageID uint64) error {
	tx := m.db.Where("message_id = ?", messageID).Delete(Message{})
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}
	return nil
}

func (m mailRepository) UpdateFolder(folder models.Folder) error {
	tx := m.db.Updates(folder)
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (m mailRepository) SelectRecipientsByMessage(messageID uint64, fromUserID uint64) ([]uint64, error) {
	var recipientsIDs []uint64

	tx := m.db.Model(Box{}).Select("user_id").Where("message_id = ?", messageID).
		Scan(&recipientsIDs)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	for i, id := range recipientsIDs {
		if id == fromUserID {
			recipientsIDs = append(recipientsIDs[:i], recipientsIDs[i+1:]...)

			return recipientsIDs, nil
		}
	}

	return recipientsIDs, nil
}

func (m mailRepository) SelectMessageByUserNMessage(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	var message *models.MessageInfo

	tx := m.db.Model(Box{}).Select("*").Joins("JOIN "+Message{}.TableName(m.cfg.DB.DBSchemaName)+" using(message_id)").
		Where("user_id = ? AND message_id = ? AND (from_user_id = user_id OR from_user_id != user_id AND is_draft = false)", userID, messageID).
		Scan(&message)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	message.Preview = common.GetInnerText(message.Text, m.cfg.Api.MailPreviewMaxLen)
	return message, nil
}

func (m mailRepository) insertMessageToMessages(fromUserID uint64, message *models.MessageInfo, tx *gorm.DB) (uint64, error) {
	convMsg := convertToMessageDB(fromUserID, message)

	tx = tx.Select("from_user_id", "title", "text", "created_at", "reply_to_message_id").Create(&convMsg)
	if err := tx.Error; err != nil {
		return 0, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	for _, v := range message.Attachments {
		var attachID uint64
		tx = tx.Raw("INSERT INTO mail.attaches(message_id, type, filename, s3_fname, size_str, size_count) VALUES ($1, $2, $3, $4, $5, $6) RETURNING attach_id;", convMsg.MessageID, v.Type, v.FileName, v.S3FName, v.SizeStr, v.SizeCount).Scan(&attachID)
		if err := tx.Error; err != nil {
			return 0, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
		}
		v.AttachID = attachID
	}

	return convMsg.MessageID, nil
}

func (m mailRepository) updateMessageInMessages(message *models.MessageInfo, tx *gorm.DB) error {
	tx = tx.Model(Message{}).Omit("message_id", "from_user_id", "size").Where("message_id = ?", message.MessageID).Updates(&message)
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}
	if message.ReplyToMessageID == nil {
		tx = tx.Model(Message{}).Select("reply_to_message_id").Where("message_id = ?", message.MessageID).
			Updates(&map[string]interface{}{"reply_to_message_id": nil})
	}

	return nil
}

func convertToMessageDB(fromUserID uint64, message *models.MessageInfo) Message {
	return Message{
		FromUserID:       fromUserID,
		Title:            message.Title,
		Text:             message.Text,
		CreatedAt:        message.CreatedAt,
		ReplyToMessageID: message.ReplyToMessageID,
	}
}

func (m mailRepository) insertMessageToBoxes(userID uint64, folderID uint64, message *models.MessageInfo, tx *gorm.DB) error {
	convMsg := convertToBoxDB(userID, folderID, message)

	tx = tx.Create(&convMsg)
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, tx.Error.Error())
	}

	return nil
}

func (m mailRepository) deleteMessageFromBoxes(userID uint64, folderID uint64, messageID uint64, tx *gorm.DB) error {
	tx = tx.Where("user_id = ? AND folder_id = ? AND message_id = ?", userID, folderID, messageID).Delete(Box{})
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, tx.Error.Error())
	}

	return nil
}

func (m mailRepository) InsertMessage(fromUserID uint64, message *models.MessageInfo, user2folder []models.User2Folder) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		messageID, err := m.insertMessageToMessages(fromUserID, message, tx)
		if err != nil {
			return pkgErrors.Wrap(err, "insert message : insert to messages")
		}

		message.MessageID = messageID
		message.Seen = true
		for _, elem := range user2folder {
			err = m.insertMessageToBoxes(elem.UserID, elem.FolderID, message, tx)
			if err != nil {
				return pkgErrors.Wrap(err, "insert message : insert to boxes")
			}

			message.Seen = false
		}

		return nil
	})
}

func (m mailRepository) SelectCustomFoldersByUser(userID uint64, defaultLocalNames []string) ([]models.Folder, error) {
	var folders []models.Folder

	tx := m.db.Where("user_id = ? AND local_name NOT IN ?", userID, defaultLocalNames).Order("name").Find(&folders)
	if err := tx.Error; err != nil {
		return make([]models.Folder, 0), pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return folders, nil
}

func (m mailRepository) UpdateMessage(message *models.MessageInfo, toInsert []models.User2Folder, toDelete []models.User2Folder) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		err := m.updateMessageInMessages(message, tx)
		if err != nil {
			return pkgErrors.Wrap(err, "update message : update in messages")
		}

		for _, elem := range toInsert {
			err := m.insertMessageToBoxes(elem.UserID, elem.FolderID, message, tx)
			if err != nil {
				return pkgErrors.Wrap(err, "update message : insert to boxes")
			}
		}

		for _, elem := range toDelete {
			err := m.deleteMessageFromBoxes(elem.UserID, elem.FolderID, message.MessageID, tx)
			if err != nil {
				return pkgErrors.Wrap(err, "update message : delete from boxes")
			}
		}

		return nil
	})
}

func convertToBoxDB(userID uint64, folderID uint64, message *models.MessageInfo) Box {
	return Box{
		UserID:    userID,
		MessageID: message.MessageID,
		FolderID:  folderID,
		Seen:      message.Seen,
		Favorite:  message.Favorite,
		Deleted:   message.Deleted,
		IsDraft:   message.IsDraft,
	}
}

func (m mailRepository) InsertFolder(folder *models.Folder) (uint64, error) {
	tx := m.db.Create(&folder)

	if err := tx.Error; err != nil {
		return 0, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return folder.FolderID, nil
}

func (m mailRepository) UpdateMessageState(userID uint64, messageID uint64, folderID uint64, stateName string, stateValue bool) error {
	tx := m.db.Model(Box{}).Where("user_id = ? AND message_id = ? AND folder_id = ?", userID, messageID, folderID).Update(stateName, stateValue)
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (m mailRepository) UpdateMessageFolder(userID uint64, messageID uint64, oldFolderID uint64, newFolderID uint64) error {
	tx := m.db.Model(Box{}).Where("user_id = ? AND message_id = ? AND folder_id = ?", userID, messageID, oldFolderID).Update("folder_id", newFolderID)
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (m mailRepository) GetAttach(attachID, userID uint64) (*models.AttachmentInfo, error) {
	type Result struct {
		Type      string
		Filename  string
		S3FName   string `gorm:"column:s3_fname"`
		SizeStr   string
		SizeCount int64
	}
	var res Result

	tx := m.db.Raw("SELECT type, filename, s3_fname, size_str, size_count from mail.attaches "+
		"JOIN mail.boxes b on attaches.message_id = b.message_id "+
		"WHERE attach_id = $1 AND user_id = $2;", attachID, userID).Scan(&res)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return &models.AttachmentInfo{
		AttachID:  attachID,
		FileName:  res.Filename,
		S3FName:   res.S3FName,
		Type:      res.Type,
		SizeStr:   res.SizeStr,
		SizeCount: res.SizeCount,
	}, nil
}

func (m mailRepository) GetMessageAttachments(messageID uint64) ([]models.AttachmentInfo, error) {
	type Result struct {
		AttachID  uint64
		Type      string
		Filename  string
		S3FName   string `gorm:"column:s3_fname"`
		SizeStr   string
		SizeCount int64
	}
	var res []Result

	tx := m.db.Raw("SELECT attach_id, type, filename, s3_fname, size_str, size_count from mail.attaches "+
		"JOIN mail.messages m on attaches.message_id = m.message_id "+
		"WHERE m.message_id = $1;", messageID).Scan(&res)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	response := make([]models.AttachmentInfo, len(res))
	for i, v := range res {
		response[i] = models.AttachmentInfo{
			AttachID:  v.AttachID,
			FileName:  v.Filename,
			S3FName:   v.S3FName,
			Type:      v.Type,
			SizeStr:   v.SizeStr,
			SizeCount: v.SizeCount,
		}
	}
	return response, nil
}
