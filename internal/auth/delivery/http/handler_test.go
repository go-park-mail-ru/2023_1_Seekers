package http

import (
	"bytes"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
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
		input  []byte
		output outputCase
		name   string
	}
	testCases := []testCase{
		{
			[]byte(`{"email": "test@example.com", "password": "1234"}`),
			outputCase{status: http.StatusOK},
			"default success",
		},
		{
			[]byte(`{"email": "test@example.com", "password": "4321"}`),
			outputCase{status: http.StatusUnauthorized},
			"default success",
		},
	}

	userRepo := _userRepo.New()
	authRepo := _authRepo.New()

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo)

	authH := New(authUCase, usersUCase)
	for _, test := range testCases {
		r := httptest.NewRequest("POST", "/api/signin", bytes.NewReader(test.input))
		w := httptest.NewRecorder()

		authH.SignIn(w, r)

		if w.Code != test.output.status {
			t.Errorf("Expected status %d, got %d ", test.output.status, w.Code)
		}
	}
}

func TestHandlers_SignUp(t *testing.T) {
	t.Parallel()
	type outputCase struct {
		status int
	}
	type testCase struct {
		input  []byte
		output outputCase
		name   string
	}
	testCases := []testCase{
		{
			[]byte(`{"email":"testing_signup@example.com",
                           "password":"4321",
                           "repeat_pw":"4321",
                           "first_name":"Ivan", 
                           "last_name":"Ivanov", 
                           "birth_date":"29.01.2002"}`),
			outputCase{status: http.StatusOK},
			"default success",
		},
		{
			[]byte(`{"email":"testing_signup@example.com",
                           "password":"4321",
                           "repeat_pw":"1231",
                           "first_name":"Ivan", 
                           "last_name":"Ivanov", 
                           "birth_date":"29.01.2002"}`),
			outputCase{status: http.StatusUnauthorized},
			"passwords dont match",
		},
	}
	userRepo := _userRepo.New()
	authRepo := _authRepo.New()

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo)

	authH := New(authUCase, usersUCase)

	for _, test := range testCases {
		r := httptest.NewRequest("POST", "/api/signup", bytes.NewReader(test.input))
		w := httptest.NewRecorder()

		authH.SignUp(w, r)

		if w.Code != test.output.status {
			t.Errorf("Expected status %d, got %d ", test.output.status, w.Code)
		}
	}
}

func TestHandlers_Auth(t *testing.T) {
	t.Parallel()
	userRepo := _userRepo.New()
	authRepo := _authRepo.New()

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo)

	authH := New(authUCase, usersUCase)

	body := bytes.NewReader([]byte(`{"email":"testing_auth@example.com",
									"password":"4321",
									"repeat_pw":"4321",
									"first_name":"Ivan", 
    								"last_name":"Ivanov", 
								    "birth_date":"29.01.2002"}`))

	exptectedSt := http.StatusOK

	r := httptest.NewRequest("POST", "/api/signup", body)
	w := httptest.NewRecorder()

	authH.SignUp(w, r)
	if w.Code != exptectedSt {
		t.Error("status is not ok")
	}
	cookie := w.Header().Get("Set-Cookie")

	r = httptest.NewRequest("GET", "/api/auth", bytes.NewReader([]byte{}))
	r.Header.Add("Cookie", cookie)
	w = httptest.NewRecorder()
	authH.Logout(w, r)
	if w.Code != exptectedSt {
		t.Error("status is not ok")
	}
}

func TestHandlers_Logout(t *testing.T) {
	t.Parallel()
	userRepo := _userRepo.New()
	authRepo := _authRepo.New()

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo)

	authH := New(authUCase, usersUCase)

	body := bytes.NewReader([]byte(`{"email":"testing_logout@example.com",
									"password":"4321",
									"repeat_pw":"4321",
									"first_name":"Ivan", 
    								"last_name":"Ivanov", 
								    "birth_date":"29.01.2002"}`))

	exptectedSt := http.StatusOK

	r := httptest.NewRequest("POST", "/api/signup", body)
	w := httptest.NewRecorder()

	authH.SignUp(w, r)
	if w.Code != exptectedSt {
		t.Error("status is not ok")
	}
	cookie := w.Header().Get("Set-Cookie")

	r = httptest.NewRequest("GET", "/api/logout", bytes.NewReader([]byte{}))
	r.Header.Add("Cookie", cookie)
	w = httptest.NewRecorder()
	authH.Auth(w, r)
	if w.Code != exptectedSt {
		t.Error("status is not ok")
	}
}
