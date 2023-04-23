package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/repository"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"strings"
)

//go:generate mockgen -destination=./mocks/mockusecase.go -package=mocks -source=../interface.go

type UseCase struct {
	mailRepo mailRepo.MailRepoI
	userUC   user.UseCaseI
}

func New(mR mailRepo.MailRepoI, uUC user.UseCaseI) mail.UseCaseI {
	return &UseCase{
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

func (uc *UseCase) GetFolders(userID uint64) ([]models.Folder, error) {
	folders, err := uc.mailRepo.SelectFoldersByUser(userID)
	if err != nil {
		return []models.Folder{}, pkgErrors.Wrap(err, "get folders")
	}

	return folders, nil
}

func (uc *UseCase) GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error) {
	folder, err := uc.mailRepo.SelectFolderByUserNFolder(userID, folderSlug)
	if err != nil {
		return folder, pkgErrors.Wrap(err, "get folder info")
	}
	if folder == nil {
		return nil, pkgErrors.WithMessage(errors.ErrFolderNotFound, "get folder info")
	}

	return folder, nil
}

func (uc *UseCase) GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error) {
	var messages []models.MessageInfo

	folder, err := uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages")
	}

	messages, err = uc.mailRepo.SelectFolderMessagesByUserNFolder(userID, folder.FolderID)
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

func (uc *UseCase) CreateDefaultFolders(userID uint64) ([]models.Folder, error) {
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

func (uc *UseCase) GetMessage(userID uint64, messageID uint64) (*models.MessageInfo, error) {
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

func (uc *UseCase) ValidateRecipients(recipients []string) ([]string, []string) {
	var validEmails []string
	var invalidEmails []string

	for _, email := range recipients {
		_, err := uc.userUC.GetInfoByEmail(email)
		if err != nil {
			invalidEmails = append(invalidEmails, email)
		} else {
			validEmails = append(validEmails, email)
		}
	}

	return validEmails, invalidEmails
}

func (uc *UseCase) SendMessage(fromUserID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	if len(message.Recipients) == 0 {
		return nil, pkgErrors.WithMessage(errors.ErrNoValidEmails, "send message")
	}

	folder, err := uc.GetFolderInfo(fromUserID, "outbox")
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
		CreatedAt:        common.GetCurrentTime(),
		Text:             message.Text,
		ReplyToMessageID: message.ReplyToMessageID,
	}

	err = uc.mailRepo.InsertMessage(fromUserID, &newMessage, user2folder)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message : insert message")
	}

	return uc.GetMessage(fromUserID, newMessage.MessageID)
}

func (uc *UseCase) SendFailedSendingMessage(recipientEmail string, invalidEmails []string) error {
	formMessage := models.FormMessage{
		Recipients: []string{recipientEmail},
		Title:      "Ваше сообщение не доставлено",
		Text: "Это письмо создано автоматически сервером Mailbox.ru, отвечать на него не нужно.\n\n" +
			"К сожалению, Ваше письмо не может быть доставлено одному или нескольким получателям:\n" +
			strings.Join(invalidEmails, "\n") + "\n\nРекомендуем Вам проверить корректность указания адресов получателей.",
		ReplyToMessageID: nil,
	}

	return uc.sendMessageFromSupport(formMessage)
}

func (uc *UseCase) SendWelcomeMessage(recipientEmail string) error {
	formMessage := models.FormMessage{
		Recipients: []string{recipientEmail},
		Title:      "Добро пожаловать в почту Mailbox",
		Text: "Это письмо создано автоматически сервером Mailbox.ru, отвечать на него не нужно.\n" +
			"Поздравляем Вас с присоединением к нашей почте. Уверены, что вы останетесь довольны ее использванием!",
		ReplyToMessageID: nil,
	}

	return uc.sendMessageFromSupport(formMessage)
}

func (uc *UseCase) sendMessageFromSupport(message models.FormMessage) error {
	supportAccount, err := uc.getSupportAccount()
	if err != nil {
		return pkgErrors.Wrap(err, "send support message : get support account")
	}

	_, err = uc.SendMessage(supportAccount.UserID, message)
	return err
}

func (uc *UseCase) getSupportAccount() (*models.UserInfo, error) {
	return uc.userUC.GetInfoByEmail("support@mailbox.ru")
}

func (uc *UseCase) MarkMessageAsSeen(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	err := uc.mailRepo.UpdateMessageState(userID, messageID, "seen", true)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "mark message seen : update state")
	}

	return uc.GetMessage(userID, messageID)
}

func (uc *UseCase) MarkMessageAsUnseen(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	err := uc.mailRepo.UpdateMessageState(userID, messageID, "seen", false)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "mark message unseen : update state")
	}

	return uc.GetMessage(userID, messageID)
}
