package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	mockUserUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestDelivery_Delete(t *testing.T) {
	userID := uint64(1)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(userUC)

	r := httptest.NewRequest("DELETE", "/api/user/", bytes.NewReader([]byte{}))
	r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, userID))
	w := httptest.NewRecorder()

	userUC.EXPECT().Delete(userID).Return(nil)
	userH.Delete(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_GetInfo(t *testing.T) {
	var fakeUser *models.User
	generateFakeData(&fakeUser)
	fakeUserInfo := &models.UserInfo{
		UserID:    fakeUser.UserID,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
		Email:     fakeUser.Email,
	}
	email := fakeUser.Email
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(userUC)

	r := httptest.NewRequest("GET", "/api/user/", bytes.NewReader([]byte{}))
	q := r.URL.Query()
	q.Add(config.RouteUserInfoQueryEmail, email)
	r.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	userUC.EXPECT().GetByEmail(email).Return(fakeUser, nil)
	userUC.EXPECT().GetInfo(fakeUser.UserID).Return(fakeUserInfo, nil)
	userH.GetInfo(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_GetPersonalInfo(t *testing.T) {
	var fakeUser *models.User
	generateFakeData(&fakeUser)
	fakeUserInfo := &models.UserInfo{
		UserID:    fakeUser.UserID,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
		Email:     fakeUser.Email,
	}
	userID := uint64(1)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(userUC)

	r := httptest.NewRequest("GET", "/api/user/info", bytes.NewReader([]byte{}))
	r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, userID))
	w := httptest.NewRecorder()

	userUC.EXPECT().GetByID(userID).Return(fakeUser, nil)
	userUC.EXPECT().GetInfo(fakeUser.UserID).Return(fakeUserInfo, nil)
	userH.GetPersonalInfo(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_EditInfo(t *testing.T) {
	var fakeUserInfo *models.UserInfo
	generateFakeData(&fakeUserInfo)
	fakeUserInfo.UserID = 0
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(userUC)

	body, err := json.Marshal(fakeUserInfo)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest("POST", "/api/user/info", bytes.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, fakeUserInfo.UserID))
	w := httptest.NewRecorder()

	userUC.EXPECT().EditInfo(fakeUserInfo.UserID, *fakeUserInfo).Return(fakeUserInfo, nil)
	userH.EditInfo(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_EditPw(t *testing.T) {
	var fakePasswordRequest models.EditPasswordRequest
	generateFakeData(&fakePasswordRequest)
	userID := uint64(1)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(userUC)

	body, err := json.Marshal(fakePasswordRequest)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest("POST", "/api/user/info", bytes.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, userID))
	w := httptest.NewRecorder()

	userUC.EXPECT().EditPw(userID, fakePasswordRequest).Return(nil)
	userH.EditPw(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_EditAvatar(t *testing.T) {
	// TODO
}

func TestDelivery_GetAvatar(t *testing.T) {
	var fakeImage *models.Image
	var email string
	generateFakeData(&fakeImage)
	generateFakeData(email)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(userUC)

	r := httptest.NewRequest("GET", "/api/user/avatar", bytes.NewReader([]byte{}))
	q := r.URL.Query()
	q.Add(config.RouteUserInfoQueryEmail, email)
	r.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	userUC.EXPECT().GetAvatar(email).Return(fakeImage, nil)
	userH.GetAvatar(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}
