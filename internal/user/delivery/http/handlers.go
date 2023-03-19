package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
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

}

func (h *handlers) GetInfo(w http.ResponseWriter, r *http.Request) {

}

func (h *handlers) EditInfo(w http.ResponseWriter, r *http.Request) {

}

func (h *handlers) EditPw(w http.ResponseWriter, r *http.Request) {

}

func (h *handlers) EditAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleUserErr(w, r, mail.ErrFailedGetUser)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		handleUserErr(w, r, user.ErrInvalidForm)
		return
	}

	file, header, err := r.FormFile(config.RouteUserAvatarFormNewAvatar)
	if err != nil {
		handleUserErr(w, r, user.ErrInvalidForm)
		return
	}

	img := models.Image{
		Data: make([]byte, header.Size),
	}
	_, err = io.ReadFull(file, img.Data)
	if err != nil {
		handleUserErr(w, r, user.ErrInvalidForm)
		return
	}

	if ok = pkg.CheckImageContentType(http.DetectContentType(img.Data)); !ok {
		handleUserErr(w, r, user.ErrWrongContentType)
		return
	}

	img.Name = header.Filename
	err = h.userUC.EditAvatar(userID, &img)
	if err != nil {
		handleUserErr(w, r, user.ErrInternal)
		return
	}

	pkg.SendImage(w, r, http.StatusOK, img.Data)
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
