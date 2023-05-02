package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/repository"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/smtp/client"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgSmtp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/smtp"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/validation"
	pkgErrors "github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
)

//go:generate mockgen -destination=./mocks/mockusecase.go -package=mocks -source=../interface.go

type mailUC struct {
	cfg      *config.Config
	mailRepo mailRepo.MailRepoI
	userUC   user.UseCaseI
}

func New(c *config.Config, mR mailRepo.MailRepoI, uUC user.UseCaseI) mail.UseCaseI {
	return &mailUC{
		cfg:      c,
		mailRepo: mR,
		userUC:   uUC,
	}
}

var defaultFolderNames = map[string]string{
	"inbox":  "Входящие",
	"outbox": "Исходящие",
	"trash":  "Корзина",
	"drafts": "Черновики",
	"spam":   "Спам",
}

func (uc *mailUC) GetFolders(userID uint64) ([]models.Folder, error) {
	folders, err := uc.mailRepo.SelectFoldersByUser(userID)
	if err != nil {
		return []models.Folder{}, pkgErrors.Wrap(err, "get folders")
	}

	return folders, nil
}

func (uc *mailUC) GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error) {
	folder, err := uc.mailRepo.SelectFolderByUserNFolderSlug(userID, folderSlug)
	if err != nil {
		return folder, pkgErrors.Wrap(err, "get folder info")
	}

	return folder, nil
}

func (uc *mailUC) GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error) {
	var messages []models.MessageInfo

	folder, err := uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages")
	}

	messages, err = uc.mailRepo.SelectFolderMessagesByUserNFolderID(userID, folder.FolderID)
	if err != nil {
		return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages : msg by user and folder")
	}

	for i, message := range messages {
		messageID := message.MessageID

		fromUser, err := uc.userUC.GetInfo(message.FromUser.UserID)
		if err != nil {
			return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages : get info by id")
		}

		messages[i].FromUser = *fromUser
		recipientsIDs, err := uc.mailRepo.SelectRecipientsByMessage(messageID, message.FromUser.UserID)
		if err != nil {
			return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages : get recipients by msg")
		}

		for _, recipientsID := range recipientsIDs {
			profile, err := uc.userUC.GetInfo(recipientsID)
			if err != nil {
				return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages : get info by id")
			}

			messages[i].Recipients = append(message.Recipients, *profile)
		}
	}

	return messages, nil
}

func (uc *mailUC) CreateDefaultFolders(userID uint64) ([]models.Folder, error) {
	for key, value := range defaultFolderNames {
		currentFolder := models.Folder{
			UserID:    userID,
			LocalName: key,
			Name:      value,
		}

		_, err := uc.mailRepo.InsertFolder(&currentFolder)
		if err != nil {
			return []models.Folder{}, pkgErrors.Wrap(err, "create default folders")
		}
	}

	return uc.GetFolders(userID)
}

func (uc *mailUC) getLastLocalFolderID(userID uint64) (uint64, error) {
	folders, err := uc.GetFolders(userID)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "get folders")
	}

	var lastLocalID uint64
	var IDs []uint64

	for _, folder := range folders {
		curID, err := strconv.ParseUint(folder.LocalName, 10, 64)
		if err == nil {
			IDs = append(IDs, curID)
		}
	}

	sort.Slice(IDs, func(i, j int) bool {
		return IDs[i] < IDs[j]
	})

	for _, id := range IDs {
		lastLocalID++

		if id != lastLocalID {
			lastLocalID--
			break
		}
	}

	return lastLocalID, nil
}

func (uc *mailUC) CreateFolder(userID uint64, form models.FormFolder) (*models.Folder, error) {
	if err := validation.FolderName(form.Name); err != nil {
		return nil, pkgErrors.Wrap(err, "validate folder name")
	}

	folder, err := uc.mailRepo.SelectFolderByUserNFolderName(userID, form.Name)
	if !pkgErrors.Is(err, errors.ErrFolderNotFound) {
		return nil, pkgErrors.WithMessage(errors.ErrFolderAlreadyExists, "select folder by user and name")
	}

	lastLocalID, err := uc.getLastLocalFolderID(userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get last local folder id")
	}

	folder = &models.Folder{
		UserID:    userID,
		LocalName: strconv.FormatUint(lastLocalID+1, 10),
		Name:      form.Name,
	}

	_, err = uc.mailRepo.InsertFolder(folder)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "insert folder")
	}

	return folder, nil
}

