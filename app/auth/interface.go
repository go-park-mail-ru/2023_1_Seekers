package auth

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/labstack/echo/v4"
)

type UseCase interface {
	SignIn(form model.FormAuth) (*model.User, *model.Session, error)
	SignUp(form model.FormReg) (*model.User, *model.Session, error)
	Logout(sessionId string) error
}

type Handlers interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
	Logout(c echo.Context) error
}
