package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"time"
)

type AuthDB struct {
	Users []model.Cookie
}

func New() auth.Repo {
	return &AuthDB{
		[]model.Cookie{
			{1, "session_cookie_1", time.Now().Add(time.Hour * 24)},
			{2, "session_cookie_2", time.Now().Add(time.Hour * 24)},
			{3, "session_cookie_3", time.Now().Add(time.Hour * 24)},
		},
	}
}

func (a *AuthDB) NewCookie(cookie model.Cookie) error {
	return nil
}

func (a *AuthDB) DeleteCookie(value string) error {
	return nil
}

func (a *AuthDB) GetCookie(value string) (*model.Cookie, error) {
	return nil, nil
}
