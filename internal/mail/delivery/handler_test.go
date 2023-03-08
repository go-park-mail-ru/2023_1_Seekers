package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/delivery/http"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/repository/inmemory"
	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type inputCase struct {
	userCr   string
	folderID uint64
}

type outputCase struct {
	status int
}

type testCases struct {
	name   string
	input  inputCase
	output outputCase
}

var credentials = map[string][]byte{
	"user1": []byte(`{"login": "test@mailbox.ru", "password": "12345"}`),
	"user2": []byte(`{"login": "gena@mailbox.ru", "password": "54321"}`),
	"user3": []byte(`{"login": "max@mailbox.ru", "password": "13795"}`),
	"user4": []byte(`{"login": "valera@mailbox.ru", "password": "12345"}`),
}

func prepare(t *testing.T, r *http.Request, testName string, userCr string) (*_middleware.Middleware, mail.HandlersI) {
	userRepo := _userRepo.New()
	authRepo := _authRepo.New()
	mailRepo := _mailRepo.New(userRepo)

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo, usersUCase)
	mailUCase := _mailUCase.New(mailRepo)

	authH := _authHandler.New(authUCase, usersUCase, mailUCase)
	mailH := New(mailUCase)
	logger := pkg.GetLogger()
	middleware := _middleware.New(authUCase, logger)

	signinReq := httptest.NewRequest("POST", "/api/signin", bytes.NewReader(credentials[userCr]))
	w := httptest.NewRecorder()

	authH.SignIn(w, signinReq)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed login for test \"%s\" with code %d", testName, w.Code)
	}

	var user models.SignInResponse
	json.NewDecoder(w.Body).Decode(&user)

	s, err := authUCase.GetSessionByEmail(user.Email)
	if err != nil {
		t.Fatalf("Failed to get session %v for test \"%s\"", err, testName)
	}

	r.AddCookie(&http.Cookie{
		Name:  config.CookieName,
		Value: s.SessionID,
	})

	return middleware, mailH
}

func TestDelivery_GetInboxMessages(t *testing.T) {
	var tests = []testCases{
		{
			name: "there are incoming messages user1",
			input: inputCase{
				userCr: "user1",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "there are incoming messages user2",
			input: inputCase{
				userCr: "user2",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "no incoming messages user4",
			input: inputCase{
				userCr: "user4",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()

	for _, test := range tests {
		r := httptest.NewRequest("GET", "/api/inbox", bytes.NewReader([]byte{}))
		w := httptest.NewRecorder()
		middleware, mailH := prepare(t, r, test.name, test.input.userCr)

		middleware.CheckAuth(mailH.GetInboxMessages)(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_GetOutboxMessages(t *testing.T) {
	var tests = []testCases{
		{
			name: "there are outgoing messages user1",
			input: inputCase{
				userCr: "user1",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "there are outgoing messages user2",
			input: inputCase{
				userCr: "user2",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "no outgoing messages user4",
			input: inputCase{
				userCr: "user4",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()

	for _, test := range tests {
		r := httptest.NewRequest("GET", "/api/outbox", bytes.NewReader([]byte{}))
		w := httptest.NewRecorder()
		middleware, mailH := prepare(t, r, test.name, test.input.userCr)

		middleware.CheckAuth(mailH.GetOutboxMessages)(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_GetFolderMessages(t *testing.T) {
	var tests = []testCases{
		{
			name: "there are messages in the folder",
			input: inputCase{
				userCr:   "user2",
				folderID: 7,
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "there are no messages in the folder",
			input: inputCase{
				userCr:   "user3",
				folderID: 8,
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "the folder does not exist",
			input: inputCase{
				userCr:   "user1",
				folderID: 7,
			},
			output: outputCase{
				status: http.StatusBadRequest,
			},
		},
	}

	t.Parallel()

	for _, test := range tests {
		r := httptest.NewRequest("GET", "/api/folder/", bytes.NewReader([]byte{}))
		r = mux.SetURLVars(r, map[string]string{"id": strconv.FormatUint(test.input.folderID, 10)})
		w := httptest.NewRecorder()
		middleware, mailH := prepare(t, r, test.name, test.input.userCr)

		middleware.CheckAuth(mailH.GetFolderMessages)(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_GetFolders(t *testing.T) {
	var tests = []testCases{
		{
			name: "only default folders",
			input: inputCase{
				userCr: "user1",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "there are not only default folders",
			input: inputCase{
				userCr: "user2",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()

	for _, test := range tests {
		r := httptest.NewRequest("GET", "/api/folders", bytes.NewReader([]byte{}))
		w := httptest.NewRecorder()
		middleware, mailH := prepare(t, r, test.name, test.input.userCr)

		middleware.CheckAuth(mailH.GetFolders)(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected status %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}
