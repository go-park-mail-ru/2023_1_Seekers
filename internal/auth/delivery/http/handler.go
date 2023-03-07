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

func New(aUC auth.UseCaseI, uUC _user.UseCaseI, mUC mail.UseCaseI) auth.HandlersI {
	return &handlers{
		authUC: aUC,
		userUC: uUC,
		mailUC: mUC,
	}
}

func handleAuthErr(w http.ResponseWriter, err error) {
	pkg.HandleError(w, auth.Errors[err], err)
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
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		handleAuthErr(w, auth.ErrInvalidForm)
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		handleAuthErr(w, auth.ErrInvalidForm)
		return
	}

	user, err := h.authUC.SignUp(form)
	if err != nil {
		if err == auth.ErrInvalidLogin || err == _user.ErrTooShortPw ||
			err == auth.ErrPwDontMatch || err == auth.ErrFailedCreateProfile {
			handleAuthErr(w, err)
			return
		}
		handleAuthErr(w, auth.ErrUserExists)
		return
	}

	session, err := h.authUC.CreateSession(user.ID)
	if err != nil {
		handleAuthErr(w, auth.ErrFailedCreateSession)
		return
	}

	err = h.mailUC.CreateHelloMessage(user.ID)
	if err != nil {
		handleAuthErr(w, auth.ErrInternalHelloMsg)
		return
	}

	setNewCookie(w, session)
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

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		handleAuthErr(w, auth.ErrInvalidForm)
		return
	}

	user, err := h.authUC.SignIn(form)
	if err != nil {
		if err == auth.ErrInvalidLogin {
			handleAuthErr(w, err)
			return
		}
		handleAuthErr(w, auth.ErrWrongPw)
		return
	}

	// когда логинимся, то обновляем куку, если ранее была, то удалится и пересоздастся
	err = h.authUC.DeleteSessionByUID(user.ID)
	session, err := h.authUC.CreateSession(user.ID)
	if err != nil {
		handleAuthErr(w, auth.ErrFailedCreateSession)
		return
	}

	setNewCookie(w, session)
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

	delCookie(w)

	w.WriteHeader(http.StatusOK)
}
