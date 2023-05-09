package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgHttp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/http"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	pkgErrors "github.com/pkg/errors"
	"net/http"
	"strconv"
)

type MailHandlersI interface {
	GetFolderMessages(w http.ResponseWriter, r *http.Request)
	SearchMessages(w http.ResponseWriter, r *http.Request)
	SearchRecipients(w http.ResponseWriter, r *http.Request)
	GetFolders(w http.ResponseWriter, r *http.Request)
	GetMessage(w http.ResponseWriter, r *http.Request)
	DeleteMessage(w http.ResponseWriter, r *http.Request)
	SendMessage(w http.ResponseWriter, r *http.Request)
	SaveDraft(w http.ResponseWriter, r *http.Request)
	ReadMessage(w http.ResponseWriter, r *http.Request)
	UnreadMessage(w http.ResponseWriter, r *http.Request)
	CreateFolder(w http.ResponseWriter, r *http.Request)
	DeleteFolder(w http.ResponseWriter, r *http.Request)
	EditFolder(w http.ResponseWriter, r *http.Request)
	MoveToFolder(w http.ResponseWriter, r *http.Request)
	EditDraft(w http.ResponseWriter, r *http.Request)
}

type mailHandlers struct {
	cfg *config.Config
	uc  mail.UseCaseI
}

func NewMailHandlers(c *config.Config, uc mail.UseCaseI) MailHandlersI {
	return &mailHandlers{
		cfg: c,
		uc:  uc,
	}
}

// GetFolderMessages godoc
// @Summary      GetFolderMessages
// @Description  List of folder messages
// @Tags     	 folders
// @Accept	 application/json
// @Produce  application/json
// @Param slug path string true "FolderSlug"
// @Success  200 {object} models.FolderResponse "success get list of folder messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid url address"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /folder/{slug} [get]
func (h *mailHandlers) GetFolderMessages(w http.ResponseWriter, r *http.Request) {
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

	folder, err := h.uc.GetFolderInfo(userID, folderSlug)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	messages, err := h.uc.GetFolderMessages(userID, folderSlug)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.FolderResponse{
		Folder:   *folder,
		Messages: messages,
	})
}

// SearchMessages godoc
// @Summary      SearchMessages
// @Description  list of filtered messages
// @Tags     	 folders
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} MessagesResponse "success get filtered messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /messages/search [get]
func (h *mailHandlers) SearchMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	filter := r.URL.Query().Get("filter")
	folder := r.URL.Query().Get("folder")

	messages, err := h.uc.SearchMessages(userID, "", "", folder, filter)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessagesResponse{
		Messages: messages,
	})
}

// SearchRecipients godoc
// @Summary      SearchRecipients
// @Description  list recipients for user
// @Tags     	 recipients
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} []UserInfo "success get recipients"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /recipients/search [get]
func (h *mailHandlers) SearchRecipients(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	usersInfo, err := h.uc.SearchRecipients(userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.Recipients{
		Users: usersInfo,
	})
}

