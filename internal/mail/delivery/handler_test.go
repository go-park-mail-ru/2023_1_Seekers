package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	mockMailUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type inputCase struct {
	userID      uint64
	folderSlug  string
	messageID   uint64
	messageForm models.FormMessage
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
	mailH := New(mailUC)

	for _, test := range tests {
		r := httptest.NewRequest("GET", "/api/folder/", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"slug": test.input.folderSlug,
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().GetFolderInfo(test.input.userID, test.input.folderSlug).Return(&models.Folder{}, nil)
		mailUC.EXPECT().GetFolderMessages(test.input.userID, test.input.folderSlug).Return([]models.MessageInfo{}, nil)

		mailH.GetFolderMessages(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_GetFolders(t *testing.T) {
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
	mailH := New(mailUC)

	for _, test := range tests {
		r := httptest.NewRequest("GET", "/api/folders", bytes.NewReader([]byte{}))
		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().GetFolders(test.input.userID).Return([]models.Folder{}, nil)

		mailH.GetFolders(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_GetMessage(t *testing.T) {
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
	mailH := New(mailUC)

	for _, test := range tests {
		r := httptest.NewRequest("GET", "/api/message", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.messageID, 10),
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().GetMessage(test.input.userID, test.input.messageID).Return(&models.MessageInfo{}, nil)

		mailH.GetMessage(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_SendMessage(t *testing.T) {
	var tests = []testCases{
		{
			name: "one recipient",
			input: inputCase{
				userID: 1,
				messageForm: models.FormMessage{
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
	mailH := New(mailUC)

	for _, test := range tests {
		body, err := json.Marshal(test.input.messageForm)
		if err != nil {
			t.Fatalf("error while marshaling to json: %v", err)
		}

		r := httptest.NewRequest("POST", "/api/message/send", bytes.NewReader(body))
		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().ValidateRecipients(test.input.messageForm.Recipients).Return(test.input.messageForm.Recipients, []string{})
		mailUC.EXPECT().SendMessage(test.input.userID, test.input.messageForm).Return(&models.MessageInfo{}, nil)

		mailH.SendMessage(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_ReadMessage(t *testing.T) {
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
	mailH := New(mailUC)

	for _, test := range tests {
		r := httptest.NewRequest("POST", "/api/v1/message/read", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.messageID, 10),
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().MarkMessageAsSeen(test.input.userID, test.input.messageID).Return(&models.MessageInfo{}, nil)

		mailH.ReadMessage(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

func TestDelivery_UnreadMessage(t *testing.T) {
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
	mailH := New(mailUC)

	for _, test := range tests {
		r := httptest.NewRequest("POST", "/api/v1/message/unread", bytes.NewReader([]byte{}))
		vars := map[string]string{
			"id": strconv.FormatUint(test.input.messageID, 10),
		}
		r = mux.SetURLVars(r, vars)
		r = r.WithContext(context.WithValue(r.Context(), pkg.ContextUser, test.input.userID))
		w := httptest.NewRecorder()

		mailUC.EXPECT().MarkMessageAsUnseen(test.input.userID, test.input.messageID).Return(&models.MessageInfo{}, nil)

		mailH.UnreadMessage(w, r)

		if w.Code != test.output.status {
			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.status, w.Code)
		}
	}
}

//func TestDelivery_GetOutboxMessages(t *testing.T) {
//	var tests = []testCases{
//		{
//			name: "there are outgoing messages user1",
//			input: inputCase{
//				userCr: "user1",
//			},
//			output: outputCase{
//				err: http.StatusOK,
//			},
//		},
//		{
//			name: "there are outgoing messages user2",
//			input: inputCase{
//				userCr: "user2",
//			},
//			output: outputCase{
//				err: http.StatusOK,
//			},
//		},
//		{
//			name: "no outgoing messages user4",
//			input: inputCase{
//				userCr: "user4",
//			},
//			output: outputCase{
//				err: http.StatusOK,
//			},
//		},
//	}
//
//	t.Parallel()
//
//	for _, test := range tests {
//		r := httptest.NewRequest("GET", "/api/outbox", bytes.NewReader([]byte{}))
//		w := httptest.NewRecorder()
//		middleware, mailH := prepare(t, r, test.name, test.input.userCr)
//
//		middleware.CheckAuth(mailH.GetOutboxMessages)(w, r)
//
//		if w.Code != test.output.err {
//			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.err, w.Code)
//		}
//	}
//}

//func TestDelivery_GetFolderMessages(t *testing.T) {
//	var tests = []testCases{
//		{
//			name: "there are messages in the folder",
//			input: inputCase{
//				userCr:   "user2",
//				folderID: 7,
//			},
//			output: outputCase{
//				err: http.StatusOK,
//			},
//		},
//		{
//			name: "there are no messages in the folder",
//			input: inputCase{
//				userCr:   "user3",
//				folderID: 8,
//			},
//			output: outputCase{
//				err: http.StatusOK,
//			},
//		},
//		{
//			name: "the folder does not exist",
//			input: inputCase{
//				userCr:   "user1",
//				folderID: 7,
//			},
//			output: outputCase{
//				err: http.StatusBadRequest,
//			},
//		},
//	}
//
//	t.Parallel()
//
//	for _, test := range tests {
//		r := httptest.NewRequest("GET", "/api/folder/", bytes.NewReader([]byte{}))
//		r = mux.SetURLVars(r, map[string]string{"id": strconv.FormatUint(test.input.folderID, 10)})
//		w := httptest.NewRecorder()
//		middleware, mailH := prepare(t, r, test.name, test.input.userCr)
//
//		middleware.CheckAuth(mailH.GetFolderMessages)(w, r)
//
//		if w.Code != test.output.err {
//			t.Errorf("[TEST] %s: Expected err %d, got %d ", test.name, test.output.err, w.Code)
//		}
//	}
//}
