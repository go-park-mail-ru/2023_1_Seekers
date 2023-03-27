package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	_ "github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-playground/validator/v10"
	pkgErrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type handlers struct {
	authUC auth.UseCaseI
}

func New(aUC auth.UseCaseI) auth.HandlersI {
	return &handlers{
		authUC: aUC,
	}
}

func setNewCookie(w http.ResponseWriter, session *models.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.CookieName,
		Value:    session.SessionID,
		Expires:  time.Now().Add(config.CookieTTL),
		HttpOnly: true,
		Path:     config.CookiePath,
	})
}

func delCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    config.CookieName,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    config.CookiePath,
	})
}

// SignUp godoc
// @Summary      SignUp
// @Description  user sign up
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Param    user body models.FormSignUp true "user info"
// @Success  200 {object} models.SignUpResponse "user created"
// @Failure 401 {object} errors.JSONError "passwords dont match"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 403 {object} errors.JSONError "password too short"
// @Failure 409 {object} errors.JSONError "user already exists"
// @Failure 500 {object} errors.JSONError "failed to create profile"
// @Failure 500 {object} errors.JSONError "failed to create session"
// @Router   /signup [post]
func (h *handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("failed to close request: ", err)
		}
	}(r.Body)

	form := models.FormSignUp{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	response, session, err := h.authUC.SignUp(form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	setNewCookie(w, session)
	pkg.SendJSON(w, r, http.StatusOK, response)
}

// SignIn godoc
// @Summary      SignIn
// @Description  user sign in
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Param    user body models.FormLogin true "user info"
// @Success  200 {object} models.SignInResponse "user created"
// @Failure 401 {object} errors.JSONError "wrong password"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 500 {object} errors.JSONError "failed to create session"
// @Router   /signin [post]
func (h *handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(fmt.Errorf("failed to close request: %w", err))
		}
	}(r.Body)
	form := models.FormLogin{}

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	response, session, err := h.authUC.SignIn(form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	setNewCookie(w, session)
	pkg.SendJSON(w, r, http.StatusOK, response)
}

// Logout godoc
// @Summary      Logout
// @Description  user log out
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Success  200 "success logout"
// @Failure 401 {object} errors.JSONError "failed auth"
// @Failure 401 {object} errors.JSONError "failed get session"
// @Router   /logout [post]
func (h *handlers) Logout(w http.ResponseWriter, _ *http.Request) {
	delCookie(w)
	w.WriteHeader(http.StatusOK)
}
