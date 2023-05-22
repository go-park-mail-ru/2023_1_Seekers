package usecase

import (
	"encoding/base64"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/api/ws"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/repository"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/smtp/client"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/rand"
	pkgSmtp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/smtp"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/validation"
	pkgErrors "github.com/pkg/errors"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:generate mockgen -destination=./mocks/mockusecase.go -package=mocks -source=../interface.go

type mailUC struct {
	cfg      *config.Config
	mailRepo mailRepo.MailRepoI
	userUC   user.UseCaseI
	fileUC   file_storage.UseCaseI
	hub      *ws.Hub
}

func New(c *config.Config, mR mailRepo.MailRepoI, uUC user.UseCaseI, fUC file_storage.UseCaseI) mail.UseCaseI {
	return &mailUC{
		cfg:      c,
		mailRepo: mR,
		userUC:   uUC,
		fileUC:   fUC,
	}
}

var defaultFolderNames = map[string]string{
	"inbox":             "Входящие",
	"outbox":            "Исходящие",
	"trash":             "Корзина",
	common.FolderDrafts: "Черновики",
	"spam":              "Спам",
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

	messages, err = uc.mailRepo.SelectFolderMessagesByUserNFolderID(userID, folder.FolderID, folderSlug == common.FolderDrafts)
	if err != nil {
		return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages : msg by user and folder")
	}

	for i, message := range messages {
		messages[i].Preview = common.GetInnerText(message.Text, uc.cfg.Api.MailPreviewMaxLen)
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

		attaches, err := uc.mailRepo.GetMessageAttachments(messageID)
		if err != nil {
			return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages : get message attachments")
		}
		var sumSize int64 = 0
		for _, v := range attaches {
			sumSize += v.SizeCount
		}

		messages[i].Attachments = attaches
		messages[i].AttachmentsSize = common.ByteSize2Str(sumSize)
	}

	return messages, nil
}

func (uc *mailUC) SearchMessages(userID uint64, fromUser, toUser, folder, filter string) ([]models.MessageInfo, error) {
	var messages []models.MessageInfo

	messages, err := uc.mailRepo.SearchMessages(userID, fromUser, toUser, folder, filter)
	if err != nil {
		return []models.MessageInfo{}, pkgErrors.Wrap(err, "SearchMessages : msg by user and folder")
	}

	for i, message := range messages {
		message.Preview = common.GetInnerText(message.Text, uc.cfg.Api.MailPreviewMaxLen)
		messageID := message.MessageID

		fromUser, err := uc.userUC.GetInfo(message.FromUser.UserID)
		if err != nil {
			return []models.MessageInfo{}, pkgErrors.Wrap(err, "SearchMessages : get info by id")
		}

		messages[i].FromUser = *fromUser
		recipientsIDs, err := uc.mailRepo.SelectRecipientsByMessage(messageID, message.FromUser.UserID)
		if err != nil {
			return []models.MessageInfo{}, pkgErrors.Wrap(err, "SearchMessages : get recipients by msg")
		}

		for _, recipientsID := range recipientsIDs {
			profile, err := uc.userUC.GetInfo(recipientsID)
			if err != nil {
				return []models.MessageInfo{}, pkgErrors.Wrap(err, "SearchMessages : get info by id")
			}

			messages[i].Recipients = append(message.Recipients, *profile)
		}
	}

	return messages, nil
}

func (uc *mailUC) SearchRecipients(userID uint64) ([]models.UserInfo, error) {
	recipesInfo, err := uc.mailRepo.SearchRecipients(userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "search recipients")
	}

	return recipesInfo, nil
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

	_, err := uc.mailRepo.SelectFolderByUserNFolderName(userID, form.Name)
	if !pkgErrors.Is(err, errors.ErrFolderNotFound) {
		return nil, pkgErrors.WithMessage(errors.ErrFolderAlreadyExists, "select folder by user and name")
	}

	lastLocalID, err := uc.getLastLocalFolderID(userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get last local folder id")
	}

	localName := strconv.FormatUint(lastLocalID+1, 10)
	folder := &models.Folder{
		UserID:    userID,
		LocalName: localName,
		Name:      form.Name,
	}

	_, err = uc.mailRepo.InsertFolder(folder)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "insert folder")
	}

	return uc.GetFolderInfo(userID, localName)
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

	messages, err := uc.mailRepo.SelectFolderMessagesByUserNFolderID(userID, folder.FolderID, false)
	if err != nil {
		return pkgErrors.Wrap(err, "get folder messages : msg by user and folder")
	}

	for _, message := range messages {
		err = uc.MoveMessageToFolder(userID, message.MessageID, folder.LocalName, "trash")
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
	attaches, err := uc.mailRepo.GetMessageAttachments(messageID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get folder messages : get message attachments")
	}
	var sumSize int64 = 0
	for _, v := range attaches {
		sumSize += v.SizeCount
	}

	firstMessage.Attachments = attaches
	firstMessage.AttachmentsSize = common.ByteSize2Str(sumSize)
	firstMessage.Preview = common.GetInnerText(firstMessage.Text, uc.cfg.Api.MailPreviewMaxLen)
	return firstMessage, nil
}

