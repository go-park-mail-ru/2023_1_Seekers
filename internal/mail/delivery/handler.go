package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	_errors "github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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

func handleMailErr(w http.ResponseWriter, r *http.Request, err error) {
	unwrappedErr := _errors.UnwrapError(err)
	pkg.HandleError(w, r, mail.GetStatusForError(unwrappedErr), err, unwrappedErr)
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
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 404 {object} errors.JSONError "invalid url address"
// @Router   /folder/{slug} [get]
func (del *delivery) GetFolderMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	folderSlug, ok := vars["slug"]
	if !ok {
		handleMailErr(w, r, mail.ErrInvalidURL)
		return
	}

	folder, err := del.uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		handleMailErr(w, r, fmt.Errorf("GetFolderInfo usecase error: %w", err))
		return
	}

	messages, err := del.uc.GetFolderMessages(userID, folderSlug)
	if err != nil {
		handleMailErr(w, r, fmt.Errorf("GetFolderMessages usecase error: %w", err))
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.FolderResponse{
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
// @Failure 404 {object} errors.JSONError "folder not found"
// @Router   /folders/ [get]
func (del *delivery) GetFolders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}

	folders, err := del.uc.GetFolders(userID)
	if err != nil {
		handleMailErr(w, r, fmt.Errorf("GetFolders usecase error: %w", err))
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.FoldersResponse{
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
// @Success  200 {object} models.FolderResponse "success get list of folder messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 404 {object} errors.JSONError "invalid url address"
// @Router   /folder/{slug} [get]
func (del *delivery) GetMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	messageID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		handleMailErr(w, r, mail.ErrInvalidURL)
		return
	}

	message, err := del.uc.GetMessage(userID, messageID)
	if err != nil {
		handleMailErr(w, r, fmt.Errorf("GetMessage usecase error: %w", err))
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

func (del *delivery) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}

	form := models.FormMessage{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		handleMailErr(w, r, mail.ErrInvalidMessageForm)
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		handleMailErr(w, r, mail.ErrInvalidMessageForm)
		return
	}

	validEmails, invalidEmails := del.uc.ValidateRecipients(form.Recipients)
	form.Recipients = validEmails

	message, err := del.uc.SendMessage(userID, form)
	if err != nil {
		handleMailErr(w, r, fmt.Errorf("SendMessage usecase error: %w", err))
		return
	}

	if len(invalidEmails) != 0 {
		err = del.uc.SendFailedSendingMessage(message.FromUser.Email, invalidEmails)

		if err != nil {
			handleMailErr(w, r, fmt.Errorf("SendFailedSendingMessage usecase error: %w", err))
			return
		}
	}

	pkg.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

func (del *delivery) ReadMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	messageID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		handleMailErr(w, r, mail.ErrInvalidURL)
		return
	}

	message, err := del.uc.MarkMessageAsSeen(userID, messageID)
	if err != nil {
		handleMailErr(w, r, fmt.Errorf("MarkMessageAsSeen usecase error: %w", err))
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

func (del *delivery) UnreadMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	messageID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		handleMailErr(w, r, mail.ErrInvalidURL)
		return
	}

	message, err := del.uc.MarkMessageAsUnseen(userID, messageID)
	if err != nil {
		handleMailErr(w, r, fmt.Errorf("MarkMessageAsUnseen usecase error: %w", err))
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}
