package delivery

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/internal/mail"
	mailUC "github.com/go-park-mail-ru/2023_1_Seekers/app/internal/mail/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/models"
	pkg2 "github.com/go-park-mail-ru/2023_1_Seekers/app/pkg"
	pkgErrors "github.com/go-park-mail-ru/2023_1_Seekers/app/pkg/errors"
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
	uc mailUC.UseCaseI
}

func New(uc mailUC.UseCaseI) DeliveryI {
	return &delivery{
		uc: uc,
	}
}

func (del *delivery) GetInboxMessages(w http.ResponseWriter, r *http.Request) {
	var userID = uint64(2) // get from auth
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if r.Method != http.MethodGet {
		methodErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.HttpGetMethodError], mail.HttpGetMethodError.Error(), errors.New(r.Method+" request received"))
		log.Error(methodErr)
		pkg2.SendError(w, methodErr)

		return
	}

	folders := del.uc.GetFolders(userID)
	messages, err := del.uc.GetIncomingMessages(userID)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrFailedGetInboxMessages], mail.ErrFailedGetInboxMessages.Error(), err)
		log.Error(mailErr)
		pkg2.SendError(w, mailErr)

		return
	}

	pkg2.SendJSON(w, http.StatusOK, models.InboxMessages{
		Folders:  folders,
		Messages: messages,
	})
}

func (del *delivery) GetOutboxMessages(w http.ResponseWriter, r *http.Request) {
	var userID = uint64(2) // get from auth
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if r.Method != http.MethodGet {
		methodErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.HttpGetMethodError], mail.HttpGetMethodError.Error(), errors.New(r.Method+" request received"))
		log.Error(methodErr)
		pkg2.SendError(w, methodErr)

		return
	}

	folders := del.uc.GetFolders(userID)
	messages, err := del.uc.GetOutgoingMessages(userID)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrFailedGetOutboxMessages], mail.ErrFailedGetOutboxMessages.Error(), err)
		log.Error(mailErr)
		pkg2.SendError(w, mailErr)

		return
	}

	pkg2.SendJSON(w, http.StatusOK, models.OutboxMessages{
		Folders:  folders,
		Messages: messages,
	})
}

func (del *delivery) GetFolderMessages(w http.ResponseWriter, r *http.Request) {
	var userID = uint64(2) // get from auth
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if r.Method != http.MethodGet {
		methodErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.HttpGetMethodError], mail.HttpGetMethodError.Error(), errors.New(r.Method+" request received"))
		log.Error(methodErr)
		pkg2.SendError(w, methodErr)

		return
	}

	vars := mux.Vars(r)
	folderID, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.InvalidURL], mail.InvalidURL.Error(), err)
		log.Error(mailErr)
		pkg2.SendError(w, mailErr)

		return
	}

	folders := del.uc.GetFolders(userID)
	messages, err := del.uc.GetFolderMessages(userID, folderID)

	if err != nil {
		mailErr := pkgErrors.NewWrappedErr(mail.MailErrors[mail.ErrFailedGetFolderMessages], mail.ErrFailedGetFolderMessages.Error(), err)
		log.Error(mailErr)
		pkg2.SendError(w, mailErr)

		return
	}

	pkg2.SendJSON(w, http.StatusOK, models.InboxMessages{
		Folders:  folders,
		Messages: messages,
	})
}