func (uc *mailUC) DeleteMessage(userID uint64, messageID uint64, folderSlug string) error {
	message, err := uc.mailRepo.SelectMessageByUserNMessage(userID, messageID)
	if err != nil {
		return pkgErrors.Wrap(err, "get message : by Uid and Mid")
	}
	if message == nil {
		return pkgErrors.WithMessage(errors.ErrMessageNotFound, "get message")
	}

	folder, err := uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		return pkgErrors.Wrap(err, "get folder info")
	}

	boxExists, err := uc.checkExistingBox(userID, messageID, folder.FolderID)
	if err != nil {
		return pkgErrors.Wrap(err, "check existing box")
	}
	if !boxExists {
		return pkgErrors.WithMessage(errors.ErrBoxNotFound, "box not found")
	}

	switch folder.LocalName {
	case "trash":
		err = uc.mailRepo.DeleteBox(userID, messageID, folder.FolderID)
		if err != nil {
			return pkgErrors.Wrap(err, "delete message for user")
		}
	case common.FolderDrafts:
		err = uc.mailRepo.DeleteMessageFromMessages(messageID)
		if err != nil {
			return pkgErrors.Wrap(err, "delete message full")
		}
	default:
		err = uc.MoveMessageToFolder(userID, messageID, folder.LocalName, "trash")
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
	folder, err := uc.GetFolderInfo(fromUserID, common.FolderDrafts)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message : get folder by UId and FolderSlug")
	}

	var user2folder []models.User2Folder
	user2folder = append(user2folder, models.User2Folder{
		UserID:   fromUserID,
		FolderID: folder.FolderID,
	})

	attachesInfo := make([]models.AttachmentInfo, len(message.Attachments))
	for i, a := range message.Attachments {
		raw, err := base64.StdEncoding.DecodeString(a.FileData)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "failed base64 decode attach")
		}

		contentType := http.DetectContentType(raw)
		fileSize := int64(len(raw))
		attachesInfo[i] = models.AttachmentInfo{
			FileName:  a.FileName,
			FileData:  raw,
			S3FName:   rand.FileName("", filepath.Ext(a.FileName)),
			Type:      contentType,
			SizeStr:   common.ByteSize2Str(fileSize),
			SizeCount: fileSize,
		}
	}

	var sumSize int64 = 0
	for _, v := range attachesInfo {
		sumSize += v.SizeCount
	}

	for _, email := range message.Recipients {
		recipient, err := uc.userUC.GetInfoByEmail(email)
		if err != nil {
			exUser, err := uc.userUC.Create(&models.User{
				Email:      email,
				Password:   uc.cfg.UserService.ExternalUserPassword,
				IsExternal: true,
			})
			if err != nil {
				return nil, pkgErrors.Wrap(err, "send message : create external user")
			}

			_, err = uc.CreateDefaultFolders(exUser.UserID)
			if err != nil {
				return nil, pkgErrors.Wrap(err, "send message : create external user")
			}

			recipient = &models.UserInfo{
				UserID:    exUser.UserID,
				FirstName: exUser.FirstName,
				LastName:  exUser.LastName,
				Email:     exUser.Email,
			}
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
		Attachments:      attachesInfo,
		AttachmentsSize:  common.ByteSize2Str(sumSize),
		Preview:          common.GetInnerText(message.Text, uc.cfg.Api.MailPreviewMaxLen),
		IsDraft:          true,
	}

	err = uc.mailRepo.InsertMessage(fromUserID, &newMessage, user2folder)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message : insert message")
	}

	for _, a := range newMessage.Attachments {
		if err := uc.fileUC.Upload(&models.S3File{
			Bucket: uc.cfg.S3.S3AttachBucket,
			Name:   a.S3FName,
			Data:   a.FileData,
		}); err != nil {
			return nil, pkgErrors.Wrap(err, "send message : put to s3")
		}
	}

	return uc.GetMessage(fromUserID, newMessage.MessageID)
}