func (uc *mailUC) isDefaultFolder(folderSlug string) bool {
	for key := range defaultFolderNames {
		if folderSlug == key {
			return true
		}
	}

	return false
}

func (uc *mailUC) DeleteFolder(userID uint64, folderSlug string) error {
	if uc.isDefaultFolder(folderSlug) {
		return pkgErrors.WithMessage(errors.ErrDeleteDefaultFolder, "is default folder")
	}

	folder, err := uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		return pkgErrors.Wrap(err, "get folder info")
	}

	messages, err := uc.mailRepo.SelectFolderMessagesByUserNFolderID(userID, folder.FolderID)
	if err != nil {
		return pkgErrors.Wrap(err, "get folder messages : msg by user and folder")
	}

	for _, message := range messages {
		err = uc.MoveMessageToFolder(userID, message.MessageID, "trash")
		if err != nil {
			return pkgErrors.Wrap(err, "move message to trash folder")
		}
	}

	err = uc.mailRepo.DeleteFolder(folder.FolderID)
	if err != nil {
		return pkgErrors.Wrap(err, "get delete folder")
	}

	return nil
}

func (uc *mailUC) EditFolder(userID uint64, folderSlug string, form models.FormFolder) (*models.Folder, error) {
	if err := validation.FolderName(form.Name); err != nil {
		return nil, pkgErrors.Wrap(err, "validate folder name")
	}

	if uc.isDefaultFolder(folderSlug) {
		return nil, pkgErrors.WithMessage(errors.ErrEditDefaultFolder, "is default folder")
	}

	folder, err := uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get folder info")
	}

	_, err = uc.mailRepo.SelectFolderByUserNFolderName(userID, form.Name)
	if !pkgErrors.Is(err, errors.ErrFolderNotFound) {
		return nil, pkgErrors.WithMessage(errors.ErrFolderAlreadyExists, "select folder by user and name")
	}

	folder.Name = form.Name

	err = uc.mailRepo.UpdateFolder(*folder)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "update folder")
	}

	return uc.GetFolderInfo(userID, folderSlug)
}

func (uc *mailUC) GetMessage(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	var firstMessage *models.MessageInfo
	var prevMessage *models.MessageInfo
	replyToMsgID := &messageID

	for replyToMsgID != nil {
		curMessage, err := uc.mailRepo.SelectMessageByUserNMessage(userID, *replyToMsgID)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "get message : by Uid and Mid")
		}
		if curMessage == nil {
			return nil, pkgErrors.WithMessage(errors.ErrMessageNotFound, "get message")
		}

		fromUser, err := uc.userUC.GetInfo(curMessage.FromUser.UserID)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "get message : get info by Uid")
		}

		curMessage.FromUser = *fromUser
		recipientsIDs, err := uc.mailRepo.SelectRecipientsByMessage(*replyToMsgID, curMessage.FromUser.UserID)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "get message : get recipients by Mid")
		}

		for _, recipientsID := range recipientsIDs {
			profile, err := uc.userUC.GetInfo(recipientsID)
			if err != nil {
				return nil, pkgErrors.Wrap(err, "get message : get recipient info by Uid")
			}

			curMessage.Recipients = append(curMessage.Recipients, *profile)
		}

		if *replyToMsgID == messageID {
			firstMessage = curMessage
		} else {
			prevMessage.ReplyTo = curMessage
		}

		replyToMsgID = curMessage.ReplyToMessageID
		prevMessage = curMessage
	}

	return firstMessage, nil
}

