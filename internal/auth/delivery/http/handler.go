package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type handlers struct {
	authUC auth.UseCaseI
	userUC _user.UseCaseI
	mailUC mail.UseCaseI
}

func New(aUC auth.UseCaseI, uUC _user.UseCaseI) auth.HandlersI {
	return &handlers{
		authUC: aUC,
		userUC: uUC,
	}
}

// SignUp godoc
// @Summary      SignUp
// @Description  user sign up
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Param    user body models.FormSignUp true "user info"
// @Success  200 {object} models.User "user created"
// @Failure 401 {object} error "passwords dont match"
// @Failure 403 {object} error "invalid form"
// @Failure 403 {object} error "password too short"
// @Failure 409 {object} error "user already exists"
// @Failure 500 {object} error "failed to create profile"
// @Failure 500 {object} error "failed to create session"
// @Router   /signup [post]
func (h *handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host, r.Header, r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("failed to close request: ", err)
		}
	}(r.Body)

	form := models.FormSignUp{}
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		authErr := errors.New(auth.Errors[auth.ErrInvalidForm], auth.ErrInvalidForm)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	validate := validator.New()
	err = validate.Struct(form)
	if err != nil {
		authErr := errors.New(auth.Errors[auth.ErrInvalidForm], err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	user, err := h.authUC.SignUp(form)
	if err != nil {
		if err == auth.ErrPwDontMatch {
			authErr := errors.New(auth.Errors[auth.ErrPwDontMatch], auth.ErrPwDontMatch)
			log.Error(authErr)
			pkg.SendError(w, authErr)
			return
		}
		if err == _user.ErrTooShortPw {
			authErr := errors.New(auth.Errors[_user.ErrTooShortPw], _user.ErrTooShortPw)
			log.Error(authErr)
			pkg.SendError(w, authErr)
			return
		}
		authErr := errors.New(auth.Errors[auth.ErrUserExists], auth.ErrUserExists)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	profile := models.Profile{
		UID:       user.ID,
		FirstName: form.FirstName,
		LastName:  form.LastName,
	}
	err = h.userUC.CreateProfile(profile)
	if err != nil {
		authErr := errors.New(auth.Errors[auth.ErrFailedCreateProfile], auth.ErrFailedCreateProfile)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	session, err := h.authUC.CreateSession(user.ID)
	if err != nil {
		authErr := errors.New(auth.Errors[auth.ErrFailedCreateSession], auth.ErrFailedCreateSession)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	//err = h.mailUC.CreateHelloMessage(user.ID)
	//if err != nil {
	//	authErr := errors.New(auth.Errors[auth.ErrInternalHelloMsg], auth.ErrInternalHelloMsg)
	//	log.Error(authErr)
	//	pkg.SendError(w, authErr)
	//	return
	//}

	http.SetCookie(w, &http.Cookie{
		Name:     config.CookieName,
		Value:    session.SessionID,
		Expires:  time.Now().Add(config.CookieTTL),
		HttpOnly: true,
		Path:     config.CookiePath,
	})
	pkg.SendJSON(w, http.StatusOK, user)
}

// SignIn godoc
// @Summary      SignIn
// @Description  user sign in
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Param    user body models.FormLogin true "user info"
// @Success  200 {object} models.User "user created"
// @Failure 401 {object} error "wrong password"
// @Failure 403 {object} error "invalid form"
// @Failure 500 {object} error "failed to create session"
// @Router   /signin [post]
func (h *handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host, r.Header, r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(fmt.Errorf("failed to close request: %w", err))
		}
	}(r.Body)
	form := models.FormLogin{}

	// TODO validate !!!

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		authErr := errors.New(auth.Errors[auth.ErrInvalidForm], auth.ErrInvalidForm)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	user, err := h.authUC.SignIn(form)
	if err != nil {
		authErr := errors.New(auth.Errors[auth.ErrWrongPw], auth.ErrWrongPw)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	// когда логинимся, то обновляем куку, если ранее была, то удалится и пересоздастся
	err = h.authUC.DeleteSessionByUID(user.ID)
	session, err := h.authUC.CreateSession(user.ID)
	if err != nil {
		authErr := errors.New(auth.Errors[auth.ErrFailedCreateSession], auth.ErrFailedCreateSession)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     config.CookieName,
		Value:    session.SessionID,
		Expires:  time.Now().Add(config.CookieTTL),
		HttpOnly: true,
		Path:     config.CookiePath,
	})
	pkg.SendJSON(w, http.StatusOK, user)
}

// Logout godoc
// @Summary      Logout
// @Description  user log out
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Success  200 "success logout"
// @Failure 401 {object} error "failed auth"
// @Failure 401 {object} error "failed get session"
// @Router   /logout [post]
func (h *handlers) Logout(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host, r.Header, r.Body)

	http.SetCookie(w, &http.Cookie{
		Name:    config.CookieName,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    config.CookiePath,
	})

	w.WriteHeader(http.StatusOK)
}
