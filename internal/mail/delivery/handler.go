package delivery

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	pkgErrors "github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type DeliveryI interface {
	GetInboxMessages(w http.ResponseWriter, r *http.Request)
	GetOutboxMessages(w http.ResponseWriter, r *http.Request)
	GetFolderMessages(w http.ResponseWriter, r *http.Request)
}

type delivery struct {
	uc mail.UseCaseI
}

func New(uc mail.UseCaseI) DeliveryI {
	return &delivery{
		uc: uc,
	}
}

func (del *delivery) GetInboxMessages(w http.ResponseWriter, r *http.Request) {
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
	messages, err := del.uc.GetIncomingMessages(userID)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrFailedGetInboxMessages], mail.ErrFailedGetInboxMessages.Error(), err)
		log.Error(mailErr)
		pkg.SendError(w, mailErr)
		return
	}

	pkg.SendJSON(w, http.StatusOK, models.InboxMessages{
		Folders:  folders,
		Messages: messages,
	})
}

func (del *delivery) GetOutboxMessages(w http.ResponseWriter, r *http.Request) {
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
	messages, err := del.uc.GetOutgoingMessages(userID)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrFailedGetOutboxMessages], mail.ErrFailedGetOutboxMessages.Error(), err)
		log.Error(mailErr)
		pkg.SendError(w, mailErr)
		return
	}

	pkg.SendJSON(w, http.StatusOK, models.OutboxMessages{
		Folders:  folders,
		Messages: messages,
	})
}

func (del *delivery) GetFolderMessages(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)
	folderID, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrInvalidURL], mail.ErrInvalidURL.Error(), err)
		log.Error(mailErr)
		pkg.SendError(w, mailErr)
		return
	}

	folders := del.uc.GetFolders(userID)
	messages, err := del.uc.GetFolderMessages(userID, folderID)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrFailedGetFolderMessages], mail.ErrFailedGetFolderMessages.Error(), err)
		log.Error(mailErr)
		pkg.SendError(w, mailErr)
		return
	}

	pkg.SendJSON(w, http.StatusOK, models.InboxMessages{
		Folders:  folders,
		Messages: messages,
	})
}
