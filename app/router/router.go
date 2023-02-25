package router

import (
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/usecase"
	_sessionRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/session/repository/inmemory"
	_sessionUcase "github.com/go-park-mail-ru/2023_1_Seekers/app/session/usecase"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/app/user/usecase"
	//userHandler "github.com/go-park-mail-ru/2023_1_Seekers/app/user/delivery/http"
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/delivery/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// TODO перейти c echo на net/http

func Register(e *echo.Echo) {
	userRepo := _userRepo.New()
	sessionRepo := _sessionRepo.New()

	usersUCase := _userUCase.New(userRepo)
	sessionUCase := _sessionUcase.New(sessionRepo)
	authUCase := _authUCase.New(sessionUCase, usersUCase)

	authH := _authHandler.New(authUCase)

	api := e.Group("/api")
	api.POST("/signin", authH.SignIn)
	api.POST("/signup", authH.SignUp)
	api.GET("/logout", authH.Logout)
}

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	// TODO ? e.Validator =
	return e
}
