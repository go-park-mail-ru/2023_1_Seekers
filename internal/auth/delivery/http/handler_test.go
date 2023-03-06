package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/repository/inmemory"
	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_SignIn(t *testing.T) {
	t.Parallel()
	type outputCase struct {
		status int
	}
	type testCase struct {
		input []byte
		outputCase
		name string
	}

	testCases := []testCase{
		{
			[]byte(`{"login": "test", "password": "12345"}`),
			outputCase{status: http.StatusOK},
			"default success",
		},
		{
			[]byte(`{"login": "test_signin1", "password": "43212"}`),
			outputCase{status: http.StatusUnauthorized},
			"no such user",
		},
		{
			[]byte(`{"login: test_signin", "password": "12334"}`),
			outputCase{status: http.StatusForbidden},
			"invalid form",
		},
	}

	userRepo := _userRepo.New()
	authRepo := _authRepo.New()
	mailRepo := _mailRepo.New(userRepo)

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo, usersUCase)
	mailUCase := _mailUCase.New(mailRepo)

	authH := New(authUCase, usersUCase, mailUCase)
	for _, test := range testCases {
		r := httptest.NewRequest("POST", "/api/signin", bytes.NewReader(test.input))
		w := httptest.NewRecorder()

		authH.SignIn(w, r)

		if w.Code != test.status {
			t.Errorf("[TEST] %s : Expected status %d, got %d ", test.name, test.status, w.Code)
		}
	}
}

func TestHandlers_SignUp(t *testing.T) {
	t.Parallel()
	type outputCase struct {
		status int
	}
	type testCase struct {
		input []byte
		outputCase
		name string
	}
	randStr, err := pkg.String(3)
	if err != nil {
		t.Errorf("failed generate rand str %v ", err)
	}
	testCases := []testCase{
		{
			[]byte(`{"login":"` + randStr + `testing_signup1",
						  "password":"54321",
						  "repeat_pw":"54321",
						  "first_name":"Ivan",
						  "last_name":"Ivanov",
						  "birth_date":"29.01.2002"}`),
			outputCase{status: http.StatusOK},
			"default success",
		},
		{
			[]byte(`{"login":"` + randStr + `testing_signup2",
                           "password":"54321",
                           "repeat_pw":"12311",
                           "first_name":"Ivan", 
                           "last_name":"Ivanov", 
                           "birth_date":"29.01.2002"}`),
			outputCase{status: http.StatusUnauthorized},
			"passwords dont match",
		},
		{
			[]byte(`{"login: + testing_signup2",
                           "password:"54321",
                           "repeat_pw":"12313",
                           "first_name":"Ivan", 
                           "last_name":"Ivanov", 
                           "birth_date":"29.01.2002"}`),
			outputCase{status: http.StatusForbidden},
			"invalid form",
		},
		{
			[]byte(`{"login":"test",
                           "password":"54321",
                           "repeat_pw":"54321",
                           "first_name":"Ivan", 
                           "last_name":"Ivanov", 
                           "birth_date":"29.01.2002"}`),
			outputCase{status: http.StatusConflict},
			"user with such login exists",
		},
	}
	userRepo := _userRepo.New()
	authRepo := _authRepo.New()
	mailRepo := _mailRepo.New(userRepo)

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo, usersUCase)
	mailUCase := _mailUCase.New(mailRepo)

	authH := New(authUCase, usersUCase, mailUCase)

	for _, test := range testCases {
		r := httptest.NewRequest("POST", "/api/signup", bytes.NewReader(test.input))
		w := httptest.NewRecorder()

		authH.SignUp(w, r)

		if w.Code != test.status {
			t.Errorf("[TEST] %s , Expected status %d, got %d ", test.name, test.status, w.Code)
		}
	}
}

func TestHandlers_Logout(t *testing.T) {
	t.Parallel()
	type inputCase struct {
		user          []byte
		createSession bool //нужно ли создавать сессию
		noCookie      bool
	}
	type outputCase struct {
		status int
	}
	type testCase struct {
		inputCase
		outputCase
		name string
	}

	randStr, err := pkg.String(3)
	if err != nil {
		t.Errorf("failed generate rand str %v ", err)
	}

	randCookie, err := pkg.String(config.CookieLen)
	if err != nil {
		t.Errorf("failed generate rand str %v ", err)
	}

	testCases := []testCase{
		{
			// регистрируем пользователя и отправляем с ним куку
			inputCase{[]byte(
				`{"login":"` + randStr + `testing_auth1",
				"password":  "54321",
				"repeat_pw":  "54321",
				"first_name": "Ivan",
				"last_name":  "Ivanov",
				"birth_date": "29.01.1999"}`), true, false},
			outputCase{status: http.StatusOK},
			"success, created cookie",
		},
		{
			// просто приходит кука которая ранее не была создана на сервере
			inputCase{[]byte(
				`{"login":"` + randStr + `testing_auth2",
				"password":  "54321",
				"repeat_pw":  "54321",
				"first_name": "Ivan",
				"last_name":  "Ivanov",
				"birth_date": "29.01.1999"}`), false, false},
			outputCase{status: http.StatusUnauthorized},
			"not valid cookie",
		},
		{
			// если вообще нет куки с таким названием
			inputCase{[]byte(
				`{"login":"` + randStr + `testing_auth3",
				"password":  "54321",
				"repeat_pw":  "54321",
				"first_name": "Ivan",
				"last_name":  "Ivanov",
				"birth_date": "29.01.1999"}`), false, true},
			outputCase{status: http.StatusUnauthorized},
			"cookie not presented",
		},
	}

	userRepo := _userRepo.New()
	authRepo := _authRepo.New()
	mailRepo := _mailRepo.New(userRepo)

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo, usersUCase)
	mailUCase := _mailUCase.New(mailRepo)

	authH := New(authUCase, usersUCase, mailUCase)
	middleware := _middleware.New(authUCase)

	for _, test := range testCases {
		r := httptest.NewRequest("POST", "/api/logout", bytes.NewReader([]byte{}))
		var user models.User
		if test.createSession && !test.noCookie {
			signupReq := httptest.NewRequest("POST", "/api/signup", bytes.NewReader(test.user))
			w := httptest.NewRecorder()

			authH.SignUp(w, signupReq)
			json.NewDecoder(w.Body).Decode(&user)
			s, err := authUCase.GetSessionByUID(user.ID)
			if err != nil {
				t.Errorf("failed to get session %v ", err)
			}

			r.AddCookie(&http.Cookie{
				Name:  config.CookieName,
				Value: s.SessionID,
			}) //необходимо проверить если нет кук,поэтому в случае пустого кейса - кука не выставится
		} else if !test.noCookie {
			r.AddCookie(&http.Cookie{
				Name:  config.CookieName,
				Value: randCookie,
			}) //создаем невалидную куку
		}
		w := httptest.NewRecorder()

		middleware.CheckAuth(authH.Logout)(w, r)

		if w.Code != test.outputCase.status {
			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.status, w.Code)
		}
	}
}