func (uc *mailUC) mapRecipients(newRecipients []string, oldMessage *models.MessageInfo) map[string]string {
	recipients := make(map[string]string)

	for _, recipient := range newRecipients {
		recipients[recipient] = common.ActionAdd
	}

	for _, recipient := range oldMessage.Recipients {
		if _, ok := recipients[recipient.Email]; ok {
			recipients[recipient.Email] = common.ActionSave
		} else {
			recipients[recipient.Email] = common.ActionDelete
		}
	}

	return recipients
}

func (uc *mailUC) EditDraft(fromUserID uint64, messageID uint64, formMessage models.FormEditMessage) (*models.MessageInfo, error) {
	//TODO
	message, err := uc.GetMessage(fromUserID, messageID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get message")
	}
	if !message.IsDraft {
		return nil, pkgErrors.WithMessage(errors.ErrCantEditSentMessage, "can't edit not draft message")
	}

	recipients := uc.mapRecipients(formMessage.Recipients, message)
	var toInsert []models.User2Folder
	var toDelete []models.User2Folder

	for email, value := range recipients {
		recipient, err := uc.userUC.GetInfoByEmail(email)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "edit draft : get user info by email")
		}

		folder, err := uc.GetFolderInfo(recipient.UserID, "inbox")
		if err != nil {
			return nil, pkgErrors.Wrap(err, "edit draft: get folder by UId and FolderSlug")
		}

		switch value {
		case common.ActionAdd:
			toInsert = append(toInsert, models.User2Folder{
				UserID:   recipient.UserID,
				FolderID: folder.FolderID,
			})
		case common.ActionDelete:
			toDelete = append(toDelete, models.User2Folder{
				UserID:   recipient.UserID,
				FolderID: folder.FolderID,
			})
		}
	}

	message.Title = formMessage.Title
	message.Text = formMessage.Text
	message.ReplyToMessageID = formMessage.ReplyToMessageID
	message.CreatedAt = common.GetCurrentTime(uc.cfg.Logger.LogsTimeFormat)
	message.Preview = common.GetInnerText(message.Text, uc.cfg.Api.MailPreviewMaxLen)

	if err := uc.mailRepo.UpdateMessage(message, toInsert, toDelete); err != nil {
		return nil, pkgErrors.Wrap(err, "edit draft : update message")
	}

	return uc.GetMessage(fromUserID, messageID)
}

