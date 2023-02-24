package auth

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/labstack/echo/v4"
)

type UseCase interface {
	SignIn(form model.FormAuth) (*model.User, *model.Cookie, error)
	SignUp(form model.FormReg) (*model.User, *model.Cookie, error)
	NewCookie(uId int) (*model.Cookie, error)
	DeleteCookie(session string) error
	GetCookie(session string) (*model.Cookie, error)
}

type Repo interface {
	NewCookie(cookie model.Cookie) error
	DeleteCookie(session string) error
	GetCookie(session string) (*model.Cookie, error)
}

type Handlers interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
	Logout(c echo.Context) error
}
