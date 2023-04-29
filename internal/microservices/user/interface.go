package user

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

type UseCaseI interface {
	Create(user *models.User) (*models.User, error)
	Delete(ID uint64) error
	GetByID(ID uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetInfo(ID uint64) (*models.UserInfo, error)
	GetInfoByEmail(email string) (*models.UserInfo, error)
	EditInfo(ID uint64, info *models.UserInfo) (*models.UserInfo, error)
	EditAvatar(ID uint64, newAvatar *models.Image, isCustom bool) error
	GetAvatar(email string) (*models.Image, error)
	EditPw(ID uint64, form *models.EditPasswordRequest) error
}
