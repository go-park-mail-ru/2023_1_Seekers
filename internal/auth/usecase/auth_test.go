package usecase

//import (
//	"bytes"
//	"encoding/json"
//	"github.com/go-faker/faker/v4"
//	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
//	mockSessionUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase/mocks_session"
//	mockMailUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase/mocks"
//	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
//	mockUserRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/postgres/mocks"
//	mockUserUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase/mocks"
//	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
//	"github.com/golang/mock/gomock"
//	pkgErr "github.com/pkg/errors"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func generateFakeData(data any) {
//	faker.SetRandomMapAndSliceMaxSize(10)
//	faker.SetRandomMapAndSliceMinSize(1)
//	faker.SetRandomStringLength(30)
//
//	faker.FakeData(data)
//}
//
//func TestUseCase_SignIn(t *testing.T) {
//	var fakeForm models.FormLogin
//	var fakeSession *models.Session
//	var fakeUser *models.User
//	generateFakeData(&fakeForm)
//	generateFakeData(&fakeSession)
//	generateFakeData(&fakeUser)
//	fakeUser.Email = fakeForm.Login + config.PostAtDomain
//	fakeUser.Password = fakeForm.Password
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
//	userRepo := mockUserRepo.NewMockRepoI(ctrl)
//	sUC := mockSessionUC.NewMockSessionUseCaseI(ctrl)
//	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
//	aUC := NewAuthUC(sUC, userRepo, mailUC)
//
//	userRepo.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
//	sUC.EXPECT().CreateSession(fakeUser.UserID).Return(fakeSession, nil)
//
//	responseAuth, responseSession, err := aUC.SignIn(fakeForm)
//	causeErr := pkgErr.Cause(err)
//
//	if causeErr != nil {
//		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
//	} else {
//		require.Equal(t, fakeAuthResponse, responseAuth)
//		require.Equal(t, fakeSession, responseSession)
//	}
//}
//
//func TestUseCase_SignUp(t *testing.T) {
//	var fakeForm models.FormSignUp
//	generateFakeData(&fakeForm)
//	fakeForm.RepeatPw = fakeForm.Password
//	email := fakeForm.Login + config.PostAtDomain
//	fakeUser := &models.User{
//		Email:     email,
//		Password:  fakeForm.Password,
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
//	userRepo := mockUserRepo.NewMockRepoI(ctrl)
//	sUC := mockSessionUC.NewMockSessionUseCaseI(ctrl)
//	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
//	aUC := NewAuthUC(sUC, userRepo, mailUC)
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
//		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
//	} else {
//		require.Equal(t, fakeAuthResponse, responseAuth)
//		require.Equal(t, fakeSession, responseSession)
//	}
//}
//
//func TestDelivery_EditPw(t *testing.T) {
//	var fakePasswordRequest models.EditPasswordRequest
//	generateFakeData(&fakePasswordRequest)
//	userID := uint64(1)
//	status := http.StatusOK
//
//	t.Parallel()
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	userUC := mockUserUC.NewMockUseCaseI(ctrl)
//	userH := New(userUC)
//
//	body, err := json.Marshal(fakePasswordRequest)
//	if err != nil {
//		t.Fatalf("error while marshaling to json: %v", err)
//	}
//
//	r := httptest.NewRequest("POST", "/api/user/info", bytes.NewReader(body))
//	r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, userID))
//	w := httptest.NewRecorder()
//
//	userUC.EXPECT().EditPw(userID, fakePasswordRequest).Return(nil)
//	userH.EditPw(w, r)
//
//	if w.Code != status {
//		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
//	}
//}