func (uc *mailUC) DeleteMessage(userID uint64, messageID uint64) error {
	curMessage, err := uc.mailRepo.SelectMessageByUserNMessage(userID, messageID)
	if err != nil {
		return pkgErrors.Wrap(err, "get message : by Uid and Mid")
	}
	if curMessage == nil {
		return pkgErrors.WithMessage(errors.ErrMessageNotFound, "get message")
	}

	fromFolder, err := uc.getFolderByMessage(userID, messageID)
	if err != nil {
		return pkgErrors.Wrap(err, "get folder by message")
	}

	switch fromFolder.LocalName {
	case "trash":
		err = uc.mailRepo.DeleteMessageForUser(userID, messageID)
		if err != nil {
			return pkgErrors.Wrap(err, "delete message for user")
		}
		break
	case "drafts":
		err = uc.mailRepo.DeleteMessageFromMessages(messageID)
		if err != nil {
			return pkgErrors.Wrap(err, "delete message full")
		}

		break
	default:
		err = uc.MoveMessageToFolder(userID, messageID, "trash")
		if err != nil {
			return pkgErrors.Wrap(err, "move message to trash folder")
		}
	}

	return nil
}

func (uc *mailUC) ValidateRecipients(recipients []string) ([]string, []string) {
	var validEmails []string
	var invalidEmails []string

	for _, email := range recipients {
		if strings.Contains(email, uc.cfg.Mail.PostDomain) {
			_, err := uc.userUC.GetInfoByEmail(email)
			if err != nil {
				invalidEmails = append(invalidEmails, email)
			} else {
				validEmails = append(validEmails, email)
			}
		} else {
			if err := validation.ValidateEmail(email); err != nil {
				invalidEmails = append(invalidEmails, email)
			} else {
				validEmails = append(validEmails, email)
			}
		}
	}

	return validEmails, invalidEmails
}

func (uc *mailUC) SaveDraft(fromUserID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	folder, err := uc.GetFolderInfo(fromUserID, "drafts")
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message : get folder by UId and FolderSlug")
	}

	var user2folder []models.User2Folder
	user2folder = append(user2folder, models.User2Folder{
		UserID:   fromUserID,
		FolderID: folder.FolderID,
	})

	for _, email := range message.Recipients {
		recipient, err := uc.userUC.GetInfoByEmail(email)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "send message : get user info by email")
		}

		folder, err = uc.GetFolderInfo(recipient.UserID, "inbox")
		if err != nil {
			return nil, pkgErrors.Wrap(err, "send message : get folder by UId and FolderSlug")
		}

		user2folder = append(user2folder, models.User2Folder{
			UserID:   recipient.UserID,
			FolderID: folder.FolderID,
		})
	}

	newMessage := models.MessageInfo{
		Title:            message.Title,
		CreatedAt:        common.GetCurrentTime(uc.cfg.Logger.LogsTimeFormat),
		Text:             message.Text,
		ReplyToMessageID: message.ReplyToMessageID,
		IsDraft:          true,
	}

	err = uc.mailRepo.InsertMessage(fromUserID, &newMessage, user2folder)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message : insert message")
	}

	return uc.GetMessage(fromUserID, newMessage.MessageID)
}

func (uc *mailUC) mapRecipients(newRecipients []string, oldMessage *models.MessageInfo) map[string]string {
	recipients := make(map[string]string)

	for _, recipient := range newRecipients {
		recipients[recipient] = "add"
	}

	for _, recipient := range oldMessage.Recipients {
		if _, ok := recipients[recipient.Email]; ok {
			recipients[recipient.Email] = "save"
		} else {
			recipients[recipient.Email] = "del"
		}
	}

	return recipients
}

