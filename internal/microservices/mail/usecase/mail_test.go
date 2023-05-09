package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	mockMailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/repository/mocks"
	mockUserUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkg "github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/golang/mock/gomock"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type inputCase struct {
	userID      uint64
	folderSlug  string
	messageID   uint64
	messageForm models.FormMessage
}

type outputCase struct {
	content any
	err     error
}

type mockCase struct {
	content any
	err     any
}

type testCases struct {
	name     string
	input    inputCase
	output   outputCase
	fromMock mockCase
}

func createConfig() *config.Config {
	cfg := new(config.Config)
	cfg.Logger.LogsTimeFormat = "2006-01-02_15:04:05_MST"

	return cfg
}

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestUseCase_GetFolders(t *testing.T) {
	cfg := createConfig()
	var mockResponse []models.Folder
	generateFakeData(&mockResponse)

	var tests = []testCases{
		{
			name: "standard test",
			input: inputCase{
				userID: 1,
			},
			output: outputCase{
				content: mockResponse,
				err:     nil,
			},
			fromMock: mockCase{
				content: mockResponse,
				err:     nil,
			},
		},
		{
			name: "internal error",
			input: inputCase{
				userID: 1,
			},
			output: outputCase{
				content: []models.Folder{},
				err:     errors.ErrInternal,
			},
			fromMock: mockCase{
				content: []models.Folder{},
				err:     errors.ErrInternal,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	for _, test := range tests {
		mailRepo.EXPECT().SelectFoldersByUser(test.input.userID).Return(test.fromMock.content, test.fromMock.err)
		response, err := mailH.GetFolders(test.input.userID)
		causeErr := pkgErr.Cause(err)

		if causeErr != test.output.err {
			t.Errorf("[TEST] %s: expected err \"%v\", got \"%v\"", test.name, test.output.err, causeErr)
		} else {
			require.Equal(t, test.output.content, response)
		}
	}
}

func TestUseCase_GetFolderInfo(t *testing.T) {
	cfg := createConfig()

	var mockResponse *models.Folder
	var emptyResponse *models.Folder
	generateFakeData(&mockResponse)

	var tests = []testCases{
		{
			name: "folder is exists",
			input: inputCase{
				userID:     1,
				folderSlug: "inbox",
			},
			output: outputCase{
				content: mockResponse,
				err:     nil,
			},
			fromMock: mockCase{
				content: mockResponse,
				err:     nil,
			},
		},
		{
			name: "internal error",
			input: inputCase{
				userID: 1,
			},
			output: outputCase{
				content: emptyResponse,
				err:     errors.ErrInternal,
			},
			fromMock: mockCase{
				content: nil,
				err:     errors.ErrInternal,
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	for _, test := range tests {
		mailRepo.EXPECT().SelectFolderByUserNFolderSlug(test.input.userID, test.input.folderSlug).Return(test.fromMock.content, test.fromMock.err)
		response, err := mailH.GetFolderInfo(test.input.userID, test.input.folderSlug)
		causeErr := pkgErr.Cause(err)

		if causeErr != test.output.err {
			t.Errorf("[TEST] %s: expected err \"%v\", got \"%v\"", test.name, test.output.err, causeErr)
		} else {
			require.Equal(t, test.output.content, response)
		}
	}
}

func TestUseCase_GetFolderMessages(t *testing.T) {
	cfg := createConfig()
	userID := uint64(1)
	folderSlug := "inbox"

	mockFolderResponse := &models.Folder{
		FolderID:       1,
		UserID:         userID,
		LocalName:      "inbox",
		Name:           "Входящие",
		MessagesUnseen: 1,
		MessagesCount:  1,
	}
	mockFolderMessagesResponse := []models.MessageInfo{{
		MessageID:        1,
		Recipients:       nil,
		Title:            "test",
		CreatedAt:        "2023-01-29",
		Text:             "test text",
		ReplyToMessageID: nil,
		Seen:             false,
		Favorite:         false,
		Deleted:          false,
	}}
	mockUserResponse := []models.UserInfo{
		{
			UserID:    2,
			FirstName: "max",
			LastName:  "vlasov",
			Email:     "max03@mailbx.ru",
		},
		{
			UserID:    userID,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     "valera03@mailbx.ru",
		},
	}
	mockRecipientsResponse := []uint64{userID}

	output := outputCase{
		content: []models.MessageInfo{{
			MessageID:  mockFolderMessagesResponse[0].MessageID,
			FromUser:   mockUserResponse[0],
			Recipients: append([]models.UserInfo{}, mockUserResponse[1]),
			Title:      mockFolderMessagesResponse[0].Title,
			CreatedAt:  mockFolderMessagesResponse[0].CreatedAt,
			Text:       mockFolderMessagesResponse[0].Text,
		}},
		err: nil,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, folderSlug).Return(mockFolderResponse, nil)
	mailRepo.EXPECT().SelectFolderMessagesByUserNFolderID(userID, mockFolderResponse.FolderID, false).Return(mockFolderMessagesResponse, nil)
	userUC.EXPECT().GetInfo(mockFolderMessagesResponse[0].FromUser.UserID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(mockFolderMessagesResponse[0].MessageID, mockFolderMessagesResponse[0].FromUser.UserID).
		Return(mockRecipientsResponse, nil)
	userUC.EXPECT().GetInfo(mockRecipientsResponse[0]).Return(&mockUserResponse[1], nil)

	response, err := mailH.GetFolderMessages(userID, folderSlug)
	causeErr := pkgErr.Cause(err)

	if causeErr != output.err {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", output.err, causeErr)
	} else {
		require.Equal(t, output.content, response)
	}
}

func TestUseCase_CreateDefaultFolders(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	output := []models.Folder{
		{
			UserID:    userID,
			LocalName: "inbox",
			Name:      "Входящие",
		},
		{
			UserID:    userID,
			LocalName: "outbox",
			Name:      "Исходящие",
		},
		{
			UserID:    userID,
			LocalName: "trash",
			Name:      "Корзина",
		},
		{
			UserID:    userID,
			LocalName: "drafts",
			Name:      "Черновики",
		},
		{
			UserID:    userID,
			LocalName: "spam",
			Name:      "Спам",
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	for i := range output {
		mailRepo.EXPECT().InsertFolder(&output[i]).Return(uint64(i+1), nil)
	}
	mailRepo.EXPECT().SelectFoldersByUser(userID).Return(output, nil)

	response, err := mailH.CreateDefaultFolders(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, output, response)
	}
}

func TestUseCase_CreateFolder(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var fakeForm models.FormFolder
	var fakeFolders []models.Folder
	generateFakeData(&fakeForm)
	generateFakeData(&fakeFolders)
	fakeFolders[0].LocalName = "1"

	newFolder := &models.Folder{
		UserID:    userID,
		LocalName: "2",
		Name:      fakeForm.Name,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderName(userID, fakeForm.Name).Return(nil, errors.ErrFolderNotFound)
	mailRepo.EXPECT().SelectFoldersByUser(userID).Return(fakeFolders, nil)
	mailRepo.EXPECT().InsertFolder(newFolder).Return(uint64(0), nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, "2").Return(newFolder, nil)

	response, err := mailH.CreateFolder(userID, fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, newFolder, response)
	}
}

func TestUseCase_DeleteFolder(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var fakeFolder *models.Folder
	generateFakeData(&fakeFolder)
	fakeFolder.LocalName = "1"

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, fakeFolder.LocalName).Return(fakeFolder, nil)
	mailRepo.EXPECT().SelectFolderMessagesByUserNFolderID(userID, fakeFolder.FolderID, false).Return([]models.MessageInfo{}, nil)
	mailRepo.EXPECT().DeleteFolder(fakeFolder.FolderID).Return(nil)

	err := mailH.DeleteFolder(userID, fakeFolder.LocalName)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestUseCase_EditFolder(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var fakeFolder *models.Folder
	generateFakeData(&fakeFolder)
	fakeForm := models.FormFolder{Name: fakeFolder.Name}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, fakeFolder.LocalName).Return(fakeFolder, nil).AnyTimes()
	mailRepo.EXPECT().SelectFolderByUserNFolderName(userID, fakeForm.Name).Return(nil, errors.ErrFolderNotFound)
	mailRepo.EXPECT().UpdateFolder(*fakeFolder).Return(nil)

	response, err := mailH.EditFolder(userID, fakeFolder.LocalName, fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, *fakeFolder, *response)
	}
}

func TestUseCase_GetMessage(t *testing.T) {
	cfg := createConfig()
	userID := uint64(1)
	messageID := uint64(1)

	mockMessageResponse := &models.MessageInfo{
		FromUser:   models.UserInfo{UserID: 2},
		MessageID:  messageID,
		Recipients: nil,
		Title:      "test",
		CreatedAt:  "2023-01-29",
		Text:       "test text",
	}
	mockUserResponse := []models.UserInfo{
		{
			UserID:    2,
			FirstName: "max",
			LastName:  "vlasov",
			Email:     "max03@mailbx.ru",
		},
		{
			UserID:    userID,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     "valera03@mailbx.ru",
		},
	}
	mockRecipientsResponse := []uint64{userID}

	output := outputCase{
		content: &models.MessageInfo{
			MessageID:  messageID,
			FromUser:   mockUserResponse[0],
			Recipients: append([]models.UserInfo{}, mockUserResponse[1]),
			Title:      mockMessageResponse.Title,
			CreatedAt:  mockMessageResponse.CreatedAt,
			Text:       mockMessageResponse.Text,
		},
		err: nil,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageID).Return(mockMessageResponse, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(messageID, mockUserResponse[0].UserID).
		Return(mockRecipientsResponse, nil)
	userUC.EXPECT().GetInfo(mockRecipientsResponse[0]).Return(&mockUserResponse[1], nil)

	response, err := mailH.GetMessage(userID, messageID)
	causeErr := pkgErr.Cause(err)

	if causeErr != output.err {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", output.err, causeErr)
	} else {
		require.Equal(t, output.content, response)
	}
}

func TestUseCase_DeleteMessage(t *testing.T) {
	cfg := createConfig()
	userID := uint64(1)
	messageID := uint64(1)

	mockMessageResponse := &models.MessageInfo{
		FromUser:   models.UserInfo{UserID: 2},
		MessageID:  messageID,
		Recipients: nil,
		Title:      "test",
		CreatedAt:  "2023-01-29",
		Text:       "test text",
	}
	var fakeFolder *models.Folder
	generateFakeData(&fakeFolder)

	var tests = []testCases{
		{
			name: "remove from trash",
			input: inputCase{
				folderSlug: "trash",
			},
		},
		{
			name: "remove from draft",
			input: inputCase{
				folderSlug: "drafts",
			},
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	for _, test := range tests {
		fakeFolder.LocalName = test.input.folderSlug
		mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageID).Return(mockMessageResponse, nil).AnyTimes()
		mailRepo.EXPECT().SelectFolderByUserNMessage(userID, messageID).Return(fakeFolder, nil).AnyTimes()
		//mailRepo.EXPECT().DeleteBox(userID, messageID).Return(nil).AnyTimes()
		mailRepo.EXPECT().DeleteMessageFromMessages(messageID).Return(nil).AnyTimes()

		// TODO
		err := mailH.DeleteMessage(userID, messageID, "TODO:folderslug")
		causeErr := pkgErr.Cause(err)

		if causeErr != test.output.err {
			t.Errorf("[TEST] %s: expected err \"%v\", got \"%v\"", test.name, test.output.err, causeErr)
		}
	}
}

func TestUseCase_ValidateRecipients(t *testing.T) {
	cfg := createConfig()

	var users [3]models.UserInfo
	generateFakeData(&users)
	emails := []string{users[0].Email, users[1].Email, users[2].Email}
	outputValid := emails[:2]
	outputInvalid := emails[2:]

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	userUC.EXPECT().GetInfoByEmail(emails[0]).Return(&users[0], nil)
	userUC.EXPECT().GetInfoByEmail(emails[1]).Return(&users[1], nil)
	userUC.EXPECT().GetInfoByEmail(emails[2]).Return(nil, errors.ErrInternal)

	responseValid, responseInvalid := mailH.ValidateRecipients(emails)

	require.Equal(t, outputValid, responseValid)
	require.Equal(t, outputInvalid, responseInvalid)
}

func TestUseCase_SaveDraft(t *testing.T) {
	cfg := createConfig()
	userID := uint64(1)

	var formMessage models.FormMessage
	generateFakeData(&formMessage)
	formMessage.Recipients = []string{"max@mailbox.ru"}
	mockUserResponse := []models.UserInfo{
		{
			UserID:    userID,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     "valera03@mailbox.ru",
		},
		{
			UserID:    2,
			FirstName: "max",
			LastName:  "vlasov",
			Email:     "max03@mailbox.ru",
		},
	}
	mockFolderResponse := []models.Folder{
		{
			FolderID:  1,
			UserID:    userID,
			LocalName: "outbox",
			Name:      "Исходящие",
		},
		{
			FolderID:  2,
			UserID:    mockUserResponse[1].UserID,
			LocalName: "inbox",
			Name:      "Входящие",
		},
	}
	user2folder := []models.User2Folder{
		{
			UserID:   userID,
			FolderID: mockFolderResponse[0].FolderID,
		},
		{
			UserID:   mockUserResponse[1].UserID,
			FolderID: mockFolderResponse[1].FolderID,
		},
	}
	newMessage := models.MessageInfo{
		Title:            formMessage.Title,
		CreatedAt:        pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
		Text:             formMessage.Text,
		ReplyToMessageID: formMessage.ReplyToMessageID,
		IsDraft:          true,
	}
	messageSelected := newMessage
	messageSelected.FromUser = mockUserResponse[0]
	messageSelected.ReplyToMessageID = nil

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, "drafts").Return(&mockFolderResponse[0], nil)
	userUC.EXPECT().GetInfoByEmail(formMessage.Recipients[0]).Return(&mockUserResponse[1], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[1].UserID, "inbox").Return(&mockFolderResponse[1], nil)
	mailRepo.EXPECT().InsertMessage(userID, &newMessage, user2folder).Return(nil).SetArg(1, messageSelected)
	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageSelected.MessageID).Return(&messageSelected, nil)
	userUC.EXPECT().GetInfo(userID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(messageSelected.MessageID, userID).Return([]uint64{mockUserResponse[1].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[0], nil)

	response, err := mailH.SaveDraft(userID, formMessage)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, messageSelected, *response)
	}
}

func TestUseCase_EditDraft(t *testing.T) {
	cfg := createConfig()
	userID := uint64(1)
	messageID := uint64(1)

	mockMessageResponse := &models.MessageInfo{
		FromUser:   models.UserInfo{UserID: 2},
		MessageID:  messageID,
		Recipients: nil,
		Title:      "test",
		CreatedAt:  pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
		Text:       "test text",
		IsDraft:    true,
	}
	mockUserResponse := []models.UserInfo{
		{
			UserID:    2,
			FirstName: "max",
			LastName:  "vlasov",
			Email:     "max03@mailbox.ru",
		},
		{
			UserID:    userID,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     "valera03@mailbox.ru",
		},
		{
			UserID:    3,
			FirstName: "oleg",
			LastName:  "kotkov",
			Email:     "oleg@mailbox.ru",
		},
	}
	mockRecipientsResponse := []uint64{userID}
	var mockFolderResponse [2]models.Folder
	generateFakeData(&mockFolderResponse)

	formMessage := models.FormMessage{
		Recipients:       []string{mockUserResponse[2].Email},
		Title:            "test",
		Text:             "test text",
		ReplyToMessageID: nil,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageID).Return(mockMessageResponse, nil).AnyTimes()
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil).AnyTimes()
	mailRepo.EXPECT().SelectRecipientsByMessage(messageID, mockUserResponse[0].UserID).
		Return(mockRecipientsResponse, nil).AnyTimes()
	userUC.EXPECT().GetInfo(mockRecipientsResponse[0]).Return(&mockUserResponse[1], nil).AnyTimes()
	userUC.EXPECT().GetInfoByEmail(mockUserResponse[2].Email).Return(&mockUserResponse[2], nil)
	userUC.EXPECT().GetInfoByEmail(mockUserResponse[1].Email).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[2].UserID, "inbox").Return(&mockFolderResponse[0], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[0].UserID, "inbox").Return(&mockFolderResponse[1], nil)

	recipients := map[string]string{
		mockUserResponse[0].Email: "del",
		mockUserResponse[2].Email: "add",
	}

	var toInsert []models.User2Folder
	var toDelete []models.User2Folder

	for _, value := range recipients {
		switch value {
		case "add":
			toInsert = append(toInsert, models.User2Folder{
				UserID:   mockUserResponse[2].UserID,
				FolderID: mockFolderResponse[0].FolderID,
			})
		case "del":
			toDelete = append(toDelete, models.User2Folder{
				UserID:   mockUserResponse[0].UserID,
				FolderID: mockFolderResponse[1].FolderID,
			})
		}
	}
	mailRepo.EXPECT().UpdateMessage(mockMessageResponse, toInsert, toDelete).Return(nil)

	response, err := mailH.EditDraft(userID, messageID, formMessage)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, *mockMessageResponse, *response)
	}
}

func TestUseCase_SendMessage(t *testing.T) {
	cfg := createConfig()
	userID := uint64(1)

	var formMessage models.FormMessage
	generateFakeData(&formMessage)
	formMessage.Recipients = []string{"max@mailbx.ru"}

	mockUserResponse := []models.UserInfo{
		{
			UserID:    userID,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     "valera03@mailbx.ru",
		},
		{
			UserID:    2,
			FirstName: "max",
			LastName:  "vlasov",
			Email:     "max03@mailbx.ru",
		},
	}
	mockFolderResponse := []models.Folder{
		{
			FolderID:  1,
			UserID:    userID,
			LocalName: "outbox",
			Name:      "Исходящие",
		},
		{
			FolderID:  2,
			UserID:    mockUserResponse[1].UserID,
			LocalName: "inbox",
			Name:      "Входящие",
		},
	}

	user2folder := []models.User2Folder{
		{
			UserID:   userID,
			FolderID: mockFolderResponse[0].FolderID,
		},
		{
			UserID:   mockUserResponse[1].UserID,
			FolderID: mockFolderResponse[1].FolderID,
		},
	}
	newMessage := models.MessageInfo{
		Title:            formMessage.Title,
		CreatedAt:        pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
		Text:             formMessage.Text,
		ReplyToMessageID: formMessage.ReplyToMessageID,
	}
	messageSelected := newMessage
	messageSelected.FromUser = mockUserResponse[0]
	messageSelected.ReplyToMessageID = nil

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, "outbox").Return(&mockFolderResponse[0], nil)
	userUC.EXPECT().GetInfoByEmail(formMessage.Recipients[0]).Return(&mockUserResponse[1], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[1].UserID, "inbox").Return(&mockFolderResponse[1], nil)
	mailRepo.EXPECT().InsertMessage(userID, &newMessage, user2folder).Return(nil).SetArg(1, messageSelected)
	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageSelected.MessageID).Return(&messageSelected, nil)
	userUC.EXPECT().GetInfo(userID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(messageSelected.MessageID, userID).Return([]uint64{mockUserResponse[1].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[0], nil)

	response, err := mailH.SendMessage(userID, formMessage)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, messageSelected, *response)
	}
}

