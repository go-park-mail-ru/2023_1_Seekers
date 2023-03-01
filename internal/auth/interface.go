package auth

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"net/http"
)

type UseCase interface {
	SignIn(form models.FormLogin) (*models.User, error)
	SignUp(form models.FormSignUp) (*models.User, error)
	CreateSession(uID uint64) (*models.Session, error)
	DeleteSession(sessionID string) error
	DeleteSessionByUID(uID uint64) error
	GetSession(sessionID string) (*models.Session, error)
	GetSessionByUID(uID uint64) (*models.Session, error)
}

type Repo interface {
	Create(user models.User) (*models.User, error)
	Delete(user models.User) error
	GetByID(ID uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	CreateSession(session models.Session) error
	DeleteSession(sessionID string) error
	DeleteSessionByUID(uID uint64) error
	GetSession(sessionID string) (*models.Session, error)
	GetSessionByUID(uID uint64) (*models.Session, error)
}

type Handlers interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}
