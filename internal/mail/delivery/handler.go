package delivery

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
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

func handleMailErr(w http.ResponseWriter, err error) {
	pkg.HandleError(w, mail.Errors[err], err)
}

// GetInboxMessages godoc
// @Summary      GetInboxMessages
// @Description  List of incoming messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} []models.IncomingMessage "success get list of incoming messages"
// @Failure 400 {object} error "failed to get user"
// @Failure 400 {object} error "failed to get inbox messages"
// @Failure 401 {object} error "failed auth"
// @Failure 401 {object} error "failed get session"
// @Router   /inbox/ [get]
func (del *delivery) GetInboxMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)

	if !ok {
		handleMailErr(w, mail.ErrFailedGetUser)
		return
	}

	messages, err := del.uc.GetIncomingMessages(userID)

	if err != nil {
		handleMailErr(w, mail.ErrFailedGetInboxMessages)
		return
	}

	pkg.SendJSON(w, http.StatusOK, messages)
}

// GetOutboxMessages godoc
// @Summary      GetOutboxMessages
// @Description  List of outgoing messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} []models.OutgoingMessage "success get list of outgoing messages"
// @Failure 400 {object} error "failed to get user"
// @Failure 400 {object} error "failed to get outbox messages"
// @Failure 401 {object} error "failed auth"
// @Failure 401 {object} error "failed get session"
// @Router   /outbox/ [get]
func (del *delivery) GetOutboxMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, mail.ErrFailedGetUser)
		return
	}

	messages, err := del.uc.GetOutgoingMessages(userID)
	if err != nil {
		handleMailErr(w, mail.ErrFailedGetOutboxMessages)
		return
	}

	pkg.SendJSON(w, http.StatusOK, messages)
}

// GetFolderMessages godoc
// @Summary      GetFolderMessages
// @Description  List of outgoing messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "FolderID"
// @Success  200 {object} []models.IncomingMessage "success get list of outgoing messages"
// @Failure 400 {object} error "failed to get user"
// @Failure 400 {object} error "failed to get folder messages"
// @Failure 401 {object} error "failed auth"
// @Failure 401 {object} error "failed get session"
// @Failure 404 {object} error "invalid url address"
// @Router   /folder/{id} [get]
func (del *delivery) GetFolderMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, mail.ErrFailedGetUser)
		return
	}
	vars := mux.Vars(r)
	folderID, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		handleMailErr(w, mail.ErrInvalidURL)
		return
	}

	messages, err := del.uc.GetFolderMessages(userID, folderID)

	if err != nil {
		handleMailErr(w, mail.ErrFailedGetFolderMessages)
		return
	}

	pkg.SendJSON(w, http.StatusOK, messages)
}

// GetFolders godoc
// @Summary      GetFolders
// @Description  List of outgoing messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} []models.Folder "success get list of outgoing messages"
// @Failure 400 {object} error "failed to get user"
// @Failure 401 {object} error "failed auth"
// @Failure 401 {object} error "failed get session"
// @Router   /folders/ [get]
func (del *delivery) GetFolders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(pkg.ContextUser).(uint64)
	if !ok {
		handleMailErr(w, mail.ErrFailedGetUser)
		return
	}

	folders := del.uc.GetFolders(userID)
	pkg.SendJSON(w, http.StatusOK, folders)
}
