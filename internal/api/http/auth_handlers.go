package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/crypto"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	_ "github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgHttp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/http"
	"github.com/go-playground/validator/v10"
	pkgErrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type AuthHandlersI interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
	GetCSRF(w http.ResponseWriter, r *http.Request)
}

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCaseI
	mailUC mail.UseCaseI
	userUC _user.UseCaseI
}

func NewAuthHandlers(c *config.Config, aUC auth.UseCaseI, mUC mail.UseCaseI, uUC _user.UseCaseI) AuthHandlersI {
	return &authHandlers{
		cfg:    c,
		authUC: aUC,
		mailUC: mUC,
		userUC: uUC,
	}
}

// SignUp godoc
// @Summary      SignUp
// @Description  user sign up
// @Tags     auth
// @Accept	 application/json
// @Produce  application/json
// @Param    user body models.FormSignUp true "user info"
// @Success  200 {object} models.AuthResponse "user created"
// @Failure 401 {object} errors.JSONError "passwords don`t match"
// @Failure 401 {object} errors.JSONError "invalid login"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 403 {object} errors.JSONError "password too short"
// @Failure 409 {object} errors.JSONError "user already exists"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /signup [post]
func (h *authHandlers) SignUp(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("failed to close request: ", err)
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed read request body"))
		return
	}

	form := models.FormSignUp{}
	if err := form.UnmarshalJSON(body); err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	form.Sanitize()

	response, session, err := h.authUC.SignUp(&form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	user, err := h.userUC.GetInfoByEmail(response.Email)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	_, err = h.mailUC.CreateDefaultFolders(user.UserID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	err = h.mailUC.SendWelcomeMessage(response.Email)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	h.setNewCookie(w, session)
	pkgHttp.SendJSON(w, r, http.StatusOK, response)
}

// SignIn godoc
// @Summary      SignIn
// @Description  user sign in
// @Tags     auth
// @Accept	 application/json
// @Produce  application/json
// @Param    user body models.FormLogin true "user info"
// @Success  200 {object} models.AuthResponse "success sign in"
// @Failure 401 {object} errors.JSONError "invalid login"
// @Failure 401 {object} errors.JSONError "wrong password"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /signin [post]
func (h *authHandlers) SignIn(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(fmt.Errorf("failed to close request: %w", err))
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed read request body"))
		return
	}

	form := models.FormLogin{}
	if err := form.UnmarshalJSON(body); err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	response, session, err := h.authUC.SignIn(&form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	h.setNewCookie(w, session)
	pkgHttp.SendJSON(w, r, http.StatusOK, response)
}

// Auth godoc
// @Summary      Auth
// @Description  check is user authorised
// @Tags     auth
// @Accept	 application/json
// @Produce  application/json
// @Success  200 "success auth"
// @Failure 401 {object} errors.JSONError "failed auth"
// @Router   /auth [get]
func (h *authHandlers) Auth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Logout godoc
// @Summary      Logout
// @Description  check is user authorised
// @Tags     auth
// @Accept	 application/json
// @Produce  application/json
// @Success  200 "success logout"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /logout [delete]
func (h *authHandlers) Logout(w http.ResponseWriter, _ *http.Request) {
	h.delCookie(w)
	w.WriteHeader(http.StatusOK)
}

// GetCSRF godoc
// @Summary      GetCSRF
// @Description  Get CSRF token
// @Tags         auth
// @Success      200    "success create csrf"
// @Failure 401 {object} errors.JSONError "failed get user"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router /create_csrf [get]
func (h *authHandlers) GetCSRF(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cfg.Sessions.CookieName)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedAuth, err.Error()))
		return
	}

	csrfToken, err := crypto.CreateCSRF(cookie.Value)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}
	w.Header().Set(h.cfg.Sessions.CSRFHeader, csrfToken)
	w.WriteHeader(http.StatusOK)
}

func (h *authHandlers) setNewCookie(w http.ResponseWriter, session *models.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.Sessions.CookieName,
		Value:    session.SessionID,
		Expires:  time.Now().Add(h.cfg.Sessions.CookieTTL),
		HttpOnly: true,
		Path:     h.cfg.Sessions.CookiePath,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *authHandlers) delCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    h.cfg.Sessions.CookieName,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    h.cfg.Sessions.CookiePath,
	})
}
