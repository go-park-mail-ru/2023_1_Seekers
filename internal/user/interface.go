package user

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

type UseCaseI interface {
	CreateUser(user models.User) (*models.User, error)
	DeleteUser(user models.User) error
	GetUserByID(ID uint64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateProfile(profile models.Profile) error
	GetProfileByID(id uint64) (*models.Profile, error)
}

type RepoI interface {
	CreateUser(user models.User) (*models.User, error)
	DeleteUser(user models.User) error
	GetUserByID(ID uint64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateProfile(profile models.Profile) error
	GetProfileByID(id uint64) (*models.Profile, error)
}
