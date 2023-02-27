package user

import "github.com/go-park-mail-ru/2023_1_Seekers/app/model"

type UseCase interface {
	CreateProfile(profile model.Profile) error
	GetProfileById(id uint64) (*model.Profile, error)
}

type Repo interface {
	CreateProfile(profile model.Profile) error
	GetProfileById(id uint64) (*model.Profile, error)
}

type Handlers interface {
	CreateProfile(profile model.Profile) error
	GetProfileById(id uint64) (*model.Profile, error)
}
