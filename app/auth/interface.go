package auth

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"net/http"
)

type UseCase interface {
	SignIn(form model.FormAuth) (*model.User, *model.Session, error)
	SignUp(form model.FormReg) (*model.User, *model.Session, error)
	Logout(sessionId string) error
	Auth(sessionId string) error
}

type Handlers interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
}
