package auth

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

type UseCaseI interface {
	SignIn(form *models.FormLogin) (*models.AuthResponse, *models.Session, error)
	SignUp(form *models.FormSignUp) (*models.AuthResponse, *models.Session, error)
	CreateSession(uID uint64) (*models.Session, error)
	DeleteSession(sessionID string) error
	GetSession(sessionID string) (*models.Session, error)
}