func (uc *mailUC) EditDraft(fromUserID uint64, messageID uint64, formMessage models.FormMessage) (*models.MessageInfo, error) {
	message, err := uc.GetMessage(fromUserID, messageID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get message")
	}

	recipients := uc.mapRecipients(formMessage.Recipients, message)
	var toInsert []models.User2Folder
	var toDelete []models.User2Folder

	for email, value := range recipients {
		switch value {
		case "add":
			recipient, err := uc.userUC.GetInfoByEmail(email)
			if err != nil {
				return nil, pkgErrors.Wrap(err, "send message : get user info by email")
			}

			folder, err := uc.GetFolderInfo(recipient.UserID, "inbox")
			if err != nil {
				return nil, pkgErrors.Wrap(err, "send message : get folder by UId and FolderSlug")
			}

			toInsert = append(toInsert, models.User2Folder{
				UserID:   recipient.UserID,
				FolderID: folder.FolderID,
			})
			break
		case "del":
			recipient, err := uc.userUC.GetInfoByEmail(email)
			if err != nil {
				return nil, pkgErrors.Wrap(err, "edit draft : get user info by email")
			}

			folder, err := uc.GetFolderInfo(recipient.UserID, "inbox")
			if err != nil {
				return nil, pkgErrors.Wrap(err, "send message : get folder by UId and FolderSlug")
			}

			toDelete = append(toDelete, models.User2Folder{
				UserID:   recipient.UserID,
				FolderID: folder.FolderID,
			})
			break
		}
	}

	message.Title = formMessage.Title
	message.Text = formMessage.Text
	message.ReplyToMessageID = formMessage.ReplyToMessageID
	message.CreatedAt = common.GetCurrentTime(uc.cfg.Logger.LogsTimeFormat)

	if err := uc.mailRepo.UpdateMessage(message, toInsert, toDelete); err != nil {
		return nil, pkgErrors.Wrap(err, "edit draft : update message")
	}

	return uc.GetMessage(fromUserID, messageID)
}

func (uc *mailUC) SendMessage(fromUserID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	if len(message.Recipients) == 0 {
		return nil, pkgErrors.WithMessage(errors.ErrNoValidEmails, "send message")
	}

	var user2folder []models.User2Folder

	fromUser, err := uc.userUC.GetByID(fromUserID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message")
	}

	if !fromUser.IsExternal {
		folder, err := uc.GetFolderInfo(fromUserID, "outbox")
		if err != nil {
			return nil, pkgErrors.Wrap(err, "send message : get folder by UId and FolderSlug")
		}

		user2folder = append(user2folder, models.User2Folder{
			UserID:   fromUserID,
			FolderID: folder.FolderID,
		})
	}

	for _, recipient := range message.Recipients {
		toDomain, err := pkgSmtp.ParseDomain(recipient)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "send message - failed get recipient domain")
		}
		if toDomain != uc.cfg.Mail.PostDomain {
			err := client.SendMail(fromUser, recipient, message.Title, message.Text, uc.cfg.Mail.PostDomain, uc.cfg.SmtpServer.SecretPassword)
			if err != nil {
				return nil, pkgErrors.Wrap(err, "send message : to other mail service")
			}
		} else {
			recipient, err := uc.userUC.GetInfoByEmail(recipient)
			if err != nil {
				return nil, pkgErrors.Wrap(err, "send message : get user info by email")
			}

			folder, err := uc.GetFolderInfo(recipient.UserID, "inbox")
			if err != nil {
				return nil, pkgErrors.Wrap(err, "send message : get folder by UId and FolderSlug")
			}

			user2folder = append(user2folder, models.User2Folder{
				UserID:   recipient.UserID,
				FolderID: folder.FolderID,
			})
		}
	}

	newMessage := models.MessageInfo{
		Title:            message.Title,
		CreatedAt:        common.GetCurrentTime(uc.cfg.Logger.LogsTimeFormat),
		Text:             message.Text,
		ReplyToMessageID: message.ReplyToMessageID,
	}

	err = uc.mailRepo.InsertMessage(fromUserID, &newMessage, user2folder)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message : insert message")
	}

	return uc.GetMessage(fromUserID, newMessage.MessageID)
}