func TestUseCase_SendFailedSendingMessage(t *testing.T) {
	cfg := createConfig()

	supportEmail := "support@mailbx.ru"
	userEmail := "valera03@mailbx.ru"
	mockUserResponse := []models.UserInfo{
		{
			UserID:    1,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     userEmail,
		},
		{
			UserID:    2,
			FirstName: "Suppot",
			LastName:  "Supportov",
			Email:     supportEmail,
		},
	}
	invalidEmails := []string{"123123123@mailbx.ru"}
	formMessage := models.FormMessage{
		Recipients: []string{userEmail},
		Title:      "Ваше сообщение не доставлено",
		Text: "Это письмо создано автоматически сервером mailbx.ru, отвечать на него не нужно.\n\n" +
			"К сожалению, Ваше письмо не может быть доставлено одному или нескольким получателям:\n" +
			strings.Join(invalidEmails, "\n") + "\n\nРекомендуем Вам проверить корректность указания адресов получателей.",
		ReplyToMessageID: nil,
	}
	mockFolderResponse := []models.Folder{
		{
			FolderID:  1,
			UserID:    mockUserResponse[0].UserID,
			LocalName: "inbox",
			Name:      "Входящие",
		},
		{
			FolderID:  2,
			UserID:    mockUserResponse[1].UserID,
			LocalName: "outbox",
			Name:      "Исходящие",
		},
	}
	user2folder := []models.User2Folder{
		{
			UserID:   mockUserResponse[1].UserID,
			FolderID: mockFolderResponse[1].FolderID,
		},
		{
			UserID:   mockUserResponse[0].UserID,
			FolderID: mockFolderResponse[0].FolderID,
		},
	}
	newMessage := models.MessageInfo{
		Title:            formMessage.Title,
		CreatedAt:        pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
		Text:             formMessage.Text,
		ReplyToMessageID: formMessage.ReplyToMessageID,
	}
	messageSelected := newMessage
	messageSelected.FromUser = mockUserResponse[1]
	messageSelected.ReplyToMessageID = nil

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	userUC.EXPECT().GetInfoByEmail(userEmail).Return(&mockUserResponse[0], nil)
	userUC.EXPECT().GetInfoByEmail(supportEmail).Return(&mockUserResponse[1], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[1].UserID, "outbox").Return(&mockFolderResponse[1], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[0].UserID, "inbox").Return(&mockFolderResponse[0], nil)
	mailRepo.EXPECT().InsertMessage(mockUserResponse[1].UserID, &newMessage, user2folder).Return(nil).SetArg(1, messageSelected)
	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[1].UserID, messageSelected.MessageID).Return(&messageSelected, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(messageSelected.MessageID, mockUserResponse[1].UserID).Return([]uint64{mockUserResponse[0].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)

	err := mailH.SendFailedSendingMessage(userEmail, invalidEmails)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	}
}

