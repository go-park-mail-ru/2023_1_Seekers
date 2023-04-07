package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-faker/faker/v4"
	mockAuthUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase/mocks_auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
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

func TestDelivery_SignUp(t *testing.T) {
	var fakeForm models.FormSignUp
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
	authH := New(authUC)

	body, err := json.Marshal(fakeForm)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest("POST", "/api/signup", bytes.NewReader(body))
	w := httptest.NewRecorder()

	authUC.EXPECT().SignUp(fakeForm).Return(fakeAuthResponse, fakeSession, nil)
	authH.SignUp(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_SignIn(t *testing.T) {
	var fakeForm models.FormLogin
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
	authH := New(authUC)

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
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUC := mockAuthUC.NewMockUseCaseI(ctrl)
	authH := New(authUC)

	r := httptest.NewRequest("POST", "/api/signin", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	authH.Logout(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}
