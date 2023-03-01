package usecase

import (
	"errors"
	mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/internal/mail/reporsitory"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/models"
)

type UseCaseI interface {
	GetIncomingMessages(userID uint64) ([]models.IncomingMessage, error)
	GetOutgoingMessages(userID uint64) ([]models.OutgoingMessage, error)
	GetFolders(userID uint64) []models.Folder
	GetFolderMessages(userID uint64, folderID uint64) ([]models.IncomingMessage, error)
}

type UseCase struct {
	repo mailRepo.RepositoryI
}

func New(rep mailRepo.RepositoryI) UseCaseI {
	return &UseCase{
		repo: rep,
	}
}

func (uc *UseCase) GetFolders(userID uint64) []models.Folder {
	folders := uc.repo.SelectFoldersByUser(userID)

	return folders
}

func (uc *UseCase) GetIncomingMessages(userID uint64) ([]models.IncomingMessage, error) {
	var messages []models.IncomingMessage
	messages, err := uc.repo.SelectIncomingMessagesByUser(userID)

	if err != nil {
		return messages, err
	}

	return messages, nil
}

func (uc *UseCase) GetOutgoingMessages(userID uint64) ([]models.OutgoingMessage, error) {
	var messages []models.OutgoingMessage
	messages, err := uc.repo.SelectOutgoingMessagesByUser(userID)

	if err != nil {
		return messages, err
	}

	return messages, nil
}

func (uc *UseCase) GetFolderMessages(userID uint64, folderID uint64) ([]models.IncomingMessage, error) {
	var messages []models.IncomingMessage

	folder, err := uc.repo.SelectFolderByUserNFolder(userID, folderID)

	if err != nil {
		return messages, err
	}

	if folder == nil {
		return messages, errors.New("folder not found")
	}

	messages, err = uc.repo.SelectMessagesByUserNFolder(userID, folderID)

	if err != nil {
		return messages, err
	}

	return messages, nil
}
