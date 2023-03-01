package inmemory

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"golang.org/x/exp/slices"
)

type MailRepository struct {
	messages   []models.Message
	recipients []models.Recipient
	folders    []models.Folder
	boxes      []models.Box
	states     []models.State
	users      []models.User
}

func NewMailRepository() mail.RepositoryI {
	return &MailRepository{
		messages: []models.Message{
			{1, 1, "2023-01-01", "Title1", "Text1"},
			{2, 1, "2023-01-02", "Title2", "Text2"},
			{3, 2, "2023-01-29", "Title3", "Text3"},
			{4, 3, "2023-01-01", "Title4", "Text4"},
			//{5, 2, "2023-03-01", "Title5", "Text5"},
		},
		recipients: []models.Recipient{
			{1, 2},
			{2, 3},
			{3, 1},
			{4, 2},
			{4, 1},
		},
		folders: []models.Folder{
			{1, "Trash", 1},
			{2, "Spam", 1},
			{3, "Trash", 2},
			{4, "Spam", 2},
			{5, "Trash", 3},
			{6, "Spam", 3},
			{7, "My", 2},
		},
		boxes: []models.Box{
			{7, 1},
		},
		states: []models.State{
			{2, 1, false, false, true},
			{3, 2, true, false, true},
			{1, 3, false, false, true},
			{2, 4, false, false, true},
			{1, 4, false, false, true},
			{1, 1, true, false, true},
			{1, 2, true, false, true},
			{2, 3, true, false, true},
			{3, 4, true, false, true},
		},
		users: []models.User{
			{1, "test@example.com", "1234"},
			{2, "gena@example.com", "4321"},
			{3, "max@example.com", "1379"},
		},
	}
}

func (db *MailRepository) messageHasFolder(userID uint64, messageID uint64) bool {
	for _, b := range db.boxes {
		if b.MessageID == messageID {
			idx := slices.IndexFunc(db.folders, func(folder models.Folder) bool {
				return folder.FolderID == b.FolderID
			})

			if idx != -1 && db.folders[idx].UserID == userID {
				return true
			}
		}
	}

	return false
}

func (db *MailRepository) findOriginalMessage(messageID uint64) (*models.Message, error) {
	idx := slices.IndexFunc(db.messages, func(message models.Message) bool {
		return message.MessageID == messageID
	})

	if idx == -1 {
		return nil, errors.New("message not found")
	}

	return &db.messages[idx], nil
}

func (db *MailRepository) findUser(userID uint64) (*models.User, error) {
	idx := slices.IndexFunc(db.users, func(user models.User) bool {
		return user.ID == userID
	})

	if idx == -1 {
		return nil, errors.New("user not found")
	}

	return &db.users[idx], nil
}

func (db *MailRepository) findState(userID uint64, messageID uint64) (*models.State, error) {
	idx := slices.IndexFunc(db.states, func(state models.State) bool {
		return state.UserID == userID && state.MessageID == messageID
	})

	if idx == -1 {
		return nil, errors.New("state not found")
	}

	return &db.states[idx], nil
}

func (db *MailRepository) findRecipientsEmails(messageID uint64) ([]string, error) {
	var recipientsEmails []string

	for _, r := range db.recipients {
		if r.MessageID == messageID {
			user, err := db.findUser(r.UserID)

			if err != nil {
				return recipientsEmails, err
			}

			recipientsEmails = append(recipientsEmails, user.Email)
		}
	}

	return recipientsEmails, nil
}

func (db *MailRepository) SelectIncomingMessagesByUser(userID uint64) ([]models.IncomingMessage, error) {
	var messages []models.IncomingMessage

	for _, r := range db.recipients {
		if r.UserID == userID {
			if db.messageHasFolder(r.UserID, r.MessageID) {
				continue
			}

			originalMessage, err := db.findOriginalMessage(r.MessageID)

			if err != nil {
				return messages, err
			}

			fromUser, err := db.findUser(originalMessage.UserID)

			if err != nil {
				return messages, err
			}

			state, err := db.findState(r.UserID, r.MessageID)

			if err != nil {
				return messages, err
			}

			currentMessage := models.IncomingMessage{
				MessageID:    r.MessageID,
				FromUser:     fromUser.Email,
				CreatingDate: originalMessage.CreatingDate,
				Title:        originalMessage.Title,
				Text:         originalMessage.Text,
				Read:         state.Read,
				Favorite:     state.Favorite,
			}

			messages = append(messages, currentMessage)
		}
	}

	return messages, nil
}

func (db *MailRepository) SelectOutgoingMessagesByUser(userID uint64) ([]models.OutgoingMessage, error) {
	var messages []models.OutgoingMessage

	for _, m := range db.messages {
		if m.UserID == userID {
			if db.messageHasFolder(userID, m.MessageID) {
				continue
			}

			recipients, err := db.findRecipientsEmails(m.MessageID)

			if len(recipients) == 0 {
				continue
			}

			state, err := db.findState(userID, m.MessageID)

			if err != nil {
				return messages, err
			}

			currentMessage := models.OutgoingMessage{
				MessageID:    m.MessageID,
				ToUsers:      recipients,
				CreatingDate: m.CreatingDate,
				Title:        m.Title,
				Text:         m.Text,
				Read:         state.Read,
				Favorite:     state.Favorite,
			}

			messages = append(messages, currentMessage)
		}
	}

	return messages, nil
}

func (db *MailRepository) SelectFolderByUserNFolder(userID uint64, folderID uint64) (*models.Folder, error) {
	idx := slices.IndexFunc(db.folders, func(folder models.Folder) bool {
		return folder.UserID == userID && folder.FolderID == folderID
	})

	if idx == -1 {
		return nil, errors.New("folder not found")
	}

	return &db.folders[idx], nil
}

func (db *MailRepository) SelectFoldersByUser(userID uint64) []models.Folder {
	var folders []models.Folder

	for _, f := range db.folders {
		if f.UserID == userID {
			folders = append(folders, f)
		}
	}

	return folders
}

func (db *MailRepository) SelectMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.IncomingMessage, error) {
	var messages []models.IncomingMessage

	for _, b := range db.boxes {
		if b.FolderID == folderID {
			originalMessage, err := db.findOriginalMessage(b.MessageID)

			if err != nil {
				return messages, err
			}

			if err != nil {
				return messages, err
			}

			fromUser, err := db.findUser(originalMessage.UserID)

			if err != nil {
				return messages, err
			}

			state, err := db.findState(userID, b.MessageID)

			if err != nil {
				return messages, err
			}

			currentMessage := models.IncomingMessage{
				MessageID:    b.MessageID,
				FromUser:     fromUser.Email,
				CreatingDate: originalMessage.CreatingDate,
				Title:        originalMessage.Title,
				Text:         originalMessage.Text,
				Read:         state.Read,
				Favorite:     state.Favorite,
			}

			messages = append(messages, currentMessage)
		}
	}

	return messages, nil
}
