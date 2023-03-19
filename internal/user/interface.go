package user

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"net/http"
)

type HandlersI interface {
	Delete(w http.ResponseWriter, r *http.Request)
	GetInfo(w http.ResponseWriter, r *http.Request)
	EditInfo(w http.ResponseWriter, r *http.Request)
	EditPw(w http.ResponseWriter, r *http.Request)
	EditAvatar(w http.ResponseWriter, r *http.Request)
	GetAvatar(w http.ResponseWriter, r *http.Request)
}

type UseCaseI interface {
	Create(user models.User) (*models.User, error)
	Delete(user models.User) error
	GetByID(ID uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	//TODO тут какая-то модель для изменения данных
	EditInfo(user models.User) (*models.User, error)
	EditPw(ID uint64, newPW string) (*models.User, error)
	EditAvatar(ID uint64, newAvatar *models.Image) error
	GetAvatar(email string) (*models.Image, error)
}

type RepoI interface {
	Create(user models.User) (*models.User, error)
	Delete(user models.User) error
	GetByID(ID uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	SetAvatar(uID uint64, avatar string) error
}
