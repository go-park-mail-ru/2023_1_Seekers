package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Middleware struct {
	uc auth.UseCaseI
}

func New(aUc auth.UseCaseI) *Middleware {
	return &Middleware{aUc}
}

func (m *Middleware) Cors(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedMethods:   []string{"POST", "GET", "PUT"},
		AllowedOrigins:   []string{"http://127.0.0.1:8002", "http://http://89.208.197.150:8002"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Content-Length", "X-Csrf-Token"},
		Debug:            true,
	})
	return c.Handler(h)
}

func (m *Middleware) CheckAuth(h http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(config.CookieName)
		if err != nil {
			authErr := errors.NewWrappedErr(auth.Errors[auth.ErrFailedAuth], auth.ErrFailedAuth.Error(), err)
			log.Error(authErr)
			pkg.SendError(w, authErr)
			return
		}
		session, err := m.uc.GetSession(cookie.Value)
		if err != nil {
			authErr := errors.NewWrappedErr(auth.Errors[auth.ErrFailedGetSession], auth.ErrFailedGetSession.Error(), err)
			log.Error(authErr)
			pkg.SendError(w, authErr)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), config.ContextUser, session.UID))

		h.ServeHTTP(w, r)
	})
	return handler
}
