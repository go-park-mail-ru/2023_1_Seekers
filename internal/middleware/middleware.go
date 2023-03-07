package middleware

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/rs/cors"
	"net/http"
)

type Middleware struct {
	uc  auth.UseCaseI
	log *pkg.Logger
}

func New(aUc auth.UseCaseI, l *pkg.Logger) *Middleware {
	return &Middleware{aUc, l}
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
			wrappedErr := fmt.Errorf("%v: %w", auth.ErrFailedAuth, err)
			pkg.HandleError(w, r, auth.Errors[auth.ErrFailedAuth], wrappedErr)
			return
		}
		session, err := m.uc.GetSession(cookie.Value)
		if err != nil {
			wrappedErr := fmt.Errorf("%v: %w", auth.ErrFailedGetSession, err)
			pkg.HandleError(w, r, auth.Errors[auth.ErrFailedGetSession], wrappedErr)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, session.UID))

		h.ServeHTTP(w, r)
	})
	return handler
}