func TestUseCase_SendWelcomeMessage(t *testing.T) {
	cfg := createConfig()

	supportEmail := "support@mailbx.ru"
	userEmail := "valera03@mailbx.ru"
	mockUserResponse := []models.UserInfo{
		{
			UserID:    1,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     userEmail,
		},
		{
			UserID:    2,
			FirstName: "Suppot",
			LastName:  "Supportov",
			Email:     supportEmail,
		},
	}
	formMessage := models.FormMessage{
		Recipients: []string{userEmail},
		Title:      "Добро пожаловать в почту Mailbx",
		Text: "Это письмо создано автоматически сервером Mailbx.ru, отвечать на него не нужно.\n" +
			"Поздравляем Вас с присоединением к нашей почте. Уверены, что вы останетесь довольны ее использованием!",
		ReplyToMessageID: nil,
	}
	mockFolderResponse := []models.Folder{
		{
			FolderID:  1,
			UserID:    mockUserResponse[0].UserID,
			LocalName: "inbox",
			Name:      "Входящие",
		},
		{
			FolderID:  2,
			UserID:    mockUserResponse[1].UserID,
			LocalName: "outbox",
			Name:      "Исходящие",
		},
	}
	user2folder := []models.User2Folder{
		{
			UserID:   mockUserResponse[1].UserID,
			FolderID: mockFolderResponse[1].FolderID,
		},
		{
			UserID:   mockUserResponse[0].UserID,
			FolderID: mockFolderResponse[0].FolderID,
		},
	}
	newMessage := models.MessageInfo{
		Title:            formMessage.Title,
		CreatedAt:        pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
		Text:             formMessage.Text,
		ReplyToMessageID: formMessage.ReplyToMessageID,
	}
	messageSelected := newMessage
	messageSelected.FromUser = mockUserResponse[1]
	messageSelected.ReplyToMessageID = nil

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	userUC.EXPECT().GetInfoByEmail(userEmail).Return(&mockUserResponse[0], nil)
	userUC.EXPECT().GetInfoByEmail(supportEmail).Return(&mockUserResponse[1], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[1].UserID, "outbox").Return(&mockFolderResponse[1], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[0].UserID, "inbox").Return(&mockFolderResponse[0], nil)
	mailRepo.EXPECT().InsertMessage(mockUserResponse[1].UserID, &newMessage, user2folder).Return(nil).SetArg(1, messageSelected)
	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[1].UserID, messageSelected.MessageID).Return(&messageSelected, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(messageSelected.MessageID, mockUserResponse[1].UserID).Return([]uint64{mockUserResponse[0].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)

	err := mailH.SendWelcomeMessage(userEmail)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	}
}

