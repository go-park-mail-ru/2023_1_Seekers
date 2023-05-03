package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	mockAuthUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/usecase/mocks"
	mockMailUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/usecase/mocks"
	mockUserUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createConfig() *config.Config {
	cfg := new(config.Config)
	cfg.Password.PasswordSaltLen = 10
	cfg.Sessions.CookieName = "MailBoxSession"
	ttl, _ := time.ParseDuration("2400h0m0s")
	cfg.Sessions.CookieTTL = ttl
	cfg.Sessions.CookieLen = 32
	cfg.Mail.PostAtDomain = "@mailbox.ru"
	cfg.Routes.RouteUserAvatarQueryEmail = "email"
	cfg.Routes.RouteUserInfoQueryEmail = "email"

	return cfg
}

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestDelivery_SignUp(t *testing.T) {
	cfg := createConfig()

	var fakeForm *models.FormSignUp
	var fakeSession *models.Session
	status := http.StatusOK
	generateFakeData(&fakeForm)
	email := fakeForm.Login + cfg.Mail.PostAtDomain
	fakeUser := &models.UserInfo{
		UserID:    uint64(1),
		Email:     email,
		FirstName: fakeForm.FirstName,
		LastName:  fakeForm.LastName,
	}
	generateFakeData(&fakeSession)
	fakeAuthResponse := &models.AuthResponse{
		Email:     fakeUser.Email,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUC := mockAuthUC.NewMockUseCaseI(ctrl)
	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	authH := NewAuthHandlers(cfg, authUC, mailUC, userUC)

	body, err := json.Marshal(fakeForm)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest("POST", "/api/signup", bytes.NewReader(body))
	w := httptest.NewRecorder()

	authUC.EXPECT().SignUp(fakeForm).Return(fakeAuthResponse, fakeSession, nil)
	userUC.EXPECT().GetInfoByEmail(fakeUser.Email).Return(fakeUser, nil)
	mailUC.EXPECT().CreateDefaultFolders(fakeUser.UserID).Return([]models.Folder{}, nil)
	mailUC.EXPECT().SendWelcomeMessage(fakeUser.Email).Return(nil)
	authH.SignUp(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_SignIn(t *testing.T) {
	cfg := createConfig()

	var fakeForm *models.FormLogin
	var fakeAuthResponse *models.AuthResponse
	var fakeSession *models.Session
	status := http.StatusOK
	generateFakeData(&fakeForm)
	generateFakeData(&fakeAuthResponse)
	generateFakeData(&fakeSession)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUC := mockAuthUC.NewMockUseCaseI(ctrl)
	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	authH := NewAuthHandlers(cfg, authUC, mailUC, userUC)

	body, err := json.Marshal(fakeForm)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest("POST", "/api/signin", bytes.NewReader(body))
	w := httptest.NewRecorder()

	authUC.EXPECT().SignIn(fakeForm).Return(fakeAuthResponse, fakeSession, nil)
	authH.SignIn(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_Logout(t *testing.T) {
	cfg := createConfig()

	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUC := mockAuthUC.NewMockUseCaseI(ctrl)
	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	authH := NewAuthHandlers(cfg, authUC, mailUC, userUC)

	r := httptest.NewRequest("POST", "/api/signin", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	authH.Logout(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_Auth(t *testing.T) {
	cfg := createConfig()

	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUC := mockAuthUC.NewMockUseCaseI(ctrl)
	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	authH := NewAuthHandlers(cfg, authUC, mailUC, userUC)

	r := httptest.NewRequest("GET", "/api/auth", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	authH.Auth(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_GetCSRF(t *testing.T) {
	cfg := createConfig()

	var fakeSession *models.Session
	generateFakeData(&fakeSession)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUC := mockAuthUC.NewMockUseCaseI(ctrl)
	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	authH := NewAuthHandlers(cfg, authUC, mailUC, userUC)

	r := httptest.NewRequest("POST", "/api/create_csrf", bytes.NewReader([]byte{}))
	r.AddCookie(&http.Cookie{
		Name:     cfg.Sessions.CookieName,
		Value:    fakeSession.SessionID,
		Expires:  time.Now().Add(cfg.Sessions.CookieTTL),
		HttpOnly: true,
		Path:     cfg.Sessions.CookiePath,
		SameSite: http.SameSiteLaxMode,
	})
	w := httptest.NewRecorder()

	authH.GetCSRF(w, r)
	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}
