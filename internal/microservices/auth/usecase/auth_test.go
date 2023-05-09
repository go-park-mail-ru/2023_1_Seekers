package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"strings"

	//"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	mockSessionRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/repository/mocks"
	mockUserUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/crypto"
	"github.com/golang/mock/gomock"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func generateFakeData(data any) {
	_ = faker.SetRandomMapAndSliceMaxSize(10)
	_ = faker.SetRandomMapAndSliceMinSize(1)
	_ = faker.SetRandomStringLength(20)

	_ = faker.FakeData(data)
}

func createConfig() *config.Config {
	cfg := new(config.Config)
	cfg.Mail.PostAtDomain = "@mailbx.ru"
	cfg.Password.PasswordSaltLen = 0
	cfg.UserService.DefaultAvatar = "default_avatar.png"

	return cfg
}

func TestUseCase_SignIn(t *testing.T) {
	cfg := createConfig()

	var fakeForm *models.FormLogin
	var fakeSession *models.Session
	var fakeUser *models.User
	generateFakeData(&fakeForm)
	generateFakeData(&fakeSession)
	generateFakeData(&fakeUser)
	fakeForm.Login = strings.ToLower(fakeForm.Login)
	fakeUser.Email = fakeForm.Login + cfg.Mail.PostAtDomain
	fakeUser.IsExternal = false

	fakeUser.Password = string(crypto.Hash([]byte{}, fakeForm.Password))

	fakeAuthResponse := &models.AuthResponse{
		Email:     fakeUser.Email,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRepo := mockSessionRepo.NewMockSessionRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	aUC := NewAuthUC(cfg, userUC, sessionRepo)

	userUC.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
	sessionRepo.EXPECT().CreateSession(fakeUser.UserID).Return(fakeSession, nil)

	responseAuth, responseSession, err := aUC.SignIn(fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, err)
	} else {
		require.Equal(t, fakeAuthResponse, responseAuth)
		require.Equal(t, fakeSession, responseSession)
	}
}

func TestUseCase_SignUp(t *testing.T) {
	cfg := createConfig()

	var fakeForm *models.FormSignUp
	generateFakeData(&fakeForm)
	fakeForm.RepeatPw = fakeForm.Password
	fakeForm.Login = strings.ToLower(fakeForm.Login)
	email := fakeForm.Login + cfg.Mail.PostAtDomain
	hashPW, err := crypto.HashPw(fakeForm.Password, cfg.Password.PasswordSaltLen)
	if err != nil {
		t.Fatalf("error while hashing pw")
	}
	fakeUser := &models.User{
		Email:     email,
		Password:  hashPW,
		FirstName: fakeForm.FirstName,
		LastName:  fakeForm.LastName,
		Avatar:    cfg.UserService.DefaultAvatar,
	}
	var fakeSession *models.Session
	generateFakeData(&fakeSession)
	fakeAuthResponse := &models.AuthResponse{
		Email:     fakeUser.Email,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRepo := mockSessionRepo.NewMockSessionRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	aUC := NewAuthUC(cfg, userUC, sessionRepo)

	userUC.EXPECT().Create(fakeUser).Return(fakeUser, nil)
	sessionRepo.EXPECT().CreateSession(fakeUser.UserID).Return(fakeSession, nil)

	responseAuth, responseSession, err := aUC.SignUp(fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeAuthResponse, responseAuth)
		require.Equal(t, fakeSession, responseSession)
	}
}

func TestUseCase_CreateSession(t *testing.T) {
	cfg := createConfig()

	var fakeSession *models.Session
	generateFakeData(&fakeSession)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRepo := mockSessionRepo.NewMockSessionRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	aUC := NewAuthUC(cfg, userUC, sessionRepo)

	sessionRepo.EXPECT().CreateSession(fakeSession.UID).Return(fakeSession, nil)
	response, err := aUC.CreateSession(fakeSession.UID)

	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeSession, response)
	}
}

func TestUseCase_DeleteSession(t *testing.T) {
	cfg := createConfig()

	sessionID := "adjfkfldlkld"

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRepo := mockSessionRepo.NewMockSessionRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	aUC := NewAuthUC(cfg, userUC, sessionRepo)

	sessionRepo.EXPECT().DeleteSession(sessionID).Return(nil)
	err := aUC.DeleteSession(sessionID)

	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestUseCase_GetSession(t *testing.T) {
	cfg := createConfig()

	var fakeSession *models.Session
	generateFakeData(&fakeSession)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRepo := mockSessionRepo.NewMockSessionRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	aUC := NewAuthUC(cfg, userUC, sessionRepo)

	sessionRepo.EXPECT().GetSession(fakeSession.SessionID).Return(fakeSession, nil)
	response, err := aUC.GetSession(fakeSession.SessionID)

	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeSession, response)
	}
}
