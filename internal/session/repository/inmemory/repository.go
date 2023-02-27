package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/session"
)

type sessionDB struct {
	sessions []model.Session
}

func New() session.Repo {
	return &sessionDB{
		[]model.Session{
			//уже есть сессия для Uid 1
			{1, "randgeneratedcookie12334524524523542"},
		},
	}
}

func (sDb *sessionDB) Create(s model.Session) error {
	if _, err := sDb.GetSessionByUId(s.UId); err == nil {
		return session.ErrSessionExists
	}
	sDb.sessions = append(sDb.sessions, s)
	return nil
}

func (sDb *sessionDB) Delete(sessionId string) error {
	for i, s := range sDb.sessions {
		if s.SessionId == sessionId {
			sDb.sessions = append(sDb.sessions[:i], sDb.sessions[i+1:]...)
			return nil
		}
	}

	return session.ErrSessionNotFound
}

func (sDb *sessionDB) DeleteByUId(uId uint64) error {
	for i, s := range sDb.sessions {
		if s.UId == uId {
			sDb.sessions = append(sDb.sessions[:i], sDb.sessions[i+1:]...)
			return nil
		}
	}
	return session.ErrSessionNotFound
}

func (sDb *sessionDB) GetSession(sessionId string) (*model.Session, error) {
	for _, s := range sDb.sessions {
		if s.SessionId == sessionId {
			return &s, nil
		}
	}
	return nil, session.ErrSessionNotFound
}

func (sDb *sessionDB) GetSessionByUId(uId uint64) (*model.Session, error) {
	for _, s := range sDb.sessions {
		if s.UId == uId {
			return &s, nil
		}
	}
	return nil, session.ErrSessionNotFound
}