func (uc *mailUC) SendMessage(userID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	usr, err := uc.userUC.GetByEmail(message.FromUser)
	if err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInvalidForm, "send message - invalid sender field")
	}
	if usr != nil {
		if usr.UserID != userID {
			return nil, pkgErrors.WithMessage(errors.ErrInvalidForm, "send message - sender not compare")
		}
	}

	if len(message.Recipients) == 0 {
		return nil, pkgErrors.WithMessage(errors.ErrNoValidEmails, "send message")
	}

	var user2folder []models.User2Folder

	fromUser, err := uc.userUC.GetByEmail(message.FromUser)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message")
	}

	if !fromUser.IsExternal {
		folder, err := uc.GetFolderInfo(fromUser.UserID, "outbox")
		if err != nil {
			return nil, pkgErrors.Wrap(err, "send message : get folder by UId and FolderSlug")
		}

		user2folder = append(user2folder, models.User2Folder{
			UserID:   fromUser.UserID,
			FolderID: folder.FolderID,
		})
	}

	attachesInfo := make([]models.AttachmentInfo, len(message.Attachments))
	for i, a := range message.Attachments {
		raw, err := base64.StdEncoding.DecodeString(a.FileData)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "failed base64 decode attach")
		}

		contentType := http.DetectContentType(raw)
		fileSize := int64(len(raw))
		attachesInfo[i] = models.AttachmentInfo{
			FileName:  a.FileName,
			FileData:  raw,
			S3FName:   rand.FileName("", filepath.Ext(a.FileName)),
			Type:      contentType,
			SizeStr:   common.ByteSize2Str(fileSize),
			SizeCount: fileSize,
		}
	}

	var sumSize int64 = 0
	for _, v := range attachesInfo {
		sumSize += v.SizeCount
	}

	newMessage := models.MessageInfo{
		Title:            message.Title,
		CreatedAt:        common.GetCurrentTime(uc.cfg.Logger.LogsTimeFormat),
		Text:             message.Text,
		Attachments:      attachesInfo,
		AttachmentsSize:  common.ByteSize2Str(sumSize),
		ReplyToMessageID: message.ReplyToMessageID,
		Preview:          common.GetInnerText(message.Text, uc.cfg.Api.MailPreviewMaxLen),
	}

	for _, recipient := range message.Recipients {
		toDomain, err := pkgSmtp.ParseDomain(recipient)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "send message - failed get recipient domain")
		}
		if toDomain != uc.cfg.Mail.PostDomain {
			if err := client.SendMail(fromUser, recipient, &newMessage, uc.cfg.Mail.PostDomain, uc.cfg.SmtpServer.SecretPassword); err != nil {
				return nil, pkgErrors.Wrap(err, "send message : to other mail service")
			}
			var recipientU *models.User
			if recipientU, err = uc.userUC.GetByEmail(recipient); err != nil {
				recipientU, err = uc.userUC.Create(&models.User{
					Email:      recipient,
					Password:   uc.cfg.UserService.ExternalUserPassword,
					FirstName:  "",
					LastName:   "",
					IsExternal: true,
				})
				if err != nil {
					return nil, pkgErrors.Wrap(err, "mail service - send message : create external user")
				}

				_, err = uc.CreateDefaultFolders(recipientU.UserID)
				if err != nil {
					return nil, pkgErrors.Wrap(err, "mail service - send message  : create default folders for external user")
				}
			}

			folder, err := uc.GetFolderInfo(recipientU.UserID, "inbox")
			if err != nil {
				return nil, pkgErrors.Wrap(err, "send message - external: get folder by UId and FolderSlug")
			}

			user2folder = append(user2folder, models.User2Folder{
				UserID:   recipientU.UserID,
				FolderID: folder.FolderID,
			})
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

	err = uc.mailRepo.InsertMessage(fromUser.UserID, &newMessage, user2folder)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message : insert message")
	}

	for _, a := range newMessage.Attachments {
		if err := uc.fileUC.Upload(&models.S3File{
			Bucket: uc.cfg.S3.S3AttachBucket,
			Name:   a.S3FName,
			Data:   a.FileData,
		}); err != nil {
			return nil, pkgErrors.Wrap(err, "send message : put to s3")
		}
	}
	m, err := uc.GetMessage(userID, newMessage.MessageID)
	if err != nil {
		return &newMessage, nil
	}
	return m, nil
}

func (uc *mailUC) SendFailedSendingMessage(recipientEmail string, invalidEmails []string) (*models.MessageInfo, error) {
	formMessage := models.FormMessage{
		Recipients: []string{recipientEmail},
		Title:      "Ваше сообщение не доставлено",
		Text: "Это письмо создано автоматически сервером mailbx.ru, отвечать на него не нужно.\n\n" +
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

	_, err := uc.sendMessageFromSupport(formMessage)
	return err
}

func (uc *mailUC) sendMessageFromSupport(message models.FormMessage) (*models.MessageInfo, error) {
	supportAccount, err := uc.getSupportAccount()
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send support message : get support account")
	}

	message.FromUser = supportAccount.Email
	return uc.SendMessage(supportAccount.UserID, message)
}

func (uc *mailUC) getSupportAccount() (*models.UserInfo, error) {
	return uc.userUC.GetInfoByEmail("support@mailbx.ru")
}

func (uc *mailUC) checkExistingBox(userID uint64, messageID uint64, folderID uint64) (bool, error) {
	boxExists, err := uc.mailRepo.CheckExistingBox(userID, messageID, folderID)
	if err != nil {
		return false, pkgErrors.Wrap(err, "check existing box")
	}

	return boxExists, nil
}

