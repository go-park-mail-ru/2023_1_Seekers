package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
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

// Delete godoc
// @Summary      Delete
// @Description  delete user
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Success  200 "success delete user"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user [delete]
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

// GetInfo godoc
// @Summary      GetInfo
// @Description  get info about user
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Param id query string true "email"
// @Success 200 {object} models.UserInfo "success get user info"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/info [get]
func (h *handlers) GetInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars[config.RouteUserInfoQueryEmail]
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

// GetPersonalInfo godoc
// @Summary      GetPersonalInfo
// @Description  get info about request creator
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Success 200 {object} models.UserInfo "success get user info"
// @Failure 401 {object} errors.JSONError "failed get user"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/info [get]
func (h *handlers) GetPersonalInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		pkg.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	u, err := h.userUC.GetByID(userID)
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

// EditInfo godoc
// @Summary      EditInfo
// @Description  edit info about user
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Success 200 {object} models.EditUserInfoResponse "success edit user info"
// @Failure 401 {object} errors.JSONError "failed to get user"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/info [post]
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

	form.Sanitize()

	info, err := h.userUC.EditInfo(userID, form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}
	pkg.SendJSON(w, r, http.StatusOK, models.EditUserInfoResponse{Email: info.Email})
}

// EditPw godoc
// @Summary      EditPw
// @Description  edit password about user
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Success 200 "success edit user password"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/pw [post]
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

	form.Sanitize()

	err = h.userUC.EditPw(userID, form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// EditAvatar godoc
// @Summary      EditAvatar
// @Description  edit user avatar
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Success 200 "success edit user avatar"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "unsupported content type"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/avatar [post]
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

// GetAvatar godoc
// @Summary      GetAvatar
// @Description  get user avatar
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Param id query string true "email"
// @Success 200 {object} []byte "success get user avatar"
// @Failure 400 {object} errors.JSONError "failed get file"
// @Failure 400 {object} errors.JSONError "no key"
// @Failure 400 {object} errors.JSONError "no bucket"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/avatar [get]
func (h *handlers) GetAvatar(w http.ResponseWriter, r *http.Request) {
	sanitizer := bluemonday.UGCPolicy()
	email := r.URL.Query().Get(config.RouteUserAvatarQueryEmail)
	img, err := h.userUC.GetAvatar(sanitizer.Sanitize(email))
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendImage(w, r, http.StatusOK, img.Data)
}