// GetFolders godoc
// @Summary      GetFolders
// @Description  List of user folders
// @Tags     	 folders
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.FoldersResponse "success get list of outgoing messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /folders [get]
func (h *mailHandlers) GetFolders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	var folders []models.Folder

	var err error
	isCustom, _ := strconv.ParseBool(r.URL.Query().Get(h.cfg.Routes.RouteGetFoldersIsCustom))
	if isCustom {
		folders, err = h.uc.GetCustomFolders(userID)
		if len(folders) == 0 || folders == nil {
			folders = make([]models.Folder, 0)
		}
	} else {
		folders, err = h.uc.GetFolders(userID)
	}

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
func (h *mailHandlers) GetMessage(w http.ResponseWriter, r *http.Request) {
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

	message, err := h.uc.GetMessage(userID, messageID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

// DeleteMessage godoc
// @Summary      DeleteMessage
// @Description  delete message for user (moving to trash or full delete - depends of folder)
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Success  200 "success delete message"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid url address"
// @Failure 404 {object} errors.JSONError "message not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/{id} [delete]
func (h *mailHandlers) DeleteMessage(w http.ResponseWriter, r *http.Request) {
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

	fromFolder := r.URL.Query().Get(h.cfg.Routes.RouteQueryFromFolderSlug)

	err = h.uc.DeleteMessage(userID, messageID, fromFolder)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SendMessage godoc
// @Summary      SendMessage
// @Description  send message
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
func (h *mailHandlers) SendMessage(w http.ResponseWriter, r *http.Request) {
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

	validEmails, invalidEmails := h.uc.ValidateRecipients(form.Recipients)
	form.Recipients = validEmails

	message, err := h.uc.SendMessage(userID, form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	if len(invalidEmails) != 0 {
		err = h.uc.SendFailedSendingMessage(message.FromUser.Email, invalidEmails)

		if err != nil {
			pkgHttp.HandleError(w, r, err)
			return
		}
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

// SaveDraft godoc
// @Summary      SaveDraft
// @Description  save draft message
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessageResponse "success save draft message"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "some emails are invalid"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 404 {object} errors.JSONError "message not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/save [post]
func (h *mailHandlers) SaveDraft(w http.ResponseWriter, r *http.Request) {
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

	for _, email := range form.Recipients {
		if err := validation.ValidateEmail(email); err != nil {
			pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrSomeEmailsAreInvalid, "validate recipients"))
			return
		}
	}

	message, err := h.uc.SaveDraft(userID, form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
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
func (h *mailHandlers) ReadMessage(w http.ResponseWriter, r *http.Request) {
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

	fromFolder := r.URL.Query().Get(h.cfg.Routes.RouteQueryFromFolderSlug)

	message, err := h.uc.MarkMessageAsSeen(userID, messageID, fromFolder)
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
func (h *mailHandlers) UnreadMessage(w http.ResponseWriter, r *http.Request) {
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

	fromFolder := r.URL.Query().Get(h.cfg.Routes.RouteQueryFromFolderSlug)

	message, err := h.uc.MarkMessageAsUnseen(userID, messageID, fromFolder)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

// CreateFolder godoc
// @Summary      CreateFolder
// @Description  creating folder
// @Tags     	 folders
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.FolderResponse "success create folder"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid folder name"
// @Failure 400 {object} errors.JSONError "folder already exists"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /folder/create [post]
func (h *mailHandlers) CreateFolder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	form := models.FormFolder{}
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

	folder, err := h.uc.CreateFolder(userID, form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.FolderResponse{
		Folder: *folder,
	})
}

// DeleteFolder godoc
// @Summary      DeleteFolder
// @Description  delete folder
// @Tags     	 folders
// @Accept	 application/json
// @Produce  application/json
// @Param slug path string true "FolderSlug"
// @Success  200 "success delete folder"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid url address"
// @Failure 400 {object} errors.JSONError "can't delete default folder"
// @Failure 400 {object} errors.JSONError "message not found"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /folder/{slug} [delete]
func (h *mailHandlers) DeleteFolder(w http.ResponseWriter, r *http.Request) {
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

	err := h.uc.DeleteFolder(userID, folderSlug)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// EditFolder godoc
// @Summary      EditFolder
// @Description  edit folder name
// @Tags     	 folders
// @Accept	 application/json
// @Produce  application/json
// @Param slug path string true "FolderSlug"
// @Success  200 {object} models.FolderResponse "success edit folder"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid url address"
// @Failure 400 {object} errors.JSONError "can't edit default folder"
// @Failure 400 {object} errors.JSONError "folder already exists"
// @Failure 400 {object} errors.JSONError "invalid folder name"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /folder/{slug} [put]
func (h *mailHandlers) EditFolder(w http.ResponseWriter, r *http.Request) {
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

	form := models.FormFolder{}
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

	folder, err := h.uc.EditFolder(userID, folderSlug, form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.FolderResponse{
		Folder: *folder,
	})
}

// MoveToFolder godoc
// @Summary      MoveToFolder
// @Description  move message to folder
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Success  200 "success changed message's folder"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "invalid url address"
// @Failure 400 {object} errors.JSONError "can't move message to same folder"
// @Failure 400 {object} errors.JSONError "can't move message to draft folder"
// @Failure 400 {object} errors.JSONError "can't move message from draft folder"
// @Failure 404 {object} errors.JSONError "message not found"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/{id}/move [put]
func (h *mailHandlers) MoveToFolder(w http.ResponseWriter, r *http.Request) {
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

	fromFolder := r.URL.Query().Get(h.cfg.Routes.RouteQueryFromFolderSlug)
	toFolder := r.URL.Query().Get(h.cfg.Routes.RouteMoveToFolderQueryToFolderSlug)

	err = h.uc.MoveMessageToFolder(userID, messageID, fromFolder, toFolder)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// EditDraft godoc
// @Summary      EditDraft
// @Description  edit draft message
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessageResponse "success edit draft message"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "some emails are invalid"
// @Failure 403 {object} errors.JSONError "invalid form"
// @Failure 404 {object} errors.JSONError "folder not found"
// @Failure 404 {object} errors.JSONError "message not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/{id} [put]
func (h *mailHandlers) EditDraft(w http.ResponseWriter, r *http.Request) {
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

	validEmails, invalidEmails := h.uc.ValidateRecipients(form.Recipients)
	if len(invalidEmails) != 0 {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrSomeEmailsAreInvalid, "validate recipients"))
		return
	}

	form.Recipients = validEmails
	message, err := h.uc.EditDraft(userID, messageID, form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}
