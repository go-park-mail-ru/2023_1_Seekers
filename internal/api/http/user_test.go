package http

import (
	"bytes"
	"context"
	mockUserUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
)

func TestDelivery_Delete(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := NewUserHandlers(cfg, userUC)

	r := httptest.NewRequest(http.MethodDelete, "/api/user/", bytes.NewReader([]byte{}))
	r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, userID))
	w := httptest.NewRecorder()

	userUC.EXPECT().Delete(userID).Return(nil)
	userH.Delete(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_GetInfo(t *testing.T) {
	cfg := createConfig()

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
	userH := NewUserHandlers(cfg, userUC)

	r := httptest.NewRequest(http.MethodGet, "/api/user/", bytes.NewReader([]byte{}))
	vars := map[string]string{
		cfg.Routes.RouteUserInfoQueryEmail: email,
	}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	userUC.EXPECT().GetByEmail(email).Return(fakeUser, nil)
	userUC.EXPECT().GetInfo(fakeUser.UserID).Return(fakeUserInfo, nil)
	userH.GetInfo(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_GetPersonalInfo(t *testing.T) {
	cfg := createConfig()

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
	userH := NewUserHandlers(cfg, userUC)

	r := httptest.NewRequest(http.MethodGet, "/api/user/info", bytes.NewReader([]byte{}))
	r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, userID))
	w := httptest.NewRecorder()

	userUC.EXPECT().GetByID(userID).Return(fakeUser, nil)
	userUC.EXPECT().GetInfo(fakeUser.UserID).Return(fakeUserInfo, nil)
	userH.GetPersonalInfo(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_EditInfo(t *testing.T) {
	cfg := createConfig()

	var fakeUserInfo *models.UserInfo
	generateFakeData(&fakeUserInfo)
	fakeUserInfo.UserID = 0
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := NewUserHandlers(cfg, userUC)

	body, err := easyjson.Marshal(fakeUserInfo)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/user/info", bytes.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, fakeUserInfo.UserID))
	w := httptest.NewRecorder()

	userUC.EXPECT().EditInfo(fakeUserInfo.UserID, fakeUserInfo).Return(fakeUserInfo, nil)
	userH.EditInfo(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_EditAvatar(t *testing.T) {
	cfg := createConfig()

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)

	img := []byte("sdsdjbasd;jrandombytes")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Add("Content-Disposition", `form-data; name="avatar"; filename="test.png"`)
	metadataHeader.Add("Content-Type", "image/png")
	part, err := writer.CreatePart(metadataHeader)
	require.Nil(t, err)
	_, err = part.Write(img)
	err = writer.Close()
	require.Nil(t, err)

	r := httptest.NewRequest(http.MethodPut, cfg.Routes.RouteUserAvatar, body)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	user := models.User{
		UserID: 1,
	}

	r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, user.UserID))

	w := httptest.NewRecorder()

	router := mux.NewRouter()
	handler := NewUserHandlers(cfg, userUC)
	router.HandleFunc(cfg.Routes.RouteUserAvatar, handler.EditAvatar).Methods(http.MethodPut)
	handler.EditAvatar(w, r)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestDelivery_GetAvatar(t *testing.T) {
	cfg := createConfig()

	var fakeImage *models.Image
	var email string
	generateFakeData(&fakeImage)
	generateFakeData(email)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := NewUserHandlers(cfg, userUC)

	r := httptest.NewRequest(http.MethodGet, "/api/user/avatar", bytes.NewReader([]byte{}))
	q := r.URL.Query()
	q.Add(cfg.Routes.RouteUserAvatarQueryEmail, email)
	r.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	userUC.EXPECT().GetAvatar(email).Return(fakeImage, nil)
	userH.GetAvatar(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}

func TestDelivery_EditPw(t *testing.T) {
	cfg := createConfig()

	var fakeForm *models.EditPasswordRequest
	generateFakeData(&fakeForm)
	userID := uint64(1)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := NewUserHandlers(cfg, userUC)

	body, err := easyjson.Marshal(fakeForm)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/user/pw", bytes.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, userID))
	w := httptest.NewRecorder()

	userUC.EXPECT().EditPw(userID, fakeForm).Return(nil)
	userH.EditPw(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", status, w.Code)
	}
}
