package user

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

type UseCase interface {
	CreateProfile(profile models.Profile) error
	GetProfileByID(id uint64) (*models.Profile, error)
}

type Repo interface {
	CreateProfile(profile models.Profile) error
	GetProfileByID(id uint64) (*models.Profile, error)
}

type Handlers interface {
	CreateProfile(profile models.Profile) error
	GetProfileByID(id uint64) (*models.Profile, error)
}
