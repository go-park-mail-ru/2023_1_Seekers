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
	//handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Access-Control-Allow-Origin", "*")
	//	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, Content-Length, User-Agent, X-CSRF-Token")
	//	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE, OPTIONS")
	//	w.Header().Set("Access-Control-Allow-Credentials", "true")
	//	h.ServeHTTP(w, r)
	//})
	//return handler

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET,POST,PUT,DELETE, OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Content-Length", "X-Csrf-Token"},
		AllowedOrigins:   []string{"http://localhost:8002"},
		AllowCredentials: true,
		Debug:            false,
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
