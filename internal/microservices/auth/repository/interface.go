package repository

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

//go:generate mockgen -destination=./mocks/mockrepo.go -source=./interface.go -package=mocks

type SessionRepoI interface {
	CreateSession(uID uint64) (*models.Session, error)
	DeleteSession(sessionID string) error
	GetSession(sessionID string) (*models.Session, error)
}
