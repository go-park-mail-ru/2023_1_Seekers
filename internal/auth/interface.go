package auth

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
	"net/http"
)

type UseCase interface {
	SignIn(form model.FormLogin) (*model.User, error)
	SignUp(form model.FormSignUp) (*model.User, error)
	CreateSession(uId uint64) (*model.Session, error)
	DeleteSession(sessionId string) error
	DeleteSessionByUId(uId uint64) error
	GetSession(sessionId string) (*model.Session, error)
	GetSessionByUId(uId uint64) (*model.Session, error)
}

type Repo interface {
	Create(user model.User) (*model.User, error)
	Delete(user model.User) error
	GetById(id uint64) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	CreateSession(session model.Session) error
	DeleteSession(sessionId string) error
	DeleteSessionByUId(uId uint64) error
	GetSession(sessionId string) (*model.Session, error)
	GetSessionByUId(uId uint64) (*model.Session, error)
}

type Handlers interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request) // del
}
