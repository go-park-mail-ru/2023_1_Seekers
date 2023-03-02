package inmemory

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"golang.org/x/exp/slices"
)

type mailDB struct {
	messages   []models.Message
	recipients []models.Recipient
	folders    []models.Folder
	boxes      []models.Box
	states     []models.State
	usersRepo  _user.RepoI
}

func New(ur _user.RepoI) mail.RepoI {
	return &mailDB{
		messages: []models.Message{
			{1, 1, "2023-01-01", "Invitation", "Hello, we decided to invite you to our party, lets go it will be fine!"},
			{2, 1, "2023-01-02", "Spam letter", "Nunc non velit commodo, vestibulum enim ullamcorper, lobortis mi. Integer eu elit nibh. Integer bibendum semper arcu, eget consectetur nisi gravida eu. Suspendisse maximus id urna a volutpat. Quisque nec iaculis purus, non facilisis massa. Maecenas finibus dui ipsum, ut tempor sapien tincidunt blandit. Ut at iaculis eros, ultrices iaculis nibh. Mauris fermentum elit erat, at cursus urna euismod vel. In congue, ipsum a fermentum semper, dolor sem scelerisque leo, a tempus risus orci eu leo. Fusce vulputate venenatis imperdiet. Vestibulum interdum pellentesque facilisis"},
			{3, 1, "2023-01-04", "Lorem", "Mauris imperdiet massa ante. Pellentesque feugiat nisl nec ultrices laoreet. Aenean a mauris mi. Sed auctor egestas nulla et vulputate. Praesent lobortis nulla ante, vel dignissim odio aliquet et. Suspendisse potenti. Donec venenatis nibh a sem consectetur, bibendum consectetur metus venenatis. Mauris lorem tellus, finibus id dui sit amet, facilisis fermentum orci. Mauris arcu ante, lacinia vitae orci in, tempus elementum lacus. Donec eu augue vulputate, tempor neque nec, efficitur purus. Mauris ut lorem non sapien placerat mattis. In in lacus a lorem viverra laoreet ut et orci. Maecenas auctor, justo nec hendrerit interdum, nibh nisi consectetur sapien, id ultrices lacus mi sed risus. "},
			{4, 1, "2023-01-05", "Very interesting letter", "Morbi sit amet porttitor sapien, eget venenatis est. Suspendisse sollicitudin elit velit, quis sodales dolor maximus id. Vestibulum gravida scelerisque nibh, sit amet tincidunt augue gravida nec. Maecenas non placerat justo, at feugiat nulla. Phasellus dapibus a mi ut interdum. Aliquam nec quam feugiat, rutrum urna ut, cursus purus. Lorem ipsum dolor sit amet, consectetur adipiscing elit. "},
			{5, 1, "2023-01-06", "Small text letter", "Hi! how are you?"},
			{5, 1, "2023-01-06", "Title", "Open this letter please"},
			{6, 1, "2023-01-07", "Advertisement", "Hi, visit our shop!"},
			{7, 2, "2023-01-29", "Title2", "Text2"},
			{8, 2, "2023-01-29", "Title3", "Text3"},
			{9, 3, "2023-01-01", "Title4", "Text4"},
			//{5, 2, "2023-03-01", "Title5", "Text5"},
		},
		recipients: []models.Recipient{
			{1, 2},
			{2, 3},
			{3, 2},
			{4, 3},
			{4, 2},
			{5, 2},
			{5, 3},
			{6, 2},
			{7, 3},
			{8, 1},
			{9, 1},
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
			{1, 1, true, false, true},
			{3, 2, false, false, true},
			{1, 2, true, false, true},
			{2, 3, false, false, true},
			{1, 3, true, false, true},
			{3, 4, false, false, true},
			{2, 4, false, false, true},
			{1, 4, true, false, true},
			{2, 5, false, false, true},
			{3, 5, false, false, true},
			{1, 5, true, false, true},
			{2, 6, false, false, true},
			{1, 6, true, false, true},
			{3, 7, false, false, true},
			{2, 7, true, false, true},
			{1, 8, false, false, true},
			{2, 8, true, false, true},
			{1, 9, false, false, true},
			{3, 9, true, false, true},
			{3, 2, true, false, true},
			{1, 3, false, false, true},
			{2, 4, false, false, true},
			{1, 4, false, false, true},
			{1, 2, true, false, true},
			{2, 3, true, false, true},
			{3, 4, true, false, true},
		},
		usersRepo: ur,
	}
}

func (db *mailDB) messageHasFolder(userID uint64, messageID uint64) bool {
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

func (db *mailDB) findOriginalMessage(messageID uint64) (*models.Message, error) {
	idx := slices.IndexFunc(db.messages, func(message models.Message) bool {
		return message.MessageID == messageID
	})

	if idx == -1 {
		return nil, errors.New("message not found")
	}

	return &db.messages[idx], nil
}

func (db *mailDB) findState(userID uint64, messageID uint64) (*models.State, error) {
	idx := slices.IndexFunc(db.states, func(state models.State) bool {
		return state.UserID == userID && state.MessageID == messageID
	})

	if idx == -1 {
		return nil, errors.New("state not found")
	}

	return &db.states[idx], nil
}

func (db *mailDB) findRecipientsEmails(messageID uint64) ([]string, error) {
	var recipientsEmails []string

	for _, r := range db.recipients {
		if r.MessageID == messageID {
			user, err := db.usersRepo.GetUserByID(r.UserID)

			if err != nil {
				return recipientsEmails, err
			}

			recipientsEmails = append(recipientsEmails, user.Email)
		}
	}

	return recipientsEmails, nil
}

func (db *mailDB) SelectIncomingMessagesByUser(userID uint64) ([]models.IncomingMessage, error) {
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

			fromUser, err := db.usersRepo.GetUserByID(originalMessage.UserID)

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

func (db *mailDB) SelectOutgoingMessagesByUser(userID uint64) ([]models.OutgoingMessage, error) {
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

func (db *mailDB) SelectFolderByUserNFolder(userID uint64, folderID uint64) (*models.Folder, error) {
	idx := slices.IndexFunc(db.folders, func(folder models.Folder) bool {
		return folder.UserID == userID && folder.FolderID == folderID
	})

	if idx == -1 {
		return nil, errors.New("folder not found")
	}

	return &db.folders[idx], nil
}

func (db *mailDB) SelectFoldersByUser(userID uint64) []models.Folder {
	var folders []models.Folder

	for _, f := range db.folders {
		if f.UserID == userID {
			folders = append(folders, f)
		}
	}

	return folders
}

func (db *mailDB) SelectMessagesByUserNFolder(userID uint64, folderID uint64) ([]models.IncomingMessage, error) {
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

			fromUser, err := db.usersRepo.GetUserByID(originalMessage.UserID)

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
