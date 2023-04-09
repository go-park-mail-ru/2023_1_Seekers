package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	mockSessionUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase/mocks_session"
	mockMailUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	mockUserRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/mocks"
	mockUserUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
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
	config.PasswordSaltLen = 0
	var fakeForm models.FormLogin
	var fakeSession *models.Session
	var fakeUser *models.User
	generateFakeData(&fakeForm)
	generateFakeData(&fakeSession)
	generateFakeData(&fakeUser)
	fakeUser.Email = fakeForm.Login + config.PostAtDomain
	var err error
	fakeUser.Password, err = pkg.HashPw(fakeForm.Password)
	if err != nil {
		t.Fatalf("error while hashing pw")
	}
	fakeAuthResponse := &models.AuthResponse{
		Email:     fakeUser.Email,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sUC := mockSessionUC.NewMockSessionUseCaseI(ctrl)
	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	userUc := mockUserUC.NewMockUseCaseI(ctrl)
	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	aUC := NewAuthUC(sUC, userRepo, mailUC, userUc)

	userRepo.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
	sUC.EXPECT().CreateSession(fakeUser.UserID).Return(fakeSession, nil)

	responseAuth, responseSession, err := aUC.SignIn(fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeAuthResponse, responseAuth)
		require.Equal(t, fakeSession, responseSession)
	}
}

//
//func TestUseCase_SignUp(t *testing.T) {
//	config.PasswordSaltLen = 0
//	var fakeForm models.FormSignUp
//	generateFakeData(&fakeForm)
//	fakeForm.RepeatPw = fakeForm.Password
//	email := fakeForm.Login + config.PostAtDomain
//	hashPW, err := pkg.HashPw(fakeForm.Password)
//	if err != nil {
//		t.Fatalf("error while hashing pw")
//	}
//	fakeUser := &models.User{
//		Email:     email,
//		Password:  hashPW,
//		FirstName: fakeForm.FirstName,
//		LastName:  fakeForm.LastName,
//		Avatar:    config.DefaultAvatar,
//	}
//	var fakeSession *models.Session
//	generateFakeData(&fakeSession)
//	fakeAuthResponse := &models.AuthResponse{
//		Email:     fakeUser.Email,
//		FirstName: fakeUser.FirstName,
//		LastName:  fakeUser.LastName,
//	}
//
//	t.Parallel()
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	sUC := mockSessionUC.NewMockSessionUseCaseI(ctrl)
//	userRepo := mockUserRepo.NewMockRepoI(ctrl)
//	userUc := mockUserUC.NewMockUseCaseI(ctrl)
//	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
//	aUC := NewAuthUC(sUC, userRepo, mailUC, userUc)
//
//	userRepo.EXPECT().Create(fakeUser).Return(fakeUser, nil)
//	mailUC.EXPECT().CreateDefaultFolders(fakeUser.UserID).Return([]models.Folder{}, nil)
//	mailUC.EXPECT().SendWelcomeMessage(fakeUser.Email).Return(nil)
//	sUC.EXPECT().CreateSession(fakeUser.UserID).Return(fakeSession, nil)
//
//	responseAuth, responseSession, err := aUC.SignUp(fakeForm)
//	causeErr := pkgErr.Cause(err)
//
//	if causeErr != nil {
//		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
//	} else {
//		require.Equal(t, fakeAuthResponse, responseAuth)
//		require.Equal(t, fakeSession, responseSession)
//	}
//}

func TestUseCase_EditPw(t *testing.T) {
	config.PasswordSaltLen = 0
	var fakeForm models.EditPasswordRequest
	generateFakeData(&fakeForm)
	fakeForm.RepeatPw = fakeForm.Password
	var fakeUser *models.User
	generateFakeData(&fakeUser)
	var err error
	fakeUser.Password, err = pkg.HashPw(fakeForm.PasswordOld)
	if err != nil {
		t.Fatalf("error while hashing pw")
	}
	newPw, err := pkg.HashPw(fakeForm.Password)
	if err != nil {
		t.Fatalf("error while hashing pw")
	}
	userID := uint64(1)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sUC := mockSessionUC.NewMockSessionUseCaseI(ctrl)
	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	userUc := mockUserUC.NewMockUseCaseI(ctrl)
	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	aUC := NewAuthUC(sUC, userRepo, mailUC, userUc)

	userRepo.EXPECT().GetByID(userID).Return(fakeUser, nil)
	userRepo.EXPECT().EditPw(userID, newPw).Return(nil)

	err = aUC.EditPw(userID, fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}
