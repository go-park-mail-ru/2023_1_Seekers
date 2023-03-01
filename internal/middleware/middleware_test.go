package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_Seekers/build/config"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_Auth(t *testing.T) {
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
			// просто приходит кука которая ранее не была создана на сервере
			inputCase{[]byte(
				`{"email":"` + randStr + `testing_auth2@example.com",
				"password":  "4321",
				"repeat_pw":  "4321",
				"first_name": "Ivan",
				"last_name":  "Ivanov",
				"birth_date": "29.01.1999"}`), false, false},
			outputCase{status: http.StatusUnauthorized},
			"not valid cookie",
		},
		{
			// если вообще нет куки с таким названием
			inputCase{[]byte(
				`{"email":"` + randStr + `testing_auth3@example.com",
				"password":  "4321",
				"repeat_pw":  "4321",
				"first_name": "Ivan",
				"last_name":  "Ivanov",
				"birth_date": "29.01.1999"}`), false, true},
			outputCase{status: http.StatusUnauthorized},
			"cookie not presented",
		},
	}

	authRepo := _authRepo.New()
	authUCase := _authUCase.New(authRepo)
	middleware := New(authUCase)

	for _, test := range testCases {
		r := httptest.NewRequest("POST", "/api/auth", bytes.NewReader([]byte{}))
		var user model.User
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

		authH.Auth(w, r)

		if w.Code != test.outputCase.status {
			t.Errorf("[TEST] %s : Expected status %d, got %d ", test.name, test.status, w.Code)
		}
	}
}