func (uc *mailUC) SendFailedSendingMessage(recipientEmail string, invalidEmails []string) error {
	formMessage := models.FormMessage{
		Recipients: []string{recipientEmail},
		Title:      "Ваше сообщение не доставлено",
		Text: "Это письмо создано автоматически сервером Mailbx.ru, отвечать на него не нужно.\n\n" +
			"К сожалению, Ваше письмо не может быть доставлено одному или нескольким получателям:\n" +
			strings.Join(invalidEmails, "\n") + "\n\nРекомендуем Вам проверить корректность указания адресов получателей.",
		ReplyToMessageID: nil,
	}

	return uc.sendMessageFromSupport(formMessage)
}

func (uc *mailUC) SendWelcomeMessage(recipientEmail string) error {
	formMessage := models.FormMessage{
		Recipients: []string{recipientEmail},
		Title:      "Добро пожаловать в почту Mailbx",
		Text: "Это письмо создано автоматически сервером Mailbx.ru, отвечать на него не нужно.\n" +
			"Поздравляем Вас с присоединением к нашей почте. Уверены, что вы останетесь довольны ее использованием!",
		ReplyToMessageID: nil,
	}

	return uc.sendMessageFromSupport(formMessage)
}

func (uc *mailUC) sendMessageFromSupport(message models.FormMessage) error {
	supportAccount, err := uc.getSupportAccount()
	if err != nil {
		return pkgErrors.Wrap(err, "send support message : get support account")
	}

	_, err = uc.SendMessage(supportAccount.UserID, message)
	return err
}

func (uc *mailUC) getSupportAccount() (*models.UserInfo, error) {
	return uc.userUC.GetInfoByEmail("support@mailbx.ru")
}

func (uc *mailUC) MarkMessageAsSeen(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	err := uc.mailRepo.UpdateMessageState(userID, messageID, "seen", true)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "mark message seen : update state")
	}

	return uc.GetMessage(userID, messageID)
}

func (uc *mailUC) MarkMessageAsUnseen(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	err := uc.mailRepo.UpdateMessageState(userID, messageID, "seen", false)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "mark message unseen : update state")
	}

	return uc.GetMessage(userID, messageID)
}

func (uc *mailUC) getFolderByMessage(userID uint64, messageID uint64) (*models.Folder, error) {
	folder, err := uc.mailRepo.SelectFolderByUserNMessage(userID, messageID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "select folder by message")
	}

	return folder, nil
}

func (uc *mailUC) MoveMessageToFolder(userID uint64, messageID uint64, folderSlug string) error {
	message, err := uc.mailRepo.SelectMessageByUserNMessage(userID, messageID)
	if err != nil {
		return pkgErrors.Wrap(err, "get message : by Uid and Mid")
	}
	if message == nil {
		return pkgErrors.WithMessage(errors.ErrMessageNotFound, "get message")
	}

	toFolder, err := uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		return pkgErrors.Wrap(err, "get golder info")
	}

	fromFolder, err := uc.getFolderByMessage(userID, messageID)
	if err != nil {
		return pkgErrors.Wrap(err, "get folder by message")
	}

	if *fromFolder == *toFolder {
		return pkgErrors.WithMessage(errors.ErrMoveToSameFolder, "new folder is equals with old folder")
	}
	if toFolder.LocalName == "drafts" {
		return pkgErrors.WithMessage(errors.ErrMoveToDraftFolder, "new folder is equals draft folder")
	}
	if fromFolder.LocalName == "drafts" {
		return pkgErrors.WithMessage(errors.ErrMoveFromDraftFolder, "old folder is equals draft folder")
	}

	err = uc.mailRepo.UpdateMessageFolder(userID, messageID, toFolder.FolderID)
	if err != nil {
		return pkgErrors.Wrap(err, "update message folder")
	}

	return nil
}

func (uc *mailUC) GetCustomFolders(userID uint64) ([]models.Folder, error) {
	localNames := make([]string, 0, len(defaultFolderNames))
	for localName := range defaultFolderNames {
		localNames = append(localNames, localName)
	}

	folders, err := uc.mailRepo.SelectCustomFoldersByUser(userID, localNames)
	if folders == nil || err != nil {
		return make([]models.Folder, 0), pkgErrors.Wrap(err, "get folders")
	}

	return folders, nil
}
