package session

import "github.com/go-park-mail-ru/2023_1_Seekers/app/model"

type UseCase interface {
	Create(uId int) (*model.Session, error)
	Delete(sessionId string) error
	GetSession(sessionId string) (*model.Session, error)
	GetSessionById(uId int) (*model.Session, error)
}

type Repo interface {
	Create(session model.Session) error
	Delete(sessionId string) error
	GetSession(sessionId string) (*model.Session, error)
	GetSessionById(uId int) (*model.Session, error)
}