func (uc *mailUC) MarkMessageAsSeen(userID uint64, messageID uint64, folderSlug string) (*models.MessageInfo, error) {
	message, err := uc.mailRepo.SelectMessageByUserNMessage(userID, messageID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get message : by Uid and Mid")
	}
	if message == nil {
		return nil, pkgErrors.WithMessage(errors.ErrMessageNotFound, "get message")
	}

	folder, err := uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get folder info")
	}

	boxExists, err := uc.checkExistingBox(userID, messageID, folder.FolderID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "check existing box")
	}
	if !boxExists {
		return nil, pkgErrors.WithMessage(errors.ErrBoxNotFound, "box not found")
	}

	err = uc.mailRepo.UpdateMessageState(userID, messageID, folder.FolderID, "seen", true)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "mark message seen : update state")
	}

	return uc.GetMessage(userID, messageID)
}

func (uc *mailUC) GetAttachInfo(attachID, userID uint64) (*models.AttachmentInfo, error) {
	attach, err := uc.mailRepo.GetAttach(attachID, userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get attach")
	}

	return attach, nil
}

func (uc *mailUC) MarkMessageAsUnseen(userID uint64, messageID uint64, folderSlug string) (*models.MessageInfo, error) {
	message, err := uc.mailRepo.SelectMessageByUserNMessage(userID, messageID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get message : by Uid and Mid")
	}
	if message == nil {
		return nil, pkgErrors.WithMessage(errors.ErrMessageNotFound, "get message")
	}

	folder, err := uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get folder info")
	}

	boxExists, err := uc.checkExistingBox(userID, messageID, folder.FolderID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "check existing box")
	}
	if !boxExists {
		return nil, pkgErrors.WithMessage(errors.ErrBoxNotFound, "box not found")
	}

	err = uc.mailRepo.UpdateMessageState(userID, messageID, folder.FolderID, "seen", false)
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

func (uc *mailUC) MoveMessageToFolder(userID uint64, messageID uint64, fromFolderSlug string, toFolderSlug string) error {
	message, err := uc.mailRepo.SelectMessageByUserNMessage(userID, messageID)
	if err != nil {
		return pkgErrors.Wrap(err, "get message : by Uid and Mid")
	}
	if message == nil {
		return pkgErrors.WithMessage(errors.ErrMessageNotFound, "get message")
	}

	fromFolder, err := uc.GetFolderInfo(userID, fromFolderSlug)
	if err != nil {
		return pkgErrors.Wrap(err, "get fromFolder info")
	}

	boxExists, err := uc.checkExistingBox(userID, messageID, fromFolder.FolderID)
	if err != nil {
		return pkgErrors.Wrap(err, "check existing box")
	}
	if !boxExists {
		return pkgErrors.WithMessage(errors.ErrBoxNotFound, "box not found")
	}

	if fromFolderSlug == toFolderSlug {
		return pkgErrors.WithMessage(errors.ErrMoveToSameFolder, "new fromFolder is equals with old fromFolder")
	}

	if fromFolderSlug == common.FolderDrafts {
		return pkgErrors.WithMessage(errors.ErrMoveFromDraftFolder, "old fromFolder is equals draft fromFolder")
	}

	toFolder, err := uc.GetFolderInfo(userID, toFolderSlug)
	if err != nil {
		return pkgErrors.Wrap(err, "get golder info (to fromFolder slug)")
	}

	if toFolder.LocalName == common.FolderDrafts {
		return pkgErrors.WithMessage(errors.ErrMoveToDraftFolder, "new fromFolder is equals draft fromFolder")
	}

	boxExists, err = uc.checkExistingBox(userID, messageID, toFolder.FolderID)
	if err != nil {
		return pkgErrors.Wrap(err, "check existing box")
	}
	if boxExists {
		err = uc.mailRepo.DeleteBox(userID, messageID, fromFolder.FolderID)
		if err != nil {
			return pkgErrors.Wrap(err, "delete message for user")
		}
	} else {
		err = uc.mailRepo.UpdateMessageFolder(userID, messageID, fromFolder.FolderID, toFolder.FolderID)
		if err != nil {
			return pkgErrors.Wrap(err, "update message fromFolder")
		}
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

func (uc *mailUC) GetAttach(attachID, userID uint64) (*models.AttachmentInfo, error) {
	attach, err := uc.mailRepo.GetAttach(attachID, userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get attach")
	}

	file, err := uc.fileUC.Get(uc.cfg.S3.S3AttachBucket, attach.S3FName)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get data from s3")
	}

	attach.FileData = file.Data

	return attach, nil
}
