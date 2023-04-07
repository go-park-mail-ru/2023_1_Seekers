package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	mockSessionUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase/mocks_session"
	mockMailUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	mockUserUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase/mocks"
	"github.com/golang/mock/gomock"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestUseCase_SignIn(t *testing.T) {
	var fakeForm models.FormLogin
	var fakeSession *models.Session
	var fakeUser *models.User
	generateFakeData(&fakeForm)
	generateFakeData(&fakeSession)
	generateFakeData(&fakeUser)
	fakeUser.Email = fakeForm.Login + config.PostAtDomain
	fakeUser.Password = fakeForm.Password
	fakeAuthResponse := &models.AuthResponse{
		Email:     fakeUser.Email,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	sUC := mockSessionUC.NewMockSessionUseCaseI(ctrl)
	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	aUC := NewAuthUC(sUC, userUC, mailUC)

	userUC.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
	sUC.EXPECT().CreateSession(fakeUser.UserID).Return(fakeSession, nil)

	responseAuth, responseSession, err := aUC.SignIn(fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeAuthResponse, responseAuth)
		require.Equal(t, fakeSession, responseSession)
	}
}

func TestUseCase_SignUp(t *testing.T) {
	var fakeForm models.FormSignUp
	generateFakeData(&fakeForm)
	fakeForm.RepeatPw = fakeForm.Password
	email := fakeForm.Login + config.PostAtDomain
	fakeUser := &models.User{
		Email:     email,
		Password:  fakeForm.Password,
		FirstName: fakeForm.FirstName,
		LastName:  fakeForm.LastName,
		Avatar:    config.DefaultAvatar,
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

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	sUC := mockSessionUC.NewMockSessionUseCaseI(ctrl)
	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	aUC := NewAuthUC(sUC, userUC, mailUC)

	userUC.EXPECT().Create(fakeUser).Return(fakeUser, nil)
	mailUC.EXPECT().CreateDefaultFolders(fakeUser.UserID).Return([]models.Folder{}, nil)
	mailUC.EXPECT().SendWelcomeMessage(fakeUser.Email).Return(nil)
	sUC.EXPECT().CreateSession(fakeUser.UserID).Return(fakeSession, nil)

	responseAuth, responseSession, err := aUC.SignUp(fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeAuthResponse, responseAuth)
		require.Equal(t, fakeSession, responseSession)
	}
}
