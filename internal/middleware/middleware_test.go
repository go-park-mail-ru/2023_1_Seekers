package middleware

import (
	"bytes"
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	mockSessionUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase/mocks_session"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/rand"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_CheckAuth(t *testing.T) {
	t.Parallel()
	type inputCase struct {
		createSession bool //нужно ли создавать сессию
		noCookie      bool
		cookie        string
	}
	type outputCase struct {
		status int
	}
	type testCase struct {
		inputCase
		outputCase
		name string
	}

	randCookie, err := rand.String(config.CookieLen)
	if err != nil {
		t.Errorf("failed generate rand str %v ", err)
	}

	testCases := []testCase{
		{
			// регистрируем пользователя и отправляем с ним куку
			inputCase:  inputCase{createSession: true, cookie: "randgeneratedcookie12334524524523542"},
			outputCase: outputCase{status: http.StatusOK},
			name:       "success, created cookie",
		},
		{
			// просто приходит кука которая ранее не была создана на сервере
			inputCase:  inputCase{cookie: randCookie},
			outputCase: outputCase{status: http.StatusUnauthorized},
			name:       "not valid cookie",
		},
		{
			// если вообще нет куки с таким названием
			inputCase:  inputCase{noCookie: true},
			outputCase: outputCase{status: http.StatusUnauthorized},
			name:       "cookie not presented",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sUC := mockSessionUC.NewMockSessionUseCaseI(ctrl)
	pkg.InitLogger()
	logger := pkg.GetLogger()
	middleware := New(sUC, logger)

	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, test := range testCases {
		r := httptest.NewRequest("POST", "/", bytes.NewReader([]byte{}))
		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextHandlerLog, logger))

		if test.createSession && !test.noCookie {
			r.AddCookie(&http.Cookie{
				Name:  config.CookieName,
				Value: test.cookie,
			})
			sUC.EXPECT().GetSession(test.cookie).Return(&models.Session{}, nil)
		} else if !test.noCookie {
			r.AddCookie(&http.Cookie{
				Name:  config.CookieName,
				Value: test.cookie,
			})
			sUC.EXPECT().GetSession(test.cookie).Return(nil, errors.ErrFailedAuth)
		}
		w := httptest.NewRecorder()

		handler := middleware.CheckAuth(wrappedHandler)
		handler(w, r)

		if w.Code != test.outputCase.status {
			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.status, w.Code)
		}
	}
}

//func TestMiddleware_CheckCSRF(t *testing.T) {
//	t.Parallel()
//	type inputCase struct {
//		token    string
//		noCookie bool
//		cookie   string
//	}
//	type outputCase struct {
//		status int
//	}
//	type testCase struct {
//		inputCase
//		outputCase
//		name string
//	}
//
//	tokenCSRF, _ := pkg.CreateCSRF("randgeneratedcookie12334524524523542")
//	testCases := []testCase{
//		{
//			inputCase: inputCase{
//				token:    tokenCSRF,
//				noCookie: false,
//				cookie:   "randgeneratedcookie12334524524523542",
//			},
//			outputCase: outputCase{status: http.StatusOK},
//			name:       "success",
//		},
//		{
//			inputCase: inputCase{
//				token:    "randgeneratedtoken.123345212",
//				noCookie: false,
//				cookie:   "randgeneratedcookie12334524524523542",
//			},
//			outputCase: outputCase{status: http.StatusBadRequest},
//			name:       "bad token time",
//		},
//		{
//			inputCase: inputCase{
//				token:    "",
//				noCookie: false,
//				cookie:   "randgeneratedcookie12334524524523542",
//			},
//			outputCase: outputCase{status: http.StatusBadRequest},
//			name:       "token not presented",
//		},
//		{
//			inputCase: inputCase{
//				token:    "",
//				noCookie: true,
//				cookie:   "",
//			},
//			outputCase: outputCase{status: http.StatusUnauthorized},
//			name:       "cookie not presented",
//		},
//	}
//	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//	})
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	sUC := mockSessionUC.NewMockSessionUseCaseI(ctrl)
//	pkg.InitLogger()
//	logger := pkg.GetLogger()
//	middleware := New(sUC, logger)
//
//	for _, test := range testCases {
//		r := httptest.NewRequest("POST", "/", bytes.NewReader([]byte{}))
//		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextHandlerLog, logger))
//
//		r.Header.Set(config.CSRFHeader, test.token)
//
//		if !test.noCookie {
//			r.AddCookie(&http.Cookie{
//				Name:  config.CookieName,
//				Value: test.cookie,
//			})
//		}
//		w := httptest.NewRecorder()
//
//		handler := middleware.CheckCSRF(wrappedHandler)
//		handler(w, r)
//
//		if w.Code != test.outputCase.status {
//			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.status, w.Code)
//		}
//	}
//}
