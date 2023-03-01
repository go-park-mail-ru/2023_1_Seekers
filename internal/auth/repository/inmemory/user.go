package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
)

type useruDb struct {
	users    []models.User
	sessions []models.Session
}

func New() auth.Repo {
	return &useruDb{
		[]models.User{
			{1, "test@example.com", "1234"},
			{2, "gena@example.com", "4321"},
			{3, "max@example.com", "1379"},
		},
		[]models.Session{
			{1, "randgeneratedcookie12334524524523542"}, //уже есть сессия для Uid 1

		},
	}
}

func (uDb *useruDb) GetByID(id uint64) (*models.User, error) {
	for i, u := range uDb.users {
		if u.ID == id {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *useruDb) GetByEmail(email string) (*models.User, error) {
	for i, u := range uDb.users {
		if u.Email == email {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *useruDb) Create(user models.User) (*models.User, error) {
	_, err := uDb.GetByEmail(user.Email)
	if err == nil {
		return nil, auth.ErrUserExists
	}
	//слой бд отвечает за присваивание id
	// TODO hash pw
	user.ID = uint64(len(uDb.users) + 1)
	uDb.users = append(uDb.users, user)
	return &user, nil
}

func (uDb *useruDb) Delete(user models.User) error {
	for i, u := range uDb.users {
		if u.ID == user.ID {
			uDb.users = append(uDb.users[:i], uDb.users[i+1:]...)
			return nil
		}
	}
	return auth.ErrUserNotFound
}

func (uDb *useruDb) CreateSession(s models.Session) error {
	if _, err := uDb.GetSessionByUID(s.UID); err == nil {
		return auth.ErrSessionExists
	}
	uDb.sessions = append(uDb.sessions, s)
	return nil
}

func (uDb *useruDb) DeleteSession(sessionID string) error {
	for i, s := range uDb.sessions {
		if s.SessionID == sessionID {
			uDb.sessions = append(uDb.sessions[:i], uDb.sessions[i+1:]...)
			return nil
		}
	}

	return auth.ErrSessionNotFound
}

func (uDb *useruDb) DeleteSessionByUID(uID uint64) error {
	for i, s := range uDb.sessions {
		if s.UID == uID {
			uDb.sessions = append(uDb.sessions[:i], uDb.sessions[i+1:]...)
			return nil
		}
	}
	return auth.ErrSessionNotFound
}

func (uDb *useruDb) GetSession(sessionID string) (*models.Session, error) {
	for _, s := range uDb.sessions {
		if s.SessionID == sessionID {
			return &s, nil
		}
	}
	return nil, auth.ErrSessionNotFound
}

func (uDb *useruDb) GetSessionByUID(uID uint64) (*models.Session, error) {
	for _, s := range uDb.sessions {
		if s.UID == uID {
			return &s, nil
		}
	}
	return nil, auth.ErrSessionNotFound
}
