package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/utils"
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
	"github.com/labstack/gommon/log"
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

func (h *handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host, r.Header, r.Body)
	if r.Method != http.MethodPost {
		// TODO pretty logger
		log.Error("sign up not post method error")
		utils.SendError(w, http.StatusMethodNotAllowed, errors.New("sign up not post method error"))
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("failed to close request: %w", err)
		}
	}(r.Body)
	form := model.FormReg{}

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		log.Error(fmt.Errorf("faliled decode sign up form %w", err))
		utils.SendError(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}
	user, session, err := h.useCase.SignUp(form)
	if err != nil {
		log.Error(fmt.Errorf("faliled to sign up %w", err))
		utils.SendError(w, http.StatusInternalServerError, fmt.Errorf("faliled to sign up %v", err.Error()))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    config.CookieName,
		Value:   session.SessionId,
		Expires: time.Now().Add(config.CookieTTL),
	})
	utils.SendJson(w, http.StatusOK, user)
}

func (h *handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host, r.Header, r.Body)
	if r.Method != http.MethodPost {
		log.Error("sign in not post method error")
		utils.SendError(w, http.StatusMethodNotAllowed, errors.New("sign in not post method error"))
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("failed to close request: %w", err)
		}
	}(r.Body)
	form := model.FormAuth{}

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		log.Error(fmt.Errorf("faliled decode sign in form %w", err))
		utils.SendError(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}
	user, session, err := h.useCase.SignIn(form)
	if err != nil {
		log.Error(fmt.Errorf("faliled to sign in %w", err))
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to sign in %v", err.Error()))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    config.CookieName,
		Value:   session.SessionId,
		Expires: time.Now().Add(config.CookieTTL),
	})
	utils.SendJson(w, http.StatusOK, user)
}

func (h *handlers) Logout(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host, r.Header, r.Body)
	cookie, err := r.Cookie(config.CookieName)
	if err == http.ErrNoCookie {
		log.Error(fmt.Errorf("faliled logout %w", err))
		utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("faliled to logout %v", err.Error()))
		return
		//return
	} else if err != nil {
		log.Error(fmt.Errorf("faliled logout %w", err))
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to logout %v", err.Error()))
		return
	}

	err = h.useCase.Logout(cookie.Value)
	if err != nil {
		log.Error(fmt.Errorf("faliled logout %w", err))
		utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("faliled to logout %v", err.Error()))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    config.CookieName,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
	})

	w.WriteHeader(http.StatusOK)
}
