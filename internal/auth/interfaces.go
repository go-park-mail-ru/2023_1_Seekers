package auth

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"net/http"
)

type HandlersI interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	EditPw(w http.ResponseWriter, r *http.Request)
	GetCSRF(w http.ResponseWriter, r *http.Request)
}

type UseCaseI interface {
	SignIn(form models.FormLogin) (*models.AuthResponse, *models.Session, error)
	SignUp(form models.FormSignUp) (*models.AuthResponse, *models.Session, error)
	EditPw(ID uint64, pw models.EditPasswordRequest) error
}

type SessionUseCaseI interface {
	CreateSession(uID uint64) (*models.Session, error)
	DeleteSession(sessionID string) error
	GetSession(sessionID string) (*models.Session, error)
}

type SessionRepoI interface {
	CreateSession(uID uint64) (*models.Session, error)
	DeleteSession(sessionID string) error
	GetSession(sessionID string) (*models.Session, error)
}
