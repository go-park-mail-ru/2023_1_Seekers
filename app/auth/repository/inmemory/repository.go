package inmemory

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"time"
)

type AuthDB struct {
	Cookies []model.Cookie
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
	if _, err := a.GetCookie(cookie.Session); err != nil {
		return fmt.Errorf("cant create cookie: %w", err)
	}
	a.Cookies = append(a.Cookies, cookie)
	return nil
}

func (a *AuthDB) DeleteCookie(session string) error {
	for i, c := range a.Cookies {
		if c.Session == session {
			a.Cookies = append(a.Cookies[:i], a.Cookies[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("cant delete cookie %s, not found", session)
}

func (a *AuthDB) GetCookie(session string) (*model.Cookie, error) {
	for _, c := range a.Cookies {
		if c.Session == session {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("cant get cookie %s", session)
}