func TestUseCase_MarkMessageAsSeen(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	messageID := uint64(1)
	folderID := uint64(1)
	state := "seen"
	stateValue := true

	mockUserResponse := []models.UserInfo{
		{
			UserID:    userID,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     "valera03@mailbx.ru",
		},
		{
			UserID:    2,
			FirstName: "max",
			LastName:  "vlasov",
			Email:     "max03@mailbx.ru",
		},
	}
	mockMessageResponse := &models.MessageInfo{
		FromUser:   mockUserResponse[0],
		MessageID:  messageID,
		Recipients: nil,
		Title:      "test",
		CreatedAt:  "2023-01-29",
		Text:       "test text",
		Seen:       stateValue,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().UpdateMessageState(userID, messageID, folderID, state, stateValue).Return(nil)
	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[0].UserID, mockMessageResponse.MessageID).Return(mockMessageResponse, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(mockMessageResponse.MessageID, mockUserResponse[0].UserID).Return([]uint64{mockUserResponse[1].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)

	// TODO
	response, err := mailH.MarkMessageAsSeen(userID, messageID, "TODO FOLDER SLUG")
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, mockMessageResponse, response)
	}
}

func TestUseCase_MarkMessageAsUnseen(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	messageID := uint64(1)
	folderID := uint64(1)
	state := "seen"
	stateValue := false

	mockUserResponse := []models.UserInfo{
		{
			UserID:    userID,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     "valera03@mailbx.ru",
		},
		{
			UserID:    2,
			FirstName: "max",
			LastName:  "vlasov",
			Email:     "max03@mailbx.ru",
		},
	}
	mockMessageResponse := &models.MessageInfo{
		FromUser:   mockUserResponse[0],
		MessageID:  messageID,
		Recipients: nil,
		Title:      "test",
		CreatedAt:  "2023-01-29",
		Text:       "test text",
		Seen:       stateValue,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().UpdateMessageState(userID, messageID, folderID, state, stateValue).Return(nil)
	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[0].UserID, mockMessageResponse.MessageID).Return(mockMessageResponse, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(mockMessageResponse.MessageID, mockUserResponse[0].UserID).Return([]uint64{mockUserResponse[1].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)

	// TODO
	response, err := mailH.MarkMessageAsUnseen(userID, messageID, "TODO FOLDER SLUG")
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, mockMessageResponse, response)
	}
}

