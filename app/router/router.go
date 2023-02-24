package router

import (
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/usecase"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/app/user/usecase"
	//userHandler "github.com/go-park-mail-ru/2023_1_Seekers/app/user/delivery/http"
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/delivery/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func Register(e *echo.Echo) {
	userRepo := _userRepo.New()
	authRepo := _authRepo.New()

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo, usersUCase)

	//userH := userHandler.
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
