package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"strings"
)

type UseCase struct {
	repoMail mail.RepoI
	repoUser user.RepoI
}

func New(repoMail mail.RepoI, repoUser user.RepoI) mail.UseCaseI {
	return &UseCase{
		repoMail: repoMail,
		repoUser: repoUser,
	}
}

func (uc *UseCase) GetFolders(userID uint64) ([]models.Folder, error) {
	folders, err := uc.repoMail.SelectFoldersByUser(userID)
	if err != nil {
		return folders, pkgErrors.Wrap(err, "get folders")
	}

	return folders, nil
}

func (uc *UseCase) GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error) {
	folder, err := uc.repoMail.SelectFolderByUserNFolder(userID, folderSlug)
	if err != nil {
		return folder, pkgErrors.Wrap(err, "get folder info")
	}
	if folder == nil {
		return folder, pkgErrors.WithMessage(errors.ErrFolderNotFound, "get folder info")
	}

	return folder, nil
}

func (uc *UseCase) GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error) {
	var messages []models.MessageInfo

	folder, err := uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		return messages, pkgErrors.Wrap(err, "get folder messages")
	}

	messages, err = uc.repoMail.SelectFolderMessagesByUserNFolder(userID, folder.FolderID)
	if err != nil {
		return messages, pkgErrors.Wrap(err, "get folder messages : msg by user and folder")
	}

	for i, message := range messages {
		messageID := message.MessageID

		fromUser, err := uc.repoUser.GetInfoByID(message.FromUser.UserID)
		if err != nil {
			return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages : get info by id")
		}

		messages[i].FromUser = *fromUser
		recipientsIDs, err := uc.repoMail.SelectRecipientsByMessage(messageID, message.FromUser.UserID)
		if err != nil {
			return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages : get recipients by msg")
		}

		for _, recipientsID := range recipientsIDs {
			profile, err := uc.repoUser.GetInfoByID(recipientsID)
			if err != nil {
				return []models.MessageInfo{}, pkgErrors.Wrap(err, "get folder messages : get info by id")
			}

			messages[i].Recipients = append(message.Recipients, *profile)
		}
	}

	return messages, nil
}

func (uc *UseCase) GetMessage(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	var firstMessage *models.MessageInfo
	var prevMessage *models.MessageInfo
	replyToMsgID := &messageID

	for replyToMsgID != nil {
		curMessage, err := uc.repoMail.SelectMessageByUserNMessage(userID, *replyToMsgID)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "get message : by Uid and Mid")
		}
		if curMessage == nil {
			return nil, pkgErrors.WithMessage(errors.ErrMessageNotFound, "get message")
		}

		fromUser, err := uc.repoUser.GetInfoByID(curMessage.FromUser.UserID)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "get message : get info by Uid")
		}

		curMessage.FromUser = *fromUser
		recipientsIDs, err := uc.repoMail.SelectRecipientsByMessage(*replyToMsgID, curMessage.FromUser.UserID)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "get message : get recipients by Mid")
		}

		for _, recipientsID := range recipientsIDs {
			profile, err := uc.repoUser.GetInfoByID(recipientsID)
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
		_, err := uc.repoUser.GetInfoByEmail(email)
		if err != nil {
			invalidEmails = append(invalidEmails, email)
		} else {
			validEmails = append(validEmails, email)
		}
	}

	return validEmails, invalidEmails
}

func (uc *UseCase) SendMessage(userID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	if len(message.Recipients) == 0 {
		return nil, pkgErrors.WithMessage(errors.ErrNoValidEmails, "send message")
	}

	newMessage, err := uc.createMessage(userID, &message)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "send message")
	}

	newMessage.Seen = true
	if err := uc.insertMessageToFolder(newMessage.FromUser.UserID, "outbox", newMessage); err != nil {
		return nil, pkgErrors.Wrap(err, "send message")
	}

	newMessage.Seen = false
	for _, email := range message.Recipients {
		recipient, err := uc.repoUser.GetInfoByEmail(email)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "send message : get recipient info by email")
		}

		if err := uc.insertMessageToFolder(recipient.UserID, "inbox", newMessage); err != nil {
			return nil, pkgErrors.Wrap(err, "send message")
		}
	}

	return uc.GetMessage(userID, newMessage.MessageID)
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

func (uc *UseCase) createMessage(userID uint64, message *models.FormMessage) (*models.MessageInfo, error) {
	fromUser, err := uc.repoUser.GetInfoByID(userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "create message : get info by Id")
	}

	newMessage := models.MessageInfo{
		FromUser:         *fromUser,
		Title:            message.Title,
		CreatedAt:        pkg.GetCurrentTime(),
		Text:             message.Text,
		ReplyToMessageID: message.ReplyToMessageID,
	}

	messageID, err := uc.repoMail.InsertMessageToMessages(&newMessage)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "create message : insert to messages")
	}

	newMessage.MessageID = messageID
	return &newMessage, nil
}

func (uc *UseCase) sendMessageFromSupport(message models.FormMessage) error {
	supportAccount, err := uc.getSupportAccount()
	if err != nil {
		return pkgErrors.Wrap(err, "send support message : get support account")
	}

	_, err = uc.SendMessage(supportAccount.UserID, message)
	return err
}

func (uc *UseCase) insertMessageToFolder(userID uint64, folderSlug string, message *models.MessageInfo) error {
	folder, err := uc.repoMail.SelectFolderByUserNFolder(userID, folderSlug)
	if err != nil {
		return pkgErrors.Wrap(err, "insert message to folder : get folder by Uid and Fid")
	}

	return uc.repoMail.InsertMessageToBoxes(userID, folder.FolderID, message)
}

func (uc *UseCase) getSupportAccount() (*models.UserInfo, error) {
	return uc.repoUser.GetInfoByEmail("support@mailbox.ru")
}

func (uc *UseCase) MarkMessageAsSeen(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	err := uc.repoMail.UpdateMessageState(userID, messageID, "seen", true)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "mark message seen : update state")
	}

	return uc.GetMessage(userID, messageID)
}

func (uc *UseCase) MarkMessageAsUnseen(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	err := uc.repoMail.UpdateMessageState(userID, messageID, "seen", false)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "mark message unseen : update state")
	}

	return uc.GetMessage(userID, messageID)
}
