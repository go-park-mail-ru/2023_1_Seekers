package http

import (
	"encoding/json"
	"errors"
	"fmt"
	auth "github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	_session "github.com/go-park-mail-ru/2023_1_Seekers/app/session"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/utils"
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"time"
)

type handlers struct {
	authUC    auth.UseCase
	sessionUC _session.UseCase
}

func New(aUC auth.UseCase, sUC _session.UseCase) auth.Handlers {
	return &handlers{
		authUC:    aUC,
		sessionUC: sUC,
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

	form := model.FormSignUp{}
	err := json.NewDecoder(r.Body).Decode(&form)
	fmt.Println(form)
	if err != nil {
		log.Error(fmt.Errorf("faliled decode sign up form %w", err))
		utils.SendError(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}
	user, err := h.authUC.SignUp(form)
	if err != nil {
		log.Error(fmt.Errorf("faliled to sign up %w", err))
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to sign up %v", err.Error()))
		return
	}

	session, err := h.sessionUC.Create(user.Id)
	if err != nil {
		log.Error(fmt.Errorf("faliled to sign up %w", err))
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to sign up %v", err.Error()))

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
	form := model.FormLogin{}

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		log.Error(fmt.Errorf("faliled decode sign in form %w", err))
		utils.SendError(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}
	user, err := h.authUC.SignIn(form)
	if err != nil {
		log.Error(fmt.Errorf("faliled to sign in %w", err))
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to sign in %v", err.Error()))
		return
	}
	// когда логинимся, то обновляем куку, если ранее была, то удалится и пересоздастся
	err = h.sessionUC.DeleteByUId(user.Id)
	fmt.Println(err)
	session, err := h.sessionUC.Create(user.Id)
	if err != nil {
		log.Error(fmt.Errorf("faliled to sign in %w", err))
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to sign in %w", err))
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

	err = h.sessionUC.Delete(cookie.Value)
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

func (h *handlers) Auth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.CookieName)
	if err != nil {
		log.Error(fmt.Errorf("faliled auth %w", err))
		utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("faliled auth %w", err))
		return
	}

	_, err = h.sessionUC.GetSession(cookie.Value)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("failed to auth: %w", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
