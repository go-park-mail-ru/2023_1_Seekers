package middleware

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/crypto"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgHttp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/http"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/logger"
	pkgErrors "github.com/pkg/errors"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Middleware struct {
	sUC auth.SessionUseCaseI
	log *logger.Logger
}

func New(sUC auth.SessionUseCaseI, l *logger.Logger) *Middleware {
	return &Middleware{sUC, l}
}

func (m *Middleware) Cors(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedMethods:   config.AllowedMethods,
		AllowedOrigins:   config.AllowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   config.AllowedHeaders,
	})
	return c.Handler(h)
}

func (m *Middleware) HandlerLogger(h http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerLogger := m.log.WithFields(map[string]any{
			"method": r.Method,
			"url":    r.URL.Path,
		})
		handlerLogger.Info("new request")
		r = r.WithContext(context.WithValue(r.Context(), common.ContextHandlerLog, handlerLogger))
		fmt.Println(r.Context())
		fmt.Println(handlerLogger)
		_, ok := r.Context().Value(common.ContextHandlerLog).(*logger.Logger)
		if !ok {
			log.Error("failed to get logger for handler")
		}
		h.ServeHTTP(w, r)
	})
	return handler
}

func (m *Middleware) CheckAuth(h http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(config.CookieName)
		if err != nil {
			pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedAuth, err.Error()))
			return
		}
		session, err := m.sUC.GetSession(cookie.Value)
		if err != nil {
			pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedAuth, err.Error()))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, session.UID))

		h.ServeHTTP(w, r)
	})
	return handler
}

func (m *Middleware) CheckCSRF(h http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(config.CookieName)
		if err != nil {
			pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedAuth, err.Error()))
			return
		}
		csrfToken := r.Header.Get(config.CSRFHeader)
		if csrfToken == "" {
			pkgHttp.HandleError(w, r, pkgErrors.WithMessage(errors.ErrWrongCSRF, "token not presented"))
			return
		}

		err = crypto.CheckCSRF(cookie.Value, csrfToken)
		if err != nil {
			pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed check csrf"))
			return
		}
		h.ServeHTTP(w, r)
	})
	return handler
}