func TestUseCase_MoveMessageToFolder(t *testing.T) {
	cfg := createConfig()
	userID := uint64(1)
	var fakeMessage *models.MessageInfo
	var fakeFromFolder, fakeToFolder *models.Folder
	generateFakeData(&fakeMessage)
	generateFakeData(&fakeFromFolder)
	generateFakeData(&fakeToFolder)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, fakeMessage.MessageID).Return(fakeMessage, nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, fakeToFolder.LocalName).Return(fakeToFolder, nil)
	mailRepo.EXPECT().SelectFolderByUserNMessage(userID, fakeMessage.MessageID).Return(fakeFromFolder, nil)
	mailRepo.EXPECT().UpdateMessageFolder(userID, fakeMessage.MessageID, fakeFromFolder.FolderID, fakeToFolder.FolderID).Return(nil)

	// TODO
	err := mailH.MoveMessageToFolder(userID, fakeMessage.MessageID, "TODO FROM FOLDER", "TODO TO FOLDER")
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	}
}

func TestUseCase_GetCustomFolders(t *testing.T) {
	cfg := createConfig()
	userID := uint64(1)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC)

	mailRepo.EXPECT().SelectCustomFoldersByUser(userID, gomock.Any()).Return([]models.Folder{}, nil)

	response, err := mailH.GetCustomFolders(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, []models.Folder{}, response)
	}
}
