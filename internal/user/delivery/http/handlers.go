package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
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

func (h *handlers) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		pkg.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}
	err := h.userUC.Delete(userID)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) GetInfo(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get(config.RouteUserInfoQueryEmail)
	u, err := h.userUC.GetByEmail(email)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}
	info, err := h.userUC.GetInfo(u.UserID)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, info)
}

func (h *handlers) EditInfo(w http.ResponseWriter, r *http.Request) {
	// тут пока что просто из body - в будущем на form data
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		pkg.HandleError(w, r, errors.ErrFailedGetUser)
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
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}
	info, err := h.userUC.EditInfo(userID, form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}
	pkg.SendJSON(w, r, http.StatusOK, models.EditUserInfoResponse{Email: info.Email})
}

func (h *handlers) EditPw(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		pkg.HandleError(w, r, errors.ErrFailedGetUser)
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
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	err = h.userUC.EditPw(userID, form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) EditAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		pkg.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	err := r.ParseMultipartForm(config.MaxImageSize)
	if err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	file, header, err := r.FormFile(config.UserFormNewAvatar)
	if err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	img, err := pkg.ReadImage(file, header)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	err = h.userUC.EditAvatar(userID, img)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	//pkg.SendImage(w, r, http.StatusOK, img.Data)
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) GetAvatar(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get(config.RouteUserAvatarQueryEmail)
	img, err := h.userUC.GetAvatar(email)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendImage(w, r, http.StatusOK, img.Data)
}
