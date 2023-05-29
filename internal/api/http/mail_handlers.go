package http

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/api/ws"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/crypto"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgHttp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/http"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/validation"
	pkgZip "github.com/go-park-mail-ru/2023_1_Seekers/pkg/zip"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	pkgErrors "github.com/pkg/errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
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
	DownloadAttach(w http.ResponseWriter, r *http.Request)
	DownloadAllAttaches(w http.ResponseWriter, r *http.Request)
	GetAttach(w http.ResponseWriter, r *http.Request)
	PreviewAttach(w http.ResponseWriter, r *http.Request)
	WSMessageHandler(w http.ResponseWriter, r *http.Request)
	DeleteDraftAttach(w http.ResponseWriter, r *http.Request)
	GetAttachB64(w http.ResponseWriter, r *http.Request)
	CreateAnonymousEmail(w http.ResponseWriter, r *http.Request)
	GetAnonymousEmails(w http.ResponseWriter, r *http.Request)
	DeleteAnonymousEmail(w http.ResponseWriter, r *http.Request)
	GetAnonymousMessages(w http.ResponseWriter, r *http.Request)
}

type mailHandlers struct {
	cfg *config.Config
	uc  mail.UseCaseI
	hub *ws.Hub
}

func NewMailHandlers(c *config.Config, uc mail.UseCaseI, h *ws.Hub) MailHandlersI {
	return &mailHandlers{
		cfg: c,
		uc:  uc,
		hub: h,
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
// @Success  200 {object} models.MessagesResponse "success get filtered messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /messages/search [get]
func (h *mailHandlers) SearchMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	fromUser := r.URL.Query().Get(h.cfg.Routes.RouteSearchQueryFromUser)
	toUser := r.URL.Query().Get(h.cfg.Routes.RouteSearchQueryToUser)
	filterText := r.URL.Query().Get(h.cfg.Routes.RouteSearchQueryFilter)
	folder := r.URL.Query().Get(h.cfg.Routes.RouteSearchQueryFolder)

	messages, err := h.uc.SearchMessages(userID, fromUser, toUser, folder, filterText)
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
// @Success  200 {object} []models.UserInfo "success get recipients"
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed read request body"))
		return
	}

	form := models.FormMessage{}
	if err := easyjson.Unmarshal(body, &form); err != nil {
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
		failedMessage, err := h.uc.SendFailedSendingMessage(message.FromUser.Email, invalidEmails)
		if err != nil {
			pkgHttp.HandleError(w, r, err)
			return
		}

		h.hub.SendNotifications(failedMessage)
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})

	h.hub.SendNotifications(message)
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed read request body"))
		return
	}

	form := models.FormMessage{}
	if err := easyjson.Unmarshal(body, &form); err != nil {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed read request body"))
		return
	}

	form := models.FormFolder{}
	if err := easyjson.Unmarshal(body, &form); err != nil {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed read request body"))
		return
	}

	form := models.FormFolder{}
	if err := easyjson.Unmarshal(body, &form); err != nil {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed read request body"))
		return
	}

	form := models.FormEditMessage{}
	if err := easyjson.Unmarshal(body, &form); err != nil {
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

	message, err := h.uc.EditDraft(userID, messageID, form)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessageResponse{
		Message: *message,
	})
}

