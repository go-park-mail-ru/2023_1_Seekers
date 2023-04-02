package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"github.com/rs/cors"
	"net/http"
)

type Middleware struct {
	sUC auth.SessionUseCaseI
	log *pkg.Logger
}

func New(sUC auth.SessionUseCaseI, l *pkg.Logger) *Middleware {
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
		handlerLogger := m.log.LoggerWithFields(map[string]any{
			"method": r.Method,
			"url":    r.URL.Path,
		})
		handlerLogger.Info("new request")
		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextHandlerLog, handlerLogger))
		h.ServeHTTP(w, r)
	})
	return handler
}

func (m *Middleware) CheckAuth(h http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(config.CookieName)
		if err != nil {
			pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedAuth, err.Error()))
			return
		}
		session, err := m.sUC.GetSession(cookie.Value)
		if err != nil {
			pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedGetSession, err.Error()))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, session.UID))

		h.ServeHTTP(w, r)
	})
	return handler
}

func (m *Middleware) CheckCSRF(h http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(config.CookieName)
		if err != nil {
			pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedAuth, err.Error()))
			return
		}
		csrfToken := r.Header.Get(config.CSRFHeader)
		if csrfToken == "" {
			pkg.HandleError(w, r, pkgErrors.WithMessage(errors.ErrWrongCSRF, "token not presented"))
			return
		}

		err = pkg.CheckCSRF(cookie.Value, csrfToken)
		if err != nil {
			pkg.HandleError(w, r, pkgErrors.Wrap(err, "failed check csrf"))
			return
		}
		h.ServeHTTP(w, r)
	})
	return handler
}
