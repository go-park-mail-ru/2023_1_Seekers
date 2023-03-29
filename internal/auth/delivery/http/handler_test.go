package http

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
//	_sessionRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/redis"
//	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
//	_fStorageRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage/repository"
//	_fStorageUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage/usecase"
//	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/repository/postgres"
//	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase"
//	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
//	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
//	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/inmemory"
//	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
//	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
//	"github.com/redis/go-redis/v9"
//	log "github.com/sirupsen/logrus"
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//	"gorm.io/gorm/schema"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestHandlers_SignIn(t *testing.T) {
//	t.Parallel()
//	type outputCase struct {
//		status int
//	}
//	type testCase struct {
//		input []byte
//		outputCase
//		name string
//	}
//
//	testCases := []testCase{
//		{
//			[]byte(`{"login": "test", "password": "12345"}`),
//			outputCase{status: http.StatusOK},
//			"default success",
//		},
//		{
//			[]byte(`{"login": "test_signin1", "password": "43212"}`),
//			outputCase{status: http.StatusUnauthorized},
//			"no such user",
//		},
//		{
//			[]byte(`{"login: test_signin", "password": "12334"}`),
//			outputCase{status: http.StatusForbidden},
//			"invalid form",
//		},
//	}
//
//	pkg.InitLogger()
//	logger := pkg.GetLogger()
//
//	var connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
//		config.DBUser,
//		config.DBPassword,
//		config.DBHost,
//		config.DBPort,
//		config.DBName,
//	)
//
//	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connStr}), &gorm.Config{NamingStrategy: schema.NamingStrategy{
//		TablePrefix:   config.DBSchemaName + ".",
//		SingularTable: false,
//	}})
//	if err != nil {
//		logger.Fatalf("db connection error %v", err)
//	}
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
//	mailRepo := _mailRepo.New(db)
//	fStorageRepo := _fStorageRepo.New()
//
//	fStorageUC := _fStorageUCase.New(fStorageRepo)
//	usersUC := _userUCase.New(userRepo, fStorageUC)
//	sessionUC := _authUCase.NewSessionUC(sessionRepo, usersUC)
//	mailUC := _mailUCase.New(mailRepo, usersUC)
//	authUC := _authUCase.NewAuthUC(sessionUC, usersUC, mailUC)
//
//	authH := New(authUC)
//	for _, test := range testCases {
//		r := httptest.NewRequest("POST", "/api/signin", bytes.NewReader(test.input))
//		w := httptest.NewRecorder()
//
//		authH.SignIn(w, r)
//
//		if w.Code != test.status {
//			t.Errorf("[TEST] %s : Expected status %d, got %d ", test.name, test.status, w.Code)
//		}
//	}
//}
//
//func TestHandlers_SignUp(t *testing.T) {
//	t.Parallel()
//	type outputCase struct {
//		status int
//	}
//	type testCase struct {
//		input []byte
//		outputCase
//		name string
//	}
//	randStr, err := pkg.String(3)
//	if err != nil {
//		t.Errorf("failed generate rand str %v ", err)
//	}
//	testCases := []testCase{
//		{
//			[]byte(`{"login":"` + randStr + `testing_signup1",
//						  "password":"54321",
//						  "repeat_pw":"54321",
//						  "first_name":"Ivan",
//						  "last_name":"Ivanov",
//						  "birth_date":"29.01.2002"}`),
//			outputCase{status: http.StatusOK},
//			"default success",
//		},
//		{
//			[]byte(`{"login":"` + randStr + `testing_signup2",
//                           "password":"54321",
//                           "repeat_pw":"12311",
//                           "first_name":"Ivan",
//                           "last_name":"Ivanov",
//                           "birth_date":"29.01.2002"}`),
//			outputCase{status: http.StatusUnauthorized},
//			"passwords dont match",
//		},
//		{
//			[]byte(`{"login: + testing_signup2",
//                           "password:"54321",
//                           "repeat_pw":"12313",
//                           "first_name":"Ivan",
//                           "last_name":"Ivanov",
//                           "birth_date":"29.01.2002"}`),
//			outputCase{status: http.StatusForbidden},
//			"invalid form",
//		},
//		{
//			[]byte(`{"login":"test",
//                           "password":"54321",
//                           "repeat_pw":"54321",
//                           "first_name":"Ivan",
//                           "last_name":"Ivanov",
//                           "birth_date":"29.01.2002"}`),
//			outputCase{status: http.StatusConflict},
//			"user with such login exists",
//		},
//	}
//	pkg.InitLogger()
//	logger := pkg.GetLogger()
//
//	var connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
//		config.DBUser,
//		config.DBPassword,
//		config.DBHost,
//		config.DBPort,
//		config.DBName,
//	)
//
//	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connStr}), &gorm.Config{NamingStrategy: schema.NamingStrategy{
//		TablePrefix:   config.DBSchemaName + ".",
//		SingularTable: false,
//	}})
//	if err != nil {
//		logger.Fatalf("db connection error %v", err)
//	}
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
//	mailRepo := _mailRepo.New(db)
//	fStorageRepo := _fStorageRepo.New()
//
//	fStorageUC := _fStorageUCase.New(fStorageRepo)
//	usersUC := _userUCase.New(userRepo, fStorageUC)
//	sessionUC := _authUCase.NewSessionUC(sessionRepo, usersUC)
//	mailUC := _mailUCase.New(mailRepo, usersUC)
//	authUC := _authUCase.NewAuthUC(sessionUC, usersUC, mailUC)
//
//	authH := New(authUC)
//
//	for _, test := range testCases {
//		r := httptest.NewRequest("POST", "/api/signup", bytes.NewReader(test.input))
//		w := httptest.NewRecorder()
//
//		authH.SignUp(w, r)
//
//		if w.Code != test.status {
//			t.Errorf("[TEST] %s , Expected status %d, got %d ", test.name, test.status, w.Code)
//		}
//	}
//}
//
//func TestHandlers_Logout(t *testing.T) {
//	t.Parallel()
//	type inputCase struct {
//		user          []byte
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
//	randStr, err := pkg.String(3)
//	if err != nil {
//		t.Errorf("failed generate rand str %v ", err)
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
//			inputCase{[]byte(
//				`{"login":"` + randStr + `testing_auth1",
//				"password":  "54321",
//				"repeat_pw":  "54321",
//				"first_name": "Ivan",
//				"last_name":  "Ivanov",
//				"birth_date": "29.01.1999"}`), true, false},
//			outputCase{status: http.StatusOK},
//			"success, created cookie",
//		},
//		{
//			// просто приходит кука которая ранее не была создана на сервере
//			inputCase{[]byte(
//				`{"login":"` + randStr + `testing_auth2",
//				"password":  "54321",
//				"repeat_pw":  "54321",
//				"first_name": "Ivan",
//				"last_name":  "Ivanov",
//				"birth_date": "29.01.1999"}`), false, false},
//			outputCase{status: http.StatusUnauthorized},
//			"not valid cookie",
//		},
//		{
//			// если вообще нет куки с таким названием
//			inputCase{[]byte(
//				`{"login":"` + randStr + `testing_auth3",
//				"password":  "54321",
//				"repeat_pw":  "54321",
//				"first_name": "Ivan",
//				"last_name":  "Ivanov",
//				"birth_date": "29.01.1999"}`), false, true},
//			outputCase{status: http.StatusUnauthorized},
//			"cookie not presented",
//		},
//	}
//
//	var connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
//		config.DBUser,
//		config.DBPassword,
//		config.DBHost,
//		config.DBPort,
//		config.DBName,
//	)
//
//	logger := pkg.GetLogger()
//
//	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connStr}), &gorm.Config{NamingStrategy: schema.NamingStrategy{
//		TablePrefix:   config.DBSchemaName + ".",
//		SingularTable: false,
//	}})
//	if err != nil {
//		logger.Fatalf("db connection error %v", err)
//	}
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
//
//	userRepo := _userRepo.New()
//	sessionRepo := _sessionRepo.NewSessionRepo(rdb)
//	mailRepo := _mailRepo.New(db)
//	fStorageRepo := _fStorageRepo.New()
//
//	fStorageUC := _fStorageUCase.New(fStorageRepo)
//	usersUC := _userUCase.New(userRepo, fStorageUC)
//	sessionUC := _authUCase.NewSessionUC(sessionRepo, usersUC)
//	mailUC := _mailUCase.New(mailRepo, usersUC)
//	authUC := _authUCase.NewAuthUC(sessionUC, usersUC, mailUC)
//
//	authH := New(authUC)
//	middleware := _middleware.New(sessionUC, logger)
//
//	for _, test := range testCases {
//		r := httptest.NewRequest("POST", "/api/logout", bytes.NewReader([]byte{}))
//		var user models.AuthResponse
//		if test.createSession && !test.noCookie {
//			signupReq := httptest.NewRequest("POST", "/api/signup", bytes.NewReader(test.user))
//			w := httptest.NewRecorder()
//
//			authH.SignUp(w, signupReq)
//			json.NewDecoder(w.Body).Decode(&user)
//			// TODO
//			//s, err := sessionUC.GetSession(user.Email)
//			//if err != nil {
//			//	t.Errorf("failed to get session %v ", err)
//			//}
//			//
//			//r.AddCookie(&http.Cookie{
//			//	Name:  config.CookieName,
//			//	Value: s.SessionID,
//			//}) //необходимо проверить если нет кук,поэтому в случае пустого кейса - кука не выставится
//		} else if !test.noCookie {
//			r.AddCookie(&http.Cookie{
//				Name:  config.CookieName,
//				Value: randCookie,
//			}) //создаем невалидную куку
//		}
//		w := httptest.NewRecorder()
//
//		middleware.CheckAuth(authH.Logout)(w, r)
//
//		if w.Code != test.outputCase.status {
//			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.status, w.Code)
//		}
//	}
//}
