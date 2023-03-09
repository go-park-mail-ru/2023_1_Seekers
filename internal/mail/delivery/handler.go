package delivery

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	_ "github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
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
	pkg.HandleError(w, r, mail.Errors[err], err)
}

// GetInboxMessages godoc
// @Summary      GetInboxMessages
// @Description  List of incoming messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.InboxResponse "success get list of incoming messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "failed to get inbox messages"
// @Failure 401 {object} errors.JSONError "failed auth"
// @Failure 401 {object} errors.JSONError "failed get session"
// @Router   /inbox/ [get]
func (del *delivery) GetInboxMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)

	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}

	messages, err := del.uc.GetIncomingMessages(userID)

	if err != nil {
		handleMailErr(w, r, mail.ErrFailedGetInboxMessages)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.InboxResponse{
		Messages: messages,
	})
}

// GetOutboxMessages godoc
// @Summary      GetOutboxMessages
// @Description  List of outgoing messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} models.OutboxResponse "success get list of outgoing messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "failed to get outbox messages"
// @Failure 401 {object} errors.JSONError "failed auth"
// @Failure 401 {object} errors.JSONError "failed get session"
// @Router   /outbox/ [get]
func (del *delivery) GetOutboxMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}

	messages, err := del.uc.GetOutgoingMessages(userID)
	if err != nil {
		handleMailErr(w, r, mail.ErrFailedGetOutboxMessages)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.OutboxResponse{
		Messages: messages,
	})
}

// GetFolderMessages godoc
// @Summary      GetFolderMessages
// @Description  List of outgoing messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "FolderID"
// @Success  200 {object} models.FolderResponse "success get list of outgoing messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "failed to get folder messages"
// @Failure 401 {object} errors.JSONError "failed auth"
// @Failure 401 {object} errors.JSONError "failed get session"
// @Failure 404 {object} errors.JSONError "invalid url address"
// @Router   /folder/{id} [get]
func (del *delivery) GetFolderMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}
	vars := mux.Vars(r)
	folderID, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		handleMailErr(w, r, mail.ErrInvalidURL)
		return
	}

	folder, err := del.uc.GetFolderInfo(userID, folderID)

	if err != nil {
		handleMailErr(w, r, mail.ErrFailedGetFolderMessages)
		return
	}

	messages, err := del.uc.GetFolderMessages(userID, folderID)

	if err != nil {
		handleMailErr(w, r, mail.ErrFailedGetFolderMessages)
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
// @Failure 401 {object} errors.JSONError "failed auth"
// @Failure 401 {object} errors.JSONError "failed get session"
// @Router   /folders/ [get]
func (del *delivery) GetFolders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, r, mail.ErrFailedGetUser)
		return
	}

	folders := del.uc.GetFolders(userID)
	pkg.SendJSON(w, r, http.StatusOK, models.FoldersResponse{
		Folders: folders,
	})
}
