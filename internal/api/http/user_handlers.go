package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgHttp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/http"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/image"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	pkgErrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type UserHandlersI interface {
	Delete(w http.ResponseWriter, r *http.Request)
	GetInfo(w http.ResponseWriter, r *http.Request)
	GetPersonalInfo(w http.ResponseWriter, r *http.Request)
	EditInfo(w http.ResponseWriter, r *http.Request)
	EditAvatar(w http.ResponseWriter, r *http.Request)
	GetAvatar(w http.ResponseWriter, r *http.Request)
	EditPw(w http.ResponseWriter, r *http.Request)
}

type userHandlers struct {
	cfg    *config.Config
	userUC user.UseCaseI
}

func NewUserHandlers(c *config.Config, uUC user.UseCaseI) UserHandlersI {
	return &userHandlers{
		cfg:    c,
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
func (h *userHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}
	err := h.userUC.Delete(userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
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
// @Param email path string true "email"
// @Success 200 {object} models.UserInfo "success get user info"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/info/{email} [get]
func (h *userHandlers) GetInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars[h.cfg.Routes.RouteUserInfoQueryEmail]
	u, err := h.userUC.GetByEmail(email)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}
	info, err := h.userUC.GetInfo(u.UserID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, info)
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
func (h *userHandlers) GetPersonalInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	u, err := h.userUC.GetByID(userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}
	info, err := h.userUC.GetInfo(u.UserID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, info)
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
// @Router   /user/info [put]
func (h *userHandlers) EditInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(fmt.Errorf("failed to close request: %w", err))
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed read request body"))
		return
	}

	form := &models.UserInfo{}
	if err := easyjson.Unmarshal(body, form); err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	form.Sanitize()

	info, err := h.userUC.EditInfo(userID, form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}
	pkgHttp.SendJSON(w, r, http.StatusOK, models.EditUserInfoResponse{Email: info.Email})
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
// @Router   /user/avatar [put]
func (h *userHandlers) EditAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	//err := r.ParseMultipartForm(config.MaxImageSize)
	//if err != nil {
	//	pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
	//	return
	//}

	file, header, err := r.FormFile(h.cfg.UserService.UserFormNewAvatar)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	img, err := image.ReadImage(file, header)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	err = h.userUC.EditAvatar(userID, img, true)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
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
func (h *userHandlers) GetAvatar(w http.ResponseWriter, r *http.Request) {
	sanitizer := bluemonday.UGCPolicy()
	email := r.URL.Query().Get(h.cfg.Routes.RouteUserAvatarQueryEmail)
	img, err := h.userUC.GetAvatar(sanitizer.Sanitize(email))
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendImage(w, r, http.StatusOK, img.Data)
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
// @Router   /user/pw [put]
func (h *userHandlers) EditPw(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(fmt.Errorf("failed to close request: %w", err))
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed read request body"))
		return
	}

	form := &models.EditPasswordRequest{}
	if err := easyjson.Unmarshal(body, form); err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	form.Sanitize()

	err = h.userUC.EditPw(userID, form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
