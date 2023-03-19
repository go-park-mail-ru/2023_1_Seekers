package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"time"
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
		return folders, err
	}

	return folders, nil
}

func (uc *UseCase) GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error) {
	var messages []models.MessageInfo

	folder, err := uc.repoMail.SelectFolderByUserNFolder(userID, folderSlug)
	if err != nil {
		return messages, err
	}

	if folder == nil {
		return messages, errors.New("folder not found")
	}

	messages, err = uc.repoMail.SelectFolderMessagesByUserNFolder(userID, folder.FolderID)
	if err != nil {
		return messages, err
	}

	for i, message := range messages {
		messageID := message.MessageID

		fromUser, err := uc.repoUser.GetUserInfoByID(message.FromUser.UserID)
		if err != nil {
			return []models.MessageInfo{}, err
		}

		messages[i].FromUser = *fromUser
		recipientsIDs, err := uc.repoMail.SelectRecipientsByMessage(messageID, message.FromUser.UserID)
		if err != nil {
			return []models.MessageInfo{}, err
		}

		for _, recipientsID := range recipientsIDs {
			profile, err := uc.repoUser.GetUserInfoByID(recipientsID)
			if err != nil {
				return []models.MessageInfo{}, err
			}

			messages[i].Recipients = append(message.Recipients, *profile)
		}
	}

	return messages, nil
}

func (uc *UseCase) GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error) {

	folder, err := uc.repoMail.SelectFolderByUserNFolder(userID, folderSlug)

	if err != nil {
		return folder, err
	}

	if folder == nil {
		return folder, errors.New("folder not found")
	}

	return folder, nil
}

func (uc *UseCase) CreateHelloMessage(to uint64) error {
	now := time.Now()
	msg := models.MessageInfo{
		FromUser: models.UserInfo{
			UserID:    1,
			FirstName: "",
			LastName:  "",
			Email:     "",
		},
		CreatedAt: now.Format("2006-02-01"),
		Title:     "Hello! Its your first mail",
		Text:      "Support of mail box is glad to see You here! Have a nice day!",
	}
	return uc.CreateMessage(msg, to)
}

func (uc *UseCase) CreateMessage(message models.MessageInfo, to ...uint64) error {
	return uc.repoMail.CreateMessage(message, to...)
}
