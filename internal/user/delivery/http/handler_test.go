package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	userUCMock "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_GetInfo(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := userUCMock.NewMockUseCaseI(ctrl)

	reqEmail := "test@example.com"
	resUser := &models.User{
		UserID:    1,
		Email:     reqEmail,
		Password:  "123456",
		FirstName: "test",
		LastName:  "test",
		Avatar:    "default_avatar",
	}
	resInfo := &models.UserInfo{
		UserID:    0,
		FirstName: "test",
		LastName:  "test",
		Email:     reqEmail,
	}

	r := httptest.NewRequest(http.MethodPost, config.RouteUserInfo, nil)
	vars := make(map[string]string)
	vars[config.RouteUserInfoQueryEmail] = reqEmail
	r = mux.SetURLVars(r, vars)

	userUC.EXPECT().GetByEmail(reqEmail).Return(resUser, nil)
	userUC.EXPECT().GetInfo(resUser.UserID).Return(resInfo, nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	handler := New(userUC)
	router.HandleFunc(config.RouteUserInfo, handler.GetInfo).Methods(http.MethodGet).Queries(config.RouteUserInfoQueryEmail, "{email}")

	handler.GetInfo(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	var respUserInfo models.UserInfo
	err := json.NewDecoder(w.Body).Decode(&respUserInfo)
	require.Nil(t, err)
	fmt.Println(respUserInfo)

	require.Equal(t, resInfo, &respUserInfo)

}

func TestHandlers_Delete(t *testing.T) {

}

func TestHandlers_EditAvatar(t *testing.T) {

}

func TestHandlers_EditInfo(t *testing.T) {

}

func TestHandlers_GetAvatar(t *testing.T) {

}

func TestHandlers_GetPersonalInfo(t *testing.T) {

}
