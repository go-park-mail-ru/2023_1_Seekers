package user

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"net/http"
)

type HandlersI interface {
	Delete(w http.ResponseWriter, r *http.Request)
	GetInfo(w http.ResponseWriter, r *http.Request)
	GetPersonalInfo(w http.ResponseWriter, r *http.Request)
	EditInfo(w http.ResponseWriter, r *http.Request)
	EditPw(w http.ResponseWriter, r *http.Request)
	EditAvatar(w http.ResponseWriter, r *http.Request)
	GetAvatar(w http.ResponseWriter, r *http.Request)
}

type UseCaseI interface {
	Create(user models.User) (*models.User, error)
	Delete(ID uint64) error
	GetByID(ID uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetInfo(ID uint64) (*models.UserInfo, error)
	EditInfo(ID uint64, info models.UserInfo) (*models.UserInfo, error)
	EditPw(ID uint64, password models.EditPasswordRequest) error
	EditAvatar(ID uint64, newAvatar *models.Image) error
	GetAvatar(email string) (*models.Image, error)
}

type RepoI interface {
	Create(user models.User) (*models.User, error)
	EditInfo(ID uint64, info models.UserInfo) error
	Delete(ID uint64) error
	GetByID(ID uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	SetAvatar(ID uint64, avatar string) error
	EditPw(ID uint64, newPW string) error
	GetInfoByID(ID uint64) (*models.UserInfo, error)
	GetInfoByEmail(email string) (*models.UserInfo, error)
}
