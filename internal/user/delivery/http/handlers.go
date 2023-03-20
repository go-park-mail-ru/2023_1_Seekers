package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type handlers struct {
	userUC user.UseCaseI
}

func New(uUC user.UseCaseI) user.HandlersI {
	return &handlers{
		userUC: uUC,
	}
}

func handleUserErr(w http.ResponseWriter, r *http.Request, err error) {
	pkg.HandleError(w, r, user.Errors[err], err)
}

func (h *handlers) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleUserErr(w, r, user.ErrFailedGetUser)
		return
	}
	err := h.userUC.Delete(userID)
	if err != nil {
		handleUserErr(w, r, user.ErrFailedDelete)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) GetInfo(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get(config.RouteUserInfoQueryEmail)
	u, err := h.userUC.GetByEmail(email)
	if err != nil {
		handleUserErr(w, r, user.ErrFailedGetUser)
		return
	}
	info, err := h.userUC.GetInfo(u.ID)
	if err != nil {
		handleUserErr(w, r, user.ErrInternal)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, info)
}

func (h *handlers) EditInfo(w http.ResponseWriter, r *http.Request) {
	// тут пока что просто из body - в будущем на form data
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleUserErr(w, r, user.ErrFailedGetUser)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(fmt.Errorf("failed to close request: %w", err))
		}
	}(r.Body)
	form := models.UserInfo{}

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		handleUserErr(w, r, user.ErrInvalidForm)
		return
	}
	info, err := h.userUC.EditInfo(userID, form)
	if err != nil {
		handleUserErr(w, r, user.ErrFailedEditInfo)
		return
	}
	pkg.SendJSON(w, r, http.StatusOK, models.EditUserInfoResponse{Email: info.Email})
}

func (h *handlers) EditPw(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleUserErr(w, r, user.ErrFailedGetUser)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(fmt.Errorf("failed to close request: %w", err))
		}
	}(r.Body)
	form := models.EditPasswordRequest{}

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		handleUserErr(w, r, user.ErrInvalidForm)
		return
	}

	err = h.userUC.EditPw(userID, form)
	if err != nil {
		handleUserErr(w, r, user.ErrFailedEditPw)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) EditAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleUserErr(w, r, user.ErrFailedGetUser)
		return
	}

	err := r.ParseMultipartForm(config.MaxImageSize)
	if err != nil {
		handleUserErr(w, r, user.ErrInvalidForm)
		return
	}

	file, header, err := r.FormFile(config.UserFormNewAvatar)
	if err != nil {
		handleUserErr(w, r, user.ErrInvalidForm)
		return
	}

	img, err := pkg.ReadImage(file, header)
	if err != nil {
		if err == user.ErrInvalidForm || err == user.ErrWrongContentType {
			handleUserErr(w, r, err)
		} else {
			handleUserErr(w, r, user.ErrInternal)
		}
		return
	}

	err = h.userUC.EditAvatar(userID, img)
	if err != nil {
		handleUserErr(w, r, user.ErrInternal)
		return
	}

	//pkg.SendImage(w, r, http.StatusOK, img.Data)
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) GetAvatar(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get(config.RouteUserAvatarQueryEmail)
	img, err := h.userUC.GetAvatar(email)
	if err != nil {
		handleUserErr(w, r, user.ErrInternal)
		return
	}

	pkg.SendImage(w, r, http.StatusOK, img.Data)
}
