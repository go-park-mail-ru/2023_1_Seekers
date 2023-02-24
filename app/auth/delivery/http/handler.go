package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/utils"
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"time"
)

type handlers struct {
	useCase auth.UseCase
}

func New(aUC auth.UseCase) auth.Handlers {
	return &handlers{
		useCase: aUC,
	}
}

func (h *handlers) SignUp(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		c.Logger().Error("sign up not post method error")
		return utils.SendError(c, http.StatusMethodNotAllowed, "sign up not post method error")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger().Error("failed to close request: %w", err)
		}
	}(c.Request().Body)
	form := model.FormReg{}

	err := json.NewDecoder(c.Request().Body).Decode(&form)
	if err != nil {
		c.Logger().Error(fmt.Errorf("faliled decode sign up form %w", err))
		return utils.SendError(c, http.StatusBadRequest, "bad request")
	}
	_, cookie, err := h.useCase.SignUp(form)
	if err != nil {
		c.Logger().Error(fmt.Errorf("faliled to sign up %w", err))
		return utils.SendError(c, http.StatusInternalServerError, fmt.Sprintf("faliled to sign up %v", err.Error()))
	}

	http.SetCookie(c.Response(), &http.Cookie{
		Name:    config.CookieName,
		Value:   cookie.Session,
		Expires: cookie.Expire,
	})
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func (h *handlers) SignIn(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		c.Logger().Error("sign in not post method error")
		return utils.SendError(c, http.StatusMethodNotAllowed, "sign in not post method error")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger().Error("failed to close request: %w", err)
		}
	}(c.Request().Body)
	form := model.FormAuth{}

	err := json.NewDecoder(c.Request().Body).Decode(&form)
	if err != nil {
		c.Logger().Error(fmt.Errorf("faliled decode sign in form %w", err))
		return utils.SendError(c, http.StatusBadRequest, "bad request")
	}
	_, cookie, err := h.useCase.SignIn(form)
	if err != nil {
		c.Logger().Error(fmt.Errorf("faliled to sign in %w", err))
		return utils.SendError(c, http.StatusBadRequest, fmt.Sprintf("faliled to sign in %v", err.Error()))
	}

	http.SetCookie(c.Response(), &http.Cookie{
		Name:    config.CookieName,
		Value:   cookie.Session,
		Expires: cookie.Expire,
	})
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func (h *handlers) Logout(c echo.Context) error {
	cookie, err := c.Request().Cookie(config.CookieName)
	if err == http.ErrNoCookie {
		c.Logger().Error(fmt.Errorf("faliled logout %w", err))
		return utils.SendError(c, http.StatusUnauthorized, fmt.Sprintf("faliled to logout %v", err.Error()))
	} else if err != nil {
		c.Logger().Error(fmt.Errorf("faliled logout %w", err))
		return utils.SendError(c, http.StatusBadRequest, fmt.Sprintf("faliled to logout %v", err.Error()))
	}

	err = h.useCase.DeleteCookie(cookie.Value)
	if err != nil {
		c.Logger().Error(fmt.Errorf("faliled logout %w", err))
		return utils.SendError(c, http.StatusUnauthorized, fmt.Sprintf("faliled to logout %v", err.Error()))
	}

	http.SetCookie(c.Response(), &http.Cookie{
		Name:    config.CookieName,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
	})

	c.Response().WriteHeader(http.StatusOK)
	return nil
}
