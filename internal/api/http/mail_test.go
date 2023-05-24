package http

import (
	"bytes"
	"context"
	mockMailUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type inputCase struct {
	userID      uint64
	folderSlug  string
	fromFolder  string
	toFolder    string
	messageID   uint64
	attachID    uint64
	messageForm models.FormMessage
	folderForm  models.FormFolder
}

type outputCase struct {
	status int
}

type testCases struct {
	name   string
	input  inputCase
	output outputCase
}

func TestDelivery_GetFolderMessages(t *testing.T) {
	cfg := createConfig()
	var tests = []testCases{
		{
			name: "default folder",
			input: inputCase{
				userID:     1,
				folderSlug: "inbox",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "custom folder",
			input: inputCase{
				userID:     1,
				folderSlug: "1",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodGet, "/api/folder/", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"slug": test.input.folderSlug,
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().GetFolderInfo(test.input.userID, test.input.folderSlug).Return(&models.Folder{}, nil)
		mailUC.EXPECT().GetFolderMessages(test.input.userID, test.input.folderSlug).Return([]models.MessageInfo{}, nil)

		mailH.GetFolderMessages(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_SearchMessages(t *testing.T) {
	cfg := createConfig()
	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID:     1,
				fromFolder: "inbox",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodGet, "/api/v2/messages/search", bytes.NewReader([]byte{}))
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		q := r.URL.Query()
		q.Set("filter", "123")
		q.Set("folder", test.input.fromFolder)
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()

		mailUC.EXPECT().SearchMessages(test.input.userID, "", "", test.input.fromFolder, "123").Return([]models.MessageInfo{}, nil)

		mailH.SearchMessages(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_SearchRecipients(t *testing.T) {
	cfg := createConfig()
	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID: 1,
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodGet, "/api/v2/recipients/search", bytes.NewReader([]byte{}))
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().SearchRecipients(test.input.userID).Return([]models.UserInfo{}, nil)

		mailH.SearchRecipients(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_GetFolders(t *testing.T) {
	cfg := createConfig()
	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID: 1,
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodGet, "/api/folders", bytes.NewReader([]byte{}))
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().GetFolders(test.input.userID).Return([]models.Folder{}, nil)

		mailH.GetFolders(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_GetMessage(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID:    1,
				messageID: 1,
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodGet, "/api/message", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.messageID, 10),
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().GetMessage(test.input.userID, test.input.messageID).Return(&models.MessageInfo{}, nil)

		mailH.GetMessage(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_DeleteMessage(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID:     1,
				messageID:  1,
				fromFolder: "outbox",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodDelete, "/api/message", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.messageID, 10),
			//"fromFolder": test.input.fromFolder,
		}
		q := r.URL.Query()
		q.Set("fromFolder", test.input.fromFolder)
		r.URL.RawQuery = q.Encode()
		r = mux.SetURLVars(r, vars)

		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().DeleteMessage(test.input.userID, test.input.messageID, test.input.fromFolder).Return(nil)

		mailH.DeleteMessage(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_SendMessage(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "one recipient",
			input: inputCase{
				userID: 1,
				messageForm: models.FormMessage{
					FromUser:         "max@mailbx.ru",
					Recipients:       []string{"valera@mailbx.ru"},
					Title:            "title test message",
					Text:             "text test message",
					ReplyToMessageID: nil,
				},
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "several recipient",
			input: inputCase{
				userID: 1,
				messageForm: models.FormMessage{
					FromUser:         "max@mailbx.ru",
					Recipients:       []string{"valera@mailbx.ru", "max@mailbx.ru"},
					Title:            "title test message",
					Text:             "text test message",
					ReplyToMessageID: nil,
				},
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		body, err := easyjson.Marshal(test.input.messageForm)
		if err != nil {
			t.Fatalf("error while marshaling to json: %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/api/message/send", bytes.NewReader(body))
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().ValidateRecipients(test.input.messageForm.Recipients).Return(test.input.messageForm.Recipients, []string{})
		mailUC.EXPECT().SendMessage(test.input.userID, test.input.messageForm).Return(&models.MessageInfo{}, nil)

		mailH.SendMessage(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_SaveDraft(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "one recipient",
			input: inputCase{
				userID: 1,
				messageForm: models.FormMessage{
					FromUser:         "max@mailbx.ru",
					Recipients:       []string{"valera@mailbox.ru"},
					Title:            "title test message",
					Text:             "text test message",
					ReplyToMessageID: nil,
				},
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
		{
			name: "several recipient",
			input: inputCase{
				userID: 1,
				messageForm: models.FormMessage{
					FromUser:         "test@mailbx.ru",
					Recipients:       []string{"valera@mailbox.ru", "max@mailbox.ru"},
					Title:            "title test message",
					Text:             "text test message",
					ReplyToMessageID: nil,
				},
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		body, err := easyjson.Marshal(test.input.messageForm)
		if err != nil {
			t.Fatalf("error while marshaling to json: %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/api/message/save", bytes.NewReader(body))
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().SaveDraft(test.input.userID, test.input.messageForm).Return(&models.MessageInfo{}, nil)

		mailH.SaveDraft(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_ReadMessage(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID:     1,
				messageID:  1,
				fromFolder: "outbox",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/message/read", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.messageID, 10),
		}

		q := r.URL.Query()
		q.Set("fromFolder", test.input.fromFolder)
		r.URL.RawQuery = q.Encode()
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().MarkMessageAsSeen(test.input.userID, test.input.messageID, test.input.fromFolder).Return(&models.MessageInfo{}, nil)

		mailH.ReadMessage(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_UnreadMessage(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID:     1,
				messageID:  1,
				fromFolder: "outbox",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/message/unread", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.messageID, 10),
		}

		q := r.URL.Query()
		q.Set("fromFolder", test.input.fromFolder)
		r.URL.RawQuery = q.Encode()
		r = mux.SetURLVars(r, vars)

		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().MarkMessageAsUnseen(test.input.userID, test.input.messageID, test.input.fromFolder).Return(&models.MessageInfo{}, nil)

		mailH.UnreadMessage(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_CreateFolder(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID:     1,
				folderForm: models.FormFolder{Name: "test"},
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		body, err := easyjson.Marshal(test.input.folderForm)
		if err != nil {
			t.Fatalf("error while marshaling to json: %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/folder/create", bytes.NewReader(body))
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().CreateFolder(test.input.userID, test.input.folderForm).Return(&models.Folder{}, nil)

		mailH.CreateFolder(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_DeleteFolder(t *testing.T) {
	cfg := createConfig()
	var tests = []testCases{
		{
			name: "custom folder",
			input: inputCase{
				userID:     1,
				folderSlug: "1",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodDelete, "/api/folder/", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"slug": test.input.folderSlug,
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().DeleteFolder(test.input.userID, test.input.folderSlug).Return(nil)
		mailH.DeleteFolder(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_EditFolder(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID:     1,
				folderSlug: "1",
				folderForm: models.FormFolder{Name: "test"},
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		body, err := easyjson.Marshal(test.input.folderForm)
		if err != nil {
			t.Fatalf("error while marshaling to json: %v", err)
		}

		r := httptest.NewRequest(http.MethodPut, "/folder/", bytes.NewReader(body))
		vars := map[string]string{
			"slug": test.input.folderSlug,
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().EditFolder(test.input.userID, test.input.folderSlug, test.input.folderForm).Return(&models.Folder{}, nil)

		mailH.EditFolder(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_MoveToFolder(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID:     1,
				messageID:  1,
				fromFolder: "inbox",
				toFolder:   "trash",
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodGet, "/api/message/move", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.messageID, 10),
		}
		r = mux.SetURLVars(r, vars)

		q := r.URL.Query()
		q.Add(cfg.Routes.RouteMoveToFolderQueryToFolderSlug, test.input.toFolder)
		q.Add(cfg.Routes.RouteQueryFromFolderSlug, test.input.fromFolder)
		r.URL.RawQuery = q.Encode()

		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().MoveMessageToFolder(test.input.userID, test.input.messageID, test.input.fromFolder, test.input.toFolder).Return(nil)

		mailH.MoveToFolder(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_EditDraft(t *testing.T) {
	cfg := createConfig()

	type editInputCase struct {
		userID      uint64
		folderSlug  string
		fromFolder  string
		toFolder    string
		messageID   uint64
		attachID    uint64
		messageForm models.FormEditMessage
		folderForm  models.FormFolder
	}

	type testCases struct {
		name   string
		input  editInputCase
		output outputCase
	}

	var tests = []testCases{
		{
			name: "standart test",
			input: editInputCase{
				userID:    1,
				messageID: 1,
				messageForm: models.FormEditMessage{
					FromUser:          "max@mailbx.ru",
					Recipients:        []string{"valera@mailbox.ru"},
					Title:             "title test message",
					Text:              "text test message",
					ReplyToMessageID:  nil,
					DeleteAttachments: nil,
					NewAttachments:    nil,
				},
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		body, err := easyjson.Marshal(test.input.messageForm)
		if err != nil {
			t.Fatalf("error while marshaling to json: %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/api/message/send", bytes.NewReader(body))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.messageID, 10),
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().ValidateRecipients(test.input.messageForm.Recipients).Return(test.input.messageForm.Recipients, []string{})
		mailUC.EXPECT().EditDraft(test.input.userID, test.input.messageID, test.input.messageForm).Return(&models.MessageInfo{}, nil)

		mailH.EditDraft(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_DownloadAttah(t *testing.T) {
	cfg := createConfig()

	var tests = []testCases{
		{
			name: "standart test",
			input: inputCase{
				userID:   1,
				attachID: 1,
			},
			output: outputCase{
				status: http.StatusOK,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/attach/", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.attachID, 10),
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().GetAttach(test.input.attachID, test.input.userID).Return(&models.AttachmentInfo{}, nil)

		mailH.DownloadAttach(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_DownloadAllAttaches(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	messageID := uint64(1)
	var msg *models.MessageInfo
	var attaches []models.AttachmentInfo
	generateFakeData(&msg)
	generateFakeData(&attaches)
	msg.Attachments = attaches[:1]

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/attach/", bytes.NewReader([]byte{}))
	vars := map[string]string{
		"id": strconv.FormatUint(messageID, 10),
	}
	r = mux.SetURLVars(r, vars)
	r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, userID))
	w := httptest.NewRecorder()

	mailUC.EXPECT().GetMessage(userID, messageID).Return(msg, nil)
	mailUC.EXPECT().GetAttach(msg.Attachments[0].AttachID, userID).Return(&msg.Attachments[0], nil)

	mailH.DownloadAllAttaches(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", http.StatusOK, w.Code)
	}
}

func TestDelivery_GetAttach(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var attach *models.AttachmentInfo
	generateFakeData(&attach)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/attach/", bytes.NewReader([]byte{}))
	vars := map[string]string{
		"id": strconv.FormatUint(attach.AttachID, 10),
	}
	r = mux.SetURLVars(r, vars)
	q := r.URL.Query()
	q.Set(cfg.Routes.QueryAccessKey, "Xi6uVrevawNwAmAcc0dp4xjzszoZfbo3wLy-MA==")
	r.URL.RawQuery = q.Encode()
	r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, userID))
	w := httptest.NewRecorder()

	mailH.GetAttach(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", http.StatusBadRequest, w.Code)
	}
}

func TestDelivery_DeleteDraftAttach(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var attach *models.AttachmentInfo
	generateFakeData(&attach)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/attach/", bytes.NewReader([]byte{}))
	vars := map[string]string{
		"id": strconv.FormatUint(attach.AttachID, 10),
	}
	r = mux.SetURLVars(r, vars)
	r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, userID))
	w := httptest.NewRecorder()

	mailUC.EXPECT().GetAttach(attach.AttachID, userID).Return(attach, nil)

	mailH.DeleteDraftAttach(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", http.StatusOK, w.Code)
	}
}

func TestDelivery_PreviewAttach(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var attach *models.AttachmentInfo
	generateFakeData(&attach)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailUC := mockMailUC.NewMockUseCaseI(ctrl)
	mailH := NewMailHandlers(cfg, mailUC)

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/attach/", bytes.NewReader([]byte{}))
	vars := map[string]string{
		"id": strconv.FormatUint(attach.AttachID, 10),
	}
	r = mux.SetURLVars(r, vars)
	r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, userID))
	w := httptest.NewRecorder()

	mailUC.EXPECT().GetAttachInfo(attach.AttachID, userID).Return(attach, nil)

	mailH.PreviewAttach(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("[TEST] simple: Expected err %d, got %d ", http.StatusInternalServerError, w.Code)
	}
}