// DownloadAttach godoc
// @Summary      DownloadAttach
// @Description  download attach, get attachID, check relation for attach and user
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessageResponse "success download attach"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "attach not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /attach/{id} [get]
func (h *mailHandlers) DownloadAttach(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	attachID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkgHttp.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	attach, err := h.uc.GetAttach(attachID, userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", attach.Type)
	w.Header().Set("Content-Disposition", "attachment; filename="+attach.FileName)

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(attach.FileData)
	if err != nil {
		pkgHttp.HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}

// GetAttachB64 godoc
// @Summary      GetAttachB64
// @Description  get attach in base64
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessageResponse "success download attach"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "attach not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /attach/{id} [get]
func (h *mailHandlers) GetAttachB64(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	attachID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkgHttp.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	attach, err := h.uc.GetAttach(attachID, userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	attachB64 := models.Attachment{FileName: attach.FileName, FileData: base64.StdEncoding.EncodeToString(attach.FileData)}
	w.Header().Set("Content-Type", common.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)

	pkgHttp.SendJSON(w, r, http.StatusOK, attachB64)
}

// DownloadAllAttaches godoc
// @Summary      DownloadAllAttaches
// @Description  download all attaches as zip, get messageID
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessageResponse "success download attach"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "attach not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/{id}/attaches [get]
func (h *mailHandlers) DownloadAllAttaches(w http.ResponseWriter, r *http.Request) {
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

	msg, err := h.uc.GetMessage(userID, messageID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	var zipArch bytes.Buffer

	zipWriter := zip.NewWriter(&zipArch)
	if len(msg.Attachments) == 0 {
		pkgHttp.HandleError(w, r, errors.ErrMessageNotFound)
		return
	} else {
		for _, a := range msg.Attachments {
			attach, err := h.uc.GetAttach(a.AttachID, userID)
			if err != nil {
				pkgHttp.HandleError(w, r, err)
				return
			}
			if err := pkgZip.Append2Zip(attach.FileName, attach.FileData, zipWriter); err != nil {
				pkgHttp.HandleError(w, r, err)
				return
			}
		}
	}
	zipWriter.Close()

	//fmt.Println("AFTER", zipArch.Bytes())

	// Use layout string for time format.
	const layout = "01-02-2006"
	// Place now in the string.
	t := time.Now()
	archiveName := "mailbx-archive-" + t.Format(layout) + ".zip"

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+archiveName)

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(zipArch.Bytes())
	if err != nil {
		pkgHttp.HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}

// GetAttach godoc
// @Summary      GetAttach
// @Description  get attach with attachID, gets access key, then validate
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessageResponse "success get attach"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "attach not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /external/attach/{id} [get]
func (h *mailHandlers) GetAttach(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	attachID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkgHttp.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	accessKey := r.URL.Query().Get(h.cfg.Routes.QueryAccessKey)

	userID, err := crypto.DecryptAccessToken(accessKey)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	attach, err := h.uc.GetAttach(attachID, userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", attach.Type)

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(attach.FileData)
	if err != nil {
		pkgHttp.HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}

// DeleteDraftAttach godoc
// @Summary      DeleteDraftAttach
// @Description  delete draft attach with attachID
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessageResponse "success get attach"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "attach not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /message/attach/{id} [delete]
func (h *mailHandlers) DeleteDraftAttach(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	attachID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkgHttp.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	attach, err := h.uc.GetAttach(attachID, userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", attach.Type)
	w.Header().Set("Content-Disposition", "attachment; filename="+attach.FileName)

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(attach.FileData)
	if err != nil {
		pkgHttp.HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}

// PreviewAttach godoc
// @Summary      PreviewAttach
// @Description  preview with attachID, returns html page
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessageResponse "success preview"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "attach not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /attach/{id}/preview [get]
func (h *mailHandlers) PreviewAttach(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	vars := mux.Vars(r)
	attachIDStr := vars["id"]
	attachID, err := strconv.ParseUint(attachIDStr, 10, 64)
	if err != nil {
		pkgHttp.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	attach, err := h.uc.GetAttachInfo(attachID, userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}
	tplFile := common.GetTplFile(attach.FileName)
	tpl, err := os.ReadFile(h.cfg.Api.MailTplDir + tplFile)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed to preview - read tpl file"))
		return
	}

	t, err := template.New("MailBx").Parse(string(tpl))
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed to preview - parse tpl file"))
		return
	}

	accessKey, err := crypto.EncryptAccessToken(userID)
	if err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed to preview - encrypt token"))
		return
	}

	data := struct {
		FileName     string
		FilePath     string
		FileDownload string
	}{
		FileName:     attach.FileName,
		FilePath:     fmt.Sprintf("%s/api/v1/external/attach/%s?accessKey=%s", h.cfg.Api.Host, attachIDStr, accessKey),
		FileDownload: fmt.Sprintf("%s/api/v1/attach/%s", h.cfg.Api.Host, attachIDStr),
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed to preview - execute tpl file"))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(buf.Bytes()); err != nil {
		pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed to preview - write body"))
		return
	}
}

func (h *mailHandlers) WSMessageHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get(h.cfg.Routes.RouteWsQueryEmail)
	ws.ServeWs(w, r, email, h.hub)

	w.Header().Set("Upgrade", "websocket")
}

// CreateAnonymousEmail godoc
// @Summary      CreateAnonymousEmail
// @Description  creating anonymous email
// @Tags     	 anonymous
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.AnonymousEmailResponse "success create anon email"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "max count anonymous emails is 5"
// @Failure 500 {object} errors.JSONError "error while generating fake email"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /anonymous/create [post]
func (h *mailHandlers) CreateAnonymousEmail(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	fakeEmail, err := h.uc.CreateAnonymousEmail(userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.AnonymousEmailResponse{Email: fakeEmail})
}

// GetAnonymousEmails godoc
// @Summary      GetAnonymousEmails
// @Description  get anonymous email by user
// @Tags     	 anonymous
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.AnonymousEmailsResponse "success get anon emails"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /anonymous [get]
func (h *mailHandlers) GetAnonymousEmails(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	fakeEmails, err := h.uc.GetAnonymousEmails(userID)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.AnonymousEmailsResponse{
		Emails: fakeEmails,
		Count:  len(fakeEmails),
	})
}

// DeleteAnonymousEmail godoc
// @Summary      DeleteAnonymousEmail
// @Description  delete anonymous email
// @Tags     	 anonymous
// @Accept	 application/json
// @Produce  application/json
// @Success  200 "success delete anonymous email"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "your anonymous email not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /anonymous [delete]
func (h *mailHandlers) DeleteAnonymousEmail(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	fakeEmail := r.URL.Query().Get(h.cfg.Routes.QueryAnonymousEmail)

	err := h.uc.DeleteAnonymousEmail(userID, fakeEmail)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetAnonymousMessages godoc
// @Summary      GetAnonymousMessages
// @Description  List of messages depends on anonymous email
// @Tags     	 anonymous
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.MessagesResponse "success get list of messages depends on anonymous email"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 404 {object} errors.JSONError "your anonymous email not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /anonymous/messages [get]
func (h *mailHandlers) GetAnonymousMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(common.ContextUser).(uint64)
	if !ok {
		pkgHttp.HandleError(w, r, errors.ErrFailedGetUser)
		return
	}

	fakeEmail := r.URL.Query().Get(h.cfg.Routes.QueryAnonymousEmail)

	messages, err := h.uc.GetMessagesByFakeEmail(userID, fakeEmail)
	if err != nil {
		pkgHttp.HandleError(w, r, err)
		return
	}

	pkgHttp.SendJSON(w, r, http.StatusOK, models.MessagesResponse{
		Messages: messages,
	})
}
