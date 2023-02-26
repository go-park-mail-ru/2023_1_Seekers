package inmemory

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/session"
)

type sessionDB struct {
	sessions []model.Session
}

func New() session.Repo {
	return &sessionDB{
		[]model.Session{
			{1, "randgeneratedcookie12334524524523542"},
		},
	}
}

func (sDb *sessionDB) Create(session model.Session) error {
	if _, err := sDb.GetSessionById(session.UId); err == nil {
		return fmt.Errorf("cant create session: %w", err)
	}
	sDb.sessions = append(sDb.sessions, session)
	return nil
}

func (sDb *sessionDB) Delete(sessionId string) error {
	for i, s := range sDb.sessions {
		if s.SessionId == sessionId {
			sDb.sessions = append(sDb.sessions[:i], sDb.sessions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("cant delete session: no session with id %s", sessionId)
}

func (sDb *sessionDB) GetSession(sessionId string) (*model.Session, error) {
	for _, s := range sDb.sessions {
		if s.SessionId == sessionId {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("no session %s", sessionId)
}

func (sDb *sessionDB) GetSessionById(uId int) (*model.Session, error) {
	for _, s := range sDb.sessions {
		if s.UId == uId {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("cant get session: no user with id %d", uId)
}
