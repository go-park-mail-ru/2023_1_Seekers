package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
)

type sessionsDB struct {
	sessions []models.Session
}

func New() auth.RepoI {
	return &sessionsDB{}
}

func (sDb *sessionsDB) CreateSession(s models.Session) error {
	if _, err := sDb.GetSessionByUID(s.UID); err == nil {
		return auth.ErrSessionExists
	}
	sDb.sessions = append(sDb.sessions, s)
	return nil
}

func (sDb *sessionsDB) DeleteSession(sessionID string) error {
	for i, s := range sDb.sessions {
		if s.SessionID == sessionID {
			sDb.sessions = append(sDb.sessions[:i], sDb.sessions[i+1:]...)
			return nil
		}
	}

	return auth.ErrSessionNotFound
}

func (sDb *sessionsDB) DeleteSessionByUID(uID uint64) error {
	for i, s := range sDb.sessions {
		if s.UID == uID {
			sDb.sessions = append(sDb.sessions[:i], sDb.sessions[i+1:]...)
			return nil
		}
	}
	return auth.ErrSessionNotFound
}

func (sDb *sessionsDB) GetSession(sessionID string) (*models.Session, error) {
	for _, s := range sDb.sessions {
		if s.SessionID == sessionID {
			return &s, nil
		}
	}
	return nil, auth.ErrSessionNotFound
}

func (sDb *sessionsDB) GetSessionByUID(uID uint64) (*models.Session, error) {
	for _, s := range sDb.sessions {
		if s.UID == uID {
			return &s, nil
		}
	}
	return nil, auth.ErrSessionNotFound
}
