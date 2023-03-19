package delivery

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	_ "github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/gorilla/mux"
	"net/http"
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

// GetFolderMessages godoc
// @Summary      GetFolderMessages
// @Description  List of folder messages
// @Tags     	 messages
// @Accept	 application/json
// @Produce  application/json
// @Param slug path string true "FolderSlug"
// @Success  200 {object} models.FolderResponse "success get list of folder messages"
// @Failure 400 {object} errors.JSONError "failed to get user"
// @Failure 400 {object} errors.JSONError "failed to get folder"
// @Failure 400 {object} errors.JSONError "failed to get folder messages"
// @Failure 401 {object} errors.JSONError "failed auth"
// @Failure 401 {object} errors.JSONError "failed get session"
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
		handleMailErr(w, r, mail.ErrFailedGetFolder)
		return
	}

	messages, err := del.uc.GetFolderMessages(userID, folderSlug)

	if err != nil {
		fmt.Println(err)
		handleMailErr(w, r, mail.ErrFailedGetFolderMessages)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.FolderResponse{
		Folder:         *folder,
		Messages:       messages,
		MessagesUnseen: folder.MessagesUnseen,
		MessagesCount:  folder.MessagesCount,
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

	folders, err := del.uc.GetFolders(userID)
	if err != nil {
		handleMailErr(w, r, mail.ErrFailedGetFolders)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.FoldersResponse{
		Folders: folders,
	})
}
