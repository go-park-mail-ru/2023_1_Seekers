package middleware

import (
	"bytes"
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	mockSessionUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/crypto"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/logger"
	promMetrics "github.com/go-park-mail-ru/2023_1_Seekers/pkg/metrics/prometheus"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/rand"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createConfig() *config.Config {
	cfg := new(config.Config)
	cfg.Sessions.CookieName = "MailBxSession"
	cfg.Sessions.CookieLen = 32
	cfg.Sessions.CSRFHeader = "Csrf-Token"
	logsUseStdOut := true
	cfg.Logger.LogsUseStdOut = &logsUseStdOut
	cfg.Logger.LogsApiFileName = "test_"
	cfg.Logger.LogsTimeFormat = "2006-01-02_15:04:05_MST"
	cfg.Project.ProjectBaseDir = "2023_1_Seekers"
	cfg.Logger.LogsDir = "logs/"

	return cfg
}

func TestMiddleware_CheckAuth(t *testing.T) {
	cfg := createConfig()
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

	randCookie, err := rand.String(cfg.Sessions.CookieLen)
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

	sUC := mockSessionUC.NewMockUseCaseI(ctrl)
	metrics, _ := promMetrics.NewMetricsHttpServer(cfg.Api.MetricsName)
	globalLogger := logger.Init(logrus.InfoLevel, *cfg.Logger.LogsUseStdOut, cfg.Logger.LogsApiFileName, cfg.Logger.LogsTimeFormat, cfg.Project.ProjectBaseDir, cfg.Logger.LogsDir)
	middleware := NewHttpMiddleware(cfg, sUC, globalLogger, metrics)

	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, test := range testCases {
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte{}))
		r = r.WithContext(context.WithValue(r.Context(), common.ContextHandlerLog, globalLogger))

		if test.createSession && !test.noCookie {
			r.AddCookie(&http.Cookie{
				Name:  cfg.Sessions.CookieName,
				Value: test.cookie,
			})
			sUC.EXPECT().GetSession(test.cookie).Return(&models.Session{}, nil)
		} else if !test.noCookie {
			r.AddCookie(&http.Cookie{
				Name:  cfg.Sessions.CookieName,
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

func TestMiddleware_CheckCSRF(t *testing.T) {
	cfg := createConfig()
	t.Parallel()
	type inputCase struct {
		token    string
		noCookie bool
		cookie   string
	}
	type outputCase struct {
		status int
	}
	type testCase struct {
		inputCase
		outputCase
		name string
	}

	tokenCSRF, _ := crypto.CreateCSRF("randgeneratedcookie12334524524523542")
	testCases := []testCase{
		{
			inputCase: inputCase{
				token:    tokenCSRF,
				noCookie: false,
				cookie:   "randgeneratedcookie12334524524523542",
			},
			outputCase: outputCase{status: http.StatusOK},
			name:       "success",
		},
		{
			inputCase: inputCase{
				token:    "randgeneratedtoken.123345212",
				noCookie: false,
				cookie:   "randgeneratedcookie12334524524523542",
			},
			outputCase: outputCase{status: http.StatusBadRequest},
			name:       "bad token time",
		},
		{
			inputCase: inputCase{
				token:    "",
				noCookie: false,
				cookie:   "randgeneratedcookie12334524524523542",
			},
			outputCase: outputCase{status: http.StatusBadRequest},
			name:       "token not presented",
		},
		{
			inputCase: inputCase{
				token:    "",
				noCookie: true,
				cookie:   "",
			},
			outputCase: outputCase{status: http.StatusUnauthorized},
			name:       "cookie not presented",
		},
	}
	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sUC := mockSessionUC.NewMockUseCaseI(ctrl)
	metrics, _ := promMetrics.NewMetricsHttpServer(cfg.Api.MetricsName)
	globalLogger := logger.Init(logrus.InfoLevel, *cfg.Logger.LogsUseStdOut, cfg.Logger.LogsApiFileName, cfg.Logger.LogsTimeFormat, cfg.Project.ProjectBaseDir, cfg.Logger.LogsDir)
	middleware := NewHttpMiddleware(cfg, sUC, globalLogger, metrics)

	for _, test := range testCases {
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte{}))
		r = r.WithContext(context.WithValue(r.Context(), common.ContextHandlerLog, globalLogger))

		r.Header.Set(cfg.Sessions.CSRFHeader, test.token)

		if !test.noCookie {
			r.AddCookie(&http.Cookie{
				Name:  cfg.Sessions.CookieName,
				Value: test.cookie,
			})
		}
		w := httptest.NewRecorder()

		handler := middleware.CheckCSRF(wrappedHandler)
		handler(w, r)

		if w.Code != test.outputCase.status {
			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.status, w.Code)
		}
	}
}
