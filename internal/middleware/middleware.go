package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Middleware struct {
	uc auth.UseCase
}

func New(aUc auth.UseCase) *Middleware {
	return &Middleware{aUc}
}

func (m *Middleware) CheckAuth(h http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(config.CookieName)
		if err != nil {
			authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedAuth], auth.ErrFailedAuth.Error(), err)
			log.Error(authErr)
			pkg.SendError(w, authErr)
			return
		}
		session, err := m.uc.GetSession(cookie.Value)
		if err != nil {
			authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedGetSession], auth.ErrFailedGetSession.Error(), err)
			log.Error(authErr)
			pkg.SendError(w, authErr)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), config.ContextUser, session.UID))

		h.ServeHTTP(w, r)
	})
	return handler
}
