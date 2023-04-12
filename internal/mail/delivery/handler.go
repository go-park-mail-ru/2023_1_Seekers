package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgHttp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	pkgErrors "github.com/pkg/errors"
	"net/http"
	"strconv"
)

type delivery struct {
	uc mail.UseCaseI
}

func New(uc mail.UseCaseI) mail.HandlersI {
	return &delivery{
		uc: uc,
	}
}

// GetFolderMessages godoc
// @Summary      GetFolderMessages
// @Description  List of folder messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Param slug path string true "FolderSlug"
// @Success  200 {object} models.FolderResponse "success get list of folder messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid url address"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /folder/{slug} [get]
func (del *delivery) GetFolderMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	folderSlug, ok := vars["slug"]
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	folder, err := del.uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	messages, err := del.uc.GetFolderMessages(userID, folderSlug)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.FolderResponse{
		Folder:   *folder,
		Messages: messages,
	})
}

// GetFolders godoc
// @Summary      GetFolders
// @Description  List of outgoing messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.FoldersResponse "success get list of outgoing messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /folders/ [get]
func (del *delivery) GetFolders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	folders, err := del.uc.GetFolders(userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.FoldersResponse{
		Folders: folders,
		Count:   len(folders),
	})
}

// GetMessage godoc
// @Summary      GetMessage
// @Description  Message
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Success  200 {object} models.MessageResponse "success get messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid url address"
// @Failure 404 {object} errors.JSONError "message not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/{id} [get]
func (del *delivery) GetMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	messageID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkgHttp.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	message, err := del.uc.GetMessage(userID, messageID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

// SendMessage godoc
// @Summary      SendMessage
// @Description  Message
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessageResponse "success send message"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "no valid emails"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 404 {object} errors.JSONError "message not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/send [post]
func (del *delivery) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	form := models.FormMessage{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	form.Sanitize()

	validEmails, invalidEmails := del.uc.ValidateRecipients(form.Recipients)
	form.Recipients = validEmails

	message, err := del.uc.SendMessage(userID, form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	if len(invalidEmails) != 0 {
		err = del.uc.SendFailedSendingMessage(message.FromUser.Email, invalidEmails)

		if err != nil {
			pkgHttp.HandleError(w, r, err)
			return
		}
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

// ReadMessage godoc
// @Summary      ReadMessage
// @Description  Message
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Success  200 {object} models.MessageResponse "success read message"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid url address"
// @Failure 404 {object} errors.JSONError "message not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/{id}/read [post]
func (del *delivery) ReadMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	messageID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkgHttp.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	message, err := del.uc.MarkMessageAsSeen(userID, messageID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

// UnreadMessage godoc
// @Summary      UnreadMessage
// @Description  Message
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Success  200 {object} models.MessageResponse "success unread message"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid url address"
// @Failure 404 {object} errors.JSONError "message not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/{id}/unread [post]
func (del *delivery) UnreadMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	messageID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkgHttp.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	message, err := del.uc.MarkMessageAsUnseen(userID, messageID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}
