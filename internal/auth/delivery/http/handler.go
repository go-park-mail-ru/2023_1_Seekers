package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type handlers struct {
	authUC auth.UseCase
	userUC _user.UseCase
}

func New(aUC auth.UseCase, uUC _user.UseCase) auth.Handlers {
	return &handlers{
		authUC: aUC,
		userUC: uUC,
	}
}

func (h *handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host, r.Header, r.Body)
	if r.Method != http.MethodPost {
		authErr := errors.New(auth.AuthErrors[auth.ErrInvalidMethodPost], auth.ErrInvalidMethodPost)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// ?
			log.Error("failed to close request: ", err)
		}
	}(r.Body)

	form := model.FormSignUp{}
	err := json.NewDecoder(r.Body).Decode(&form)
	fmt.Println(form)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrInvalidForm], auth.ErrInvalidForm.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}
	user, err := h.authUC.SignUp(form)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedSignUp], auth.ErrFailedSignUp.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	profile := model.Profile{
		UId:       user.Id,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		BirthDate: form.BirthDate,
	}
	err = h.userUC.CreateProfile(profile)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedCreateProfile], auth.ErrFailedCreateProfile.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	session, err := h.authUC.CreateSession(user.Id)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedCreateSession], auth.ErrFailedCreateSession.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    config.CookieName,
		Value:   session.SessionId,
		Expires: time.Now().Add(config.CookieTTL),
	})
	pkg.SendJson(w, http.StatusOK, user)
}

func (h *handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host, r.Header, r.Body)
	if r.Method != http.MethodPost {
		authErr := errors.New(auth.AuthErrors[auth.ErrInvalidMethodPost], auth.ErrInvalidMethodPost)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(fmt.Errorf("failed to close request: %w", err))
		}
	}(r.Body)
	form := model.FormLogin{}

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrInvalidForm], auth.ErrInvalidForm.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	user, err := h.authUC.SignIn(form)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedSignIn], auth.ErrFailedSignIn.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	// когда логинимся, то обновляем куку, если ранее была, то удалится и пересоздастся
	err = h.authUC.DeleteSessionByUId(user.Id)
	session, err := h.authUC.CreateSession(user.Id)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedCreateSession], auth.ErrFailedCreateSession.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    config.CookieName,
		Value:   session.SessionId,
		Expires: time.Now().Add(config.CookieTTL),
	})
	pkg.SendJson(w, http.StatusOK, user)
}

func (h *handlers) Logout(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host, r.Header, r.Body)
	cookie, err := r.Cookie(config.CookieName)
	if err == http.ErrNoCookie {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedLogoutNoCookie], auth.ErrFailedLogoutNoCookie.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	} else if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedLogout], auth.ErrFailedLogout.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
		return
	}

	err = h.authUC.DeleteSession(cookie.Value)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedLogout], auth.ErrFailedLogout.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    config.CookieName,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
	})

	w.WriteHeader(http.StatusOK)
}

func (h *handlers) Auth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.CookieName)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedAuth], auth.ErrFailedAuth.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
	}

	_, err = h.authUC.GetSession(cookie.Value)
	if err != nil {
		authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedGetSession], auth.ErrFailedGetSession.Error(), err)
		log.Error(authErr)
		pkg.SendError(w, authErr)
	}

	w.WriteHeader(http.StatusOK)
}
