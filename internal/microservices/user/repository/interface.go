package repository

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

//go:generate mockgen -destination=./mocks/mockrepo.go -source=./interface.go -package=mocks

type UserRepoI interface {
	Create(user *models.User) (*models.User, error)
	EditInfo(ID uint64, info *models.UserInfo) error
	Delete(ID uint64) error
	GetByID(ID uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	SetAvatar(ID uint64, avatar string) error
	EditPw(ID uint64, newPW string) error
	GetInfoByID(ID uint64) (*models.UserInfo, error)
	GetInfoByEmail(email string) (*models.UserInfo, error)
	IsCustomAvatar(ID uint64) (bool, error)
	SetCustomAvatar(ID uint64) error
}
