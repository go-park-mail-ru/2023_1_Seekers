package middleware

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/build/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Middleware struct {
	auth.UseCase
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
		_, err = m.GetSession(cookie.Value)
		if err != nil {
			authErr := errors.NewWrappedErr(auth.AuthErrors[auth.ErrFailedGetSession], auth.ErrFailedGetSession.Error(), err)
			log.Error(authErr)
			pkg.SendError(w, authErr)
			return
		}
		h.ServeHTTP(w, r)
	})
	return handler
}
