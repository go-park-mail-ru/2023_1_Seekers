package middleware

//
//import (
//	"bytes"
//	"context"
//	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
//	_sessionRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/redis"
//	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
//	_fStorageRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage/repository"
//	_fStorageUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage/usecase"
//
//	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/inmemory"
//	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
//	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
//	"github.com/redis/go-redis/v9"
//	log "github.com/sirupsen/logrus"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestHandlers_Auth(t *testing.T) {
//	t.Parallel()
//	type inputCase struct {
//		createSession bool //нужно ли создавать сессию
//		noCookie      bool
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
//	randCookie, err := pkg.String(config.CookieLen)
//	if err != nil {
//		t.Errorf("failed generate rand str %v ", err)
//	}
//
//	testCases := []testCase{
//		{
//			// регистрируем пользователя и отправляем с ним куку
//			inputCase{true, false},
//			outputCase{status: http.StatusOK},
//			"success, created cookie",
//		},
//		{
//			// просто приходит кука которая ранее не была создана на сервере
//			inputCase{false, false},
//			outputCase{status: http.StatusUnauthorized},
//			"not valid cookie",
//		},
//		{
//			// если вообще нет куки с таким названием
//			inputCase{false, true},
//			outputCase{status: http.StatusUnauthorized},
//			"cookie not presented",
//		},
//	}
//
//	pkg.InitLogger()
//	logger := pkg.GetLogger()
//
//	if err != nil {
//		logger.Fatalf("db connection error %v", err)
//	}
//
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     config.RedisAddr,
//		Password: config.RedisPassword,
//	})
//
//	_, err = rdb.Ping(context.Background()).Result()
//	if err != nil {
//		log.Fatalf("failed connect to redis : %v", err)
//	}
//	userRepo := _userRepo.New()
//	sessionRepo := _sessionRepo.NewSessionRepo(rdb)
//	fStorageRepo := _fStorageRepo.New()
//
//	fStorageUC := _fStorageUCase.New(fStorageRepo)
//	usersUC := _userUCase.New(userRepo, fStorageUC)
//	sessionUC := _authUCase.NewSessionUC(sessionRepo, usersUC)
//
//	middleware := New(sessionUC, logger)
//
//	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//	})
//
//	for _, test := range testCases {
//		r := httptest.NewRequest("POST", "/", bytes.NewReader([]byte{}))
//
//		if test.createSession && !test.noCookie {
//			r.AddCookie(&http.Cookie{
//				Name:  config.CookieName,
//				Value: "randgeneratedcookie12334524524523542",
//			}) //выставляем авторизованную куку (ранее созданную дл тестирования)
//		} else if !test.noCookie {
//			r.AddCookie(&http.Cookie{
//				Name:  config.CookieName,
//				Value: randCookie,
//			}) //создаем невалидную куку
//		}
//		w := httptest.NewRecorder()
//
//		handler := middleware.CheckAuth(wrappedHandler)
//		handler(w, r)
//
//		if w.Code != test.outputCase.status {
//			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.status, w.Code)
//		}
//	}
//}
