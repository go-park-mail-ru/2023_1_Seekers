package delivery

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	pkgErrors "github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

// GetInboxMessages godoc
// @Summary      GetInboxMessages
// @Description  List of incoming messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Success  200 {object} []models.IncomingMessage "success get list of incoming messages"
// @Failure 400 {object} error "a get request was expected"
// @Failure 400 {object} error "failed to get user"
// @Failure 400 {object} error "failed to get inbox messages"
// @Failure 401 {object} error "failed auth"
// @Failure 401 {object} error "failed get session"
// @Router   /inbox/ [get]
func (del *delivery) GetInboxMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrHttpGetMethod], mail.ErrHttpGetMethod.Error(), errors.New(r.Method+" request received"))
		log.Error(methodErr)
		pkg.SendError(w, methodErr)
		return
	}

	userID, ok := r.Context().Value(config.ContextUser).(uint64)

	if !ok {
		methodErr := pkgErrors.New(mail.MailErrors[mail.ErrFailedGetUser], mail.ErrFailedGetUser)
		log.Error(methodErr)
		pkg.SendError(w, methodErr)
		return
	}

	messages, err := del.uc.GetIncomingMessages(userID)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrFailedGetInboxMessages], mail.ErrFailedGetInboxMessages.Error(), err)
		log.Error(mailErr)
		pkg.SendError(w, mailErr)
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
// @Failure 400 {object} error "a get request was expected"
// @Failure 400 {object} error "failed to get user"
// @Failure 400 {object} error "failed to get outbox messages"
// @Failure 401 {object} error "failed auth"
// @Failure 401 {object} error "failed get session"
// @Router   /outbox/ [get]
func (del *delivery) GetOutboxMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrHttpGetMethod], mail.ErrHttpGetMethod.Error(), errors.New(r.Method+" request received"))
		log.Error(methodErr)
		pkg.SendError(w, methodErr)
		return
	}

	userID, ok := r.Context().Value(config.ContextUser).(uint64)

	if !ok {
		methodErr := pkgErrors.New(mail.MailErrors[mail.ErrFailedGetUser], mail.ErrFailedGetUser)
		log.Error(methodErr)
		pkg.SendError(w, methodErr)
		return
	}

	messages, err := del.uc.GetOutgoingMessages(userID)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrFailedGetOutboxMessages], mail.ErrFailedGetOutboxMessages.Error(), err)
		log.Error(mailErr)
		pkg.SendError(w, mailErr)
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
// @Failure 400 {object} error "a get request was expected"
// @Failure 400 {object} error "failed to get user"
// @Failure 400 {object} error "failed to get folder messages"
// @Failure 401 {object} error "failed auth"
// @Failure 401 {object} error "failed get session"
// @Failure 404 {object} error "invalid url address"
// @Router   /folder/{id} [get]
func (del *delivery) GetFolderMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrHttpGetMethod], mail.ErrHttpGetMethod.Error(), errors.New(r.Method+" request received"))
		log.Error(methodErr)
		pkg.SendError(w, methodErr)
		return
	}

	userID, ok := r.Context().Value(config.ContextUser).(uint64)
	if !ok {
		methodErr := pkgErrors.New(mail.MailErrors[mail.ErrFailedGetUser], mail.ErrFailedGetUser)
		log.Error(methodErr)
		pkg.SendError(w, methodErr)
		return
	}
	vars := mux.Vars(r)
	folderID, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrInvalidURL], mail.ErrInvalidURL.Error(), err)
		log.Error(mailErr)
		pkg.SendError(w, mailErr)
		return
	}

	messages, err := del.uc.GetFolderMessages(userID, folderID)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrFailedGetFolderMessages], mail.ErrFailedGetFolderMessages.Error(), err)
		log.Error(mailErr)
		pkg.SendError(w, mailErr)
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
// @Failure 400 {object} error "a get request was expected"
// @Failure 400 {object} error "failed to get user"
// @Failure 401 {object} error "failed auth"
// @Failure 401 {object} error "failed get session"
// @Router   /folders/ [get]
func (del *delivery) GetFolders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if r.Method != http.MethodGet {
		methodErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrHttpGetMethod], mail.ErrHttpGetMethod.Error(), errors.New(r.Method+" request received"))
		log.Error(methodErr)
		pkg.SendError(w, methodErr)
		return
	}

	userID, ok := r.Context().Value(config.ContextUser).(uint64)

	if !ok {
		methodErr := pkgErrors.New(mail.MailErrors[mail.ErrFailedGetUser], mail.ErrFailedGetUser)
		log.Error(methodErr)
		pkg.SendError(w, methodErr)
		return
	}

	folders := del.uc.GetFolders(userID)
	pkg.SendJSON(w, http.StatusOK, folders)
}
