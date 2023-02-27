package auth

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"net/http"
)

type UseCase interface {
	SignIn(form model.FormLogin) (*model.User, error)
	SignUp(form model.FormSignUp) (*model.User, error)
	Logout(sessionId string) error
	Auth(sessionId string) error
}

type Repo interface {
	Create(user model.User) (*model.User, error)
	Delete(user model.User) error
	GetById(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type Handlers interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
}
