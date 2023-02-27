package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
)

type useruDb struct {
	users    []model.User
	sessions []model.Session
}

func New() auth.Repo {
	return &useruDb{
		[]model.User{
			{1, "test@example.com", "1234"},
			{2, "gena@example.com", "4321"},
			{3, "max@example.com", "1379"},
		},
		[]model.Session{
			//уже есть сессия для Uid 1
			{1, "randgeneratedcookie12334524524523542"},
		},
	}
}

func (uDb *useruDb) GetById(id uint64) (*model.User, error) {
	for i, u := range uDb.users {
		if u.Id == id {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *useruDb) GetByEmail(email string) (*model.User, error) {
	for i, u := range uDb.users {
		if u.Email == email {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *useruDb) Create(user model.User) (*model.User, error) {
	_, err := uDb.GetById(user.Id)
	if err == nil {
		return nil, auth.ErrUserExists
	}
	_, err = uDb.GetByEmail(user.Email)
	if err == nil {
		return nil, auth.ErrUserExists
	}
	//слой бд отвечает за присваивание id
	// TODO hash pw
	user.Id = uint64(len(uDb.users) + 1)
	uDb.users = append(uDb.users, user)
	return &user, nil
}

func (uDb *useruDb) Delete(user model.User) error {
	for i, u := range uDb.users {
		if u.Id == user.Id {
			uDb.users = append(uDb.users[:i], uDb.users[i+1:]...)
			return nil
		}
	}
	return auth.ErrUserNotFound
}

func (uDb *useruDb) CreateSession(s model.Session) error {
	if _, err := uDb.GetSessionByUId(s.UId); err == nil {
		return auth.ErrSessionExists
	}
	uDb.sessions = append(uDb.sessions, s)
	return nil
}

func (uDb *useruDb) DeleteSession(sessionId string) error {
	for i, s := range uDb.sessions {
		if s.SessionId == sessionId {
			uDb.sessions = append(uDb.sessions[:i], uDb.sessions[i+1:]...)
			return nil
		}
	}

	return auth.ErrSessionNotFound
}

func (uDb *useruDb) DeleteSessionByUId(uId uint64) error {
	for i, s := range uDb.sessions {
		if s.UId == uId {
			uDb.sessions = append(uDb.sessions[:i], uDb.sessions[i+1:]...)
			return nil
		}
	}
	return auth.ErrSessionNotFound
}

func (uDb *useruDb) GetSession(sessionId string) (*model.Session, error) {
	for _, s := range uDb.sessions {
		if s.SessionId == sessionId {
			return &s, nil
		}
	}
	return nil, auth.ErrSessionNotFound
}

func (uDb *useruDb) GetSessionByUId(uId uint64) (*model.Session, error) {
	for _, s := range uDb.sessions {
		if s.UId == uId {
			return &s, nil
		}
	}
	return nil, auth.ErrSessionNotFound
}
