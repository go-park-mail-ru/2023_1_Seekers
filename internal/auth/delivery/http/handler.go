package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
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
		log.Error(auth.ErrInvalidMethodPost)
		pkg.SendError(w, http.StatusMethodNotAllowed, auth.ErrInvalidMethodPost)
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
		log.Error(auth.ErrInvalidForm)
		pkg.SendError(w, http.StatusBadRequest, auth.ErrInvalidForm)
		return
	}
	user, err := h.authUC.SignUp(form)
	if err != nil {
		// ?
		log.Error(fmt.Errorf("faliled to sign up %w", err))
		pkg.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to sign up %v", err.Error()))
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
		log.Error(fmt.Errorf("faliled to create profile %w", err))
		pkg.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to create profile %w", err))
		return
	}

	session, err := h.authUC.CreateSession(user.Id)
	if err != nil {
		log.Error(fmt.Errorf("faliled to sign up %w", err))
		pkg.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to sign up %w", err))
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
		log.Error(auth.ErrInvalidMethodPost)
		pkg.SendError(w, http.StatusMethodNotAllowed, auth.ErrInvalidMethodPost)
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
		pkg.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled decode sign in form %w", err))
		return
	}

	user, err := h.authUC.SignIn(form)
	if err != nil {
		log.Error(fmt.Errorf("faliled to sign in %w", err))
		pkg.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to sign in %w", err))
		return
	}

	// когда логинимся, то обновляем куку, если ранее была, то удалится и пересоздастся
	err = h.authUC.DeleteSessionByUId(user.Id)
	session, err := h.authUC.CreateSession(user.Id)
	if err != nil {
		log.Error(fmt.Errorf("faliled to sign in %w", err))
		pkg.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to sign in %w", err))
		return
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
		log.Error(fmt.Errorf("faliled logout %w", err))
		pkg.SendError(w, http.StatusUnauthorized, fmt.Errorf("faliled to logout %v", err.Error()))
		return
		//return
	} else if err != nil {
		log.Error(fmt.Errorf("faliled logout %w", err))
		pkg.SendError(w, http.StatusBadRequest, fmt.Errorf("faliled to logout %v", err.Error()))
		return
	}

	err = h.authUC.DeleteSession(cookie.Value)
	if err != nil {
		log.Error(fmt.Errorf("faliled logout %w", err))
		pkg.SendError(w, http.StatusUnauthorized, fmt.Errorf("faliled to logout %v", err.Error()))
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
		pkg.SendError(w, http.StatusUnauthorized, fmt.Errorf("faliled auth %w", err))
		return
	}

	_, err = h.authUC.GetSession(cookie.Value)
	if err != nil {
		pkg.SendError(w, http.StatusUnauthorized, fmt.Errorf("failed to auth: %w", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
