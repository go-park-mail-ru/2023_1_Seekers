package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	mockFileUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/usecase/mocks"
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
	cfg.S3.S3AttachBucket = "test"

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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

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
			MessageID:       mockFolderMessagesResponse[0].MessageID,
			FromUser:        mockUserResponse[0],
			Recipients:      append([]models.UserInfo{}, mockUserResponse[1]),
			Title:           mockFolderMessagesResponse[0].Title,
			CreatedAt:       mockFolderMessagesResponse[0].CreatedAt,
			Text:            mockFolderMessagesResponse[0].Text,
			AttachmentsSize: "0 Б",
		}},
		err: nil,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, folderSlug).Return(mockFolderResponse, nil)
	mailRepo.EXPECT().SelectFolderMessagesByUserNFolderID(userID, mockFolderResponse.FolderID, false).Return(mockFolderMessagesResponse, nil)
	mailRepo.EXPECT().GetMessageAttachments(mockFolderMessagesResponse[0].MessageID).Return(nil, nil)
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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

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
			MessageID:       messageID,
			FromUser:        mockUserResponse[0],
			Recipients:      append([]models.UserInfo{}, mockUserResponse[1]),
			Title:           mockMessageResponse.Title,
			CreatedAt:       mockMessageResponse.CreatedAt,
			Text:            mockMessageResponse.Text,
			AttachmentsSize: "0 Б",
		},
		err: nil,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageID).Return(mockMessageResponse, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(messageID, mockUserResponse[0].UserID).
		Return(mockRecipientsResponse, nil)
	mailRepo.EXPECT().GetMessageAttachments(messageID).Return(nil, nil)
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
	folderID := uint64(1)

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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailUCase := New(cfg, mailRepo, userUC, fileUC)

	for _, test := range tests {
		fakeFolder.LocalName = test.input.folderSlug
		mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageID).Return(mockMessageResponse, nil).AnyTimes()
		mailRepo.EXPECT().SelectFolderByUserNMessage(userID, messageID).Return(fakeFolder, nil).AnyTimes()
		mailRepo.EXPECT().DeleteMessageFromMessages(messageID).Return(nil).AnyTimes()
		mailRepo.EXPECT().CheckExistingBox(userID, messageID, folderID).Return(true, nil)
		mailRepo.EXPECT().DeleteBox(userID, messageID, folderID).Return(nil).AnyTimes()

		mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, test.input.folderSlug).Return(&models.Folder{
			FolderID:       folderID,
			UserID:         userID,
			LocalName:      test.input.folderSlug,
			Name:           test.input.folderSlug,
			MessagesUnseen: 0,
			MessagesCount:  0,
		}, nil)

		err := mailUCase.DeleteMessage(userID, messageID, test.input.folderSlug)
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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

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
	formMessage.Recipients = []string{"max@mailbx.ru"}
	formMessage.ReplyToMessageID = nil
	formMessage.Attachments = nil
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
		IsDraft:          true,
		Attachments:      nil,
		AttachmentsSize:  "0 Б",
	}

	messageSelected := newMessage
	//messageSelected.FromUser = mockUserResponse[0]
	messageSelected.ReplyToMessageID = nil
	messageSelected.Attachments = []models.AttachmentInfo{}
	messageSelected.AttachmentsSize = "0 Б"

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, "drafts").Return(&mockFolderResponse[0], nil)
	userUC.EXPECT().GetInfoByEmail(formMessage.Recipients[0]).Return(&mockUserResponse[1], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[1].UserID, "inbox").Return(&mockFolderResponse[1], nil)
	mailRepo.EXPECT().InsertMessage(userID, &messageSelected, user2folder).Return(nil)
	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageSelected.MessageID).Return(&messageSelected, nil)
	mailRepo.EXPECT().GetMessageAttachments(messageSelected.MessageID).Return(nil, nil)
	userUC.EXPECT().GetInfo(uint64(0)).Return(&mockUserResponse[0], nil).Times(1)
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
		FromUser:        models.UserInfo{UserID: 2},
		MessageID:       messageID,
		Recipients:      nil,
		Title:           "test",
		CreatedAt:       pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
		Text:            "test text",
		IsDraft:         true,
		AttachmentsSize: "0 Б",
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
		{
			UserID:    3,
			FirstName: "oleg",
			LastName:  "kotkov",
			Email:     "oleg@mailbx.ru",
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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageID).Return(mockMessageResponse, nil).AnyTimes()
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil).AnyTimes()
	mailRepo.EXPECT().GetMessageAttachments(messageID).Return(nil, nil).AnyTimes()
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
	formMessage.FromUser = "valera03@mailbx.ru"
	formMessage.Attachments = nil
	formMessage.ReplyToMessageID = nil

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
	mockUsers := []models.User{
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
		AttachmentsSize:  "0 Б",
		Attachments:      []models.AttachmentInfo{},
	}
	messageSelected := newMessage
	messageSelected.MessageID = uint64(1)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	userUC.EXPECT().GetByEmail(formMessage.FromUser).Return(&mockUsers[0], nil).Times(1)
	userUC.EXPECT().GetByEmail(formMessage.Recipients[0]).Return(&mockUsers[1], nil).Times(1)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, "outbox").Return(&mockFolderResponse[0], nil)
	userUC.EXPECT().GetInfoByEmail(formMessage.Recipients[0]).Return(&mockUserResponse[1], nil).AnyTimes()
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[1].UserID, "inbox").Return(&mockFolderResponse[1], nil)
	mailRepo.EXPECT().InsertMessage(userID, &newMessage, user2folder).Return(nil).SetArg(1, messageSelected)
	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, messageSelected.MessageID).Return(&messageSelected, nil)
	userUC.EXPECT().GetInfo(uint64(0)).Return(&mockUserResponse[0], nil).Times(1)
	mailRepo.EXPECT().SelectRecipientsByMessage(messageSelected.MessageID, userID).Return([]uint64{mockUserResponse[1].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().GetMessageAttachments(messageSelected.MessageID).Return([]models.AttachmentInfo{}, nil)

	response, err := mailH.SendMessage(userID, formMessage)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, messageSelected, *response)
	}
}

func TestUseCase_SendFailedSendingMessage(t *testing.T) {
	cfg := createConfig()

	supportEmail := "support@mailbx.ru"
	userEmail := "valera03@mailbx.ru"
	mockUserInfoResponse := []models.UserInfo{
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
	mockUserResponse := []models.User{
		{
			UserID:    mockUserInfoResponse[0].UserID,
			FirstName: mockUserInfoResponse[0].FirstName,
			LastName:  mockUserInfoResponse[0].LastName,
			Email:     mockUserInfoResponse[0].Email,
		},
		{
			UserID:    mockUserInfoResponse[1].UserID,
			FirstName: mockUserInfoResponse[1].FirstName,
			LastName:  mockUserInfoResponse[1].LastName,
			Email:     mockUserInfoResponse[1].Email,
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
		Attachments:      nil,
	}
	mockFolderResponse := []models.Folder{
		{
			FolderID:  1,
			UserID:    mockUserInfoResponse[0].UserID,
			LocalName: "inbox",
			Name:      "Входящие",
		},
		{
			FolderID:  2,
			UserID:    mockUserInfoResponse[1].UserID,
			LocalName: "outbox",
			Name:      "Исходящие",
		},
	}
	user2folder := []models.User2Folder{
		{
			UserID:   mockUserInfoResponse[1].UserID,
			FolderID: mockFolderResponse[1].FolderID,
		},
		{
			UserID:   mockUserInfoResponse[0].UserID,
			FolderID: mockFolderResponse[0].FolderID,
		},
	}
	newMessage := models.MessageInfo{
		Title:            formMessage.Title,
		CreatedAt:        pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
		Text:             formMessage.Text,
		ReplyToMessageID: formMessage.ReplyToMessageID,
		AttachmentsSize:  "0 Б",
		Attachments:      []models.AttachmentInfo{},
	}
	messageSelected := newMessage
	messageSelected.FromUser = mockUserInfoResponse[1]
	messageSelected.ReplyToMessageID = nil

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	userUC.EXPECT().GetInfoByEmail(supportEmail).Return(&mockUserInfoResponse[1], nil)
	userUC.EXPECT().GetByEmail(supportEmail).Return(&mockUserResponse[1], nil)
	userUC.EXPECT().GetByEmail(userEmail).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserInfoResponse[1].UserID, "outbox").Return(&mockFolderResponse[1], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserInfoResponse[0].UserID, "inbox").Return(&mockFolderResponse[0], nil)
	mailRepo.EXPECT().InsertMessage(mockUserInfoResponse[1].UserID, &newMessage, user2folder).Return(nil).SetArg(1, messageSelected)
	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserInfoResponse[1].UserID, messageSelected.MessageID).Return(&messageSelected, nil)
	userUC.EXPECT().GetInfo(mockUserInfoResponse[1].UserID).Return(&mockUserInfoResponse[1], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(messageSelected.MessageID, mockUserInfoResponse[1].UserID).Return([]uint64{mockUserInfoResponse[0].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserInfoResponse[0].UserID).Return(&mockUserInfoResponse[0], nil)
	mailRepo.EXPECT().GetMessageAttachments(messageSelected.MessageID).Return([]models.AttachmentInfo{}, nil)

	_, err := mailH.SendFailedSendingMessage(userEmail, invalidEmails)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestUseCase_SendWelcomeMessage(t *testing.T) {
	cfg := createConfig()

	supportEmail := "support@mailbx.ru"
	userEmail := "valera03@mailbx.ru"
	mockUserInfoResponse := []models.UserInfo{
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
	mockUserResponse := []models.User{
		{
			UserID:    mockUserInfoResponse[0].UserID,
			FirstName: mockUserInfoResponse[0].FirstName,
			LastName:  mockUserInfoResponse[0].LastName,
			Email:     mockUserInfoResponse[0].Email,
		},
		{
			UserID:    mockUserInfoResponse[1].UserID,
			FirstName: mockUserInfoResponse[1].FirstName,
			LastName:  mockUserInfoResponse[1].LastName,
			Email:     mockUserInfoResponse[1].Email,
		},
	}
	formMessage := models.FormMessage{
		Recipients: []string{userEmail},
		Title:      "Добро пожаловать в почту Mailbx",
		Text: "Это письмо создано автоматически сервером Mailbx.ru, отвечать на него не нужно.\n" +
			"Поздравляем Вас с присоединением к нашей почте. Уверены, что вы останетесь довольны ее использованием!",
		ReplyToMessageID: nil,
		Attachments:      nil,
	}
	mockFolderResponse := []models.Folder{
		{
			FolderID:  1,
			UserID:    mockUserInfoResponse[0].UserID,
			LocalName: "inbox",
			Name:      "Входящие",
		},
		{
			FolderID:  2,
			UserID:    mockUserInfoResponse[1].UserID,
			LocalName: "outbox",
			Name:      "Исходящие",
		},
	}
	user2folder := []models.User2Folder{
		{
			UserID:   mockUserInfoResponse[1].UserID,
			FolderID: mockFolderResponse[1].FolderID,
		},
		{
			UserID:   mockUserInfoResponse[0].UserID,
			FolderID: mockFolderResponse[0].FolderID,
		},
	}
	newMessage := models.MessageInfo{
		Title:            formMessage.Title,
		CreatedAt:        pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
		Text:             formMessage.Text,
		ReplyToMessageID: formMessage.ReplyToMessageID,
		AttachmentsSize:  "0 Б",
		Attachments:      []models.AttachmentInfo{},
	}
	messageSelected := newMessage
	messageSelected.FromUser = mockUserInfoResponse[1]
	messageSelected.ReplyToMessageID = nil

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	userUC.EXPECT().GetInfoByEmail(supportEmail).Return(&mockUserInfoResponse[1], nil)
	userUC.EXPECT().GetByEmail(supportEmail).Return(&mockUserResponse[1], nil)
	userUC.EXPECT().GetByEmail(userEmail).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserInfoResponse[1].UserID, "outbox").Return(&mockFolderResponse[1], nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserInfoResponse[0].UserID, "inbox").Return(&mockFolderResponse[0], nil)
	mailRepo.EXPECT().InsertMessage(mockUserInfoResponse[1].UserID, &newMessage, user2folder).Return(nil).SetArg(1, messageSelected)
	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserInfoResponse[1].UserID, messageSelected.MessageID).Return(&messageSelected, nil)
	userUC.EXPECT().GetInfo(mockUserInfoResponse[1].UserID).Return(&mockUserInfoResponse[1], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(messageSelected.MessageID, mockUserInfoResponse[1].UserID).Return([]uint64{mockUserInfoResponse[0].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserInfoResponse[0].UserID).Return(&mockUserInfoResponse[0], nil)
	mailRepo.EXPECT().GetMessageAttachments(messageSelected.MessageID).Return([]models.AttachmentInfo{}, nil)

	err := mailH.SendWelcomeMessage(userEmail)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
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
	mockFolderResponse := &models.Folder{
		FolderID:  folderID,
		UserID:    mockUserResponse[0].UserID,
		LocalName: "inbox",
		Name:      "Входящие",
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, mockFolderResponse.LocalName).Return(mockFolderResponse, nil)
	mailRepo.EXPECT().CheckExistingBox(userID, messageID, folderID).Return(true, nil)
	mailRepo.EXPECT().UpdateMessageState(userID, messageID, folderID, state, stateValue).Return(nil)
	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[0].UserID, mockMessageResponse.MessageID).Return(mockMessageResponse, nil).AnyTimes()
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(mockMessageResponse.MessageID, mockUserResponse[0].UserID).Return([]uint64{mockUserResponse[1].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)
	mailRepo.EXPECT().GetMessageAttachments(messageID).Return([]models.AttachmentInfo{}, nil).AnyTimes()

	response, err := mailH.MarkMessageAsSeen(userID, messageID, mockFolderResponse.LocalName)
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
	mockFolderResponse := &models.Folder{
		FolderID:  folderID,
		UserID:    mockUserResponse[0].UserID,
		LocalName: "inbox",
		Name:      "Входящие",
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, mockFolderResponse.LocalName).Return(mockFolderResponse, nil)
	mailRepo.EXPECT().CheckExistingBox(userID, messageID, folderID).Return(true, nil)
	mailRepo.EXPECT().UpdateMessageState(userID, messageID, folderID, state, stateValue).Return(nil)
	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[0].UserID, mockMessageResponse.MessageID).Return(mockMessageResponse, nil).AnyTimes()
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(mockMessageResponse.MessageID, mockUserResponse[0].UserID).Return([]uint64{mockUserResponse[1].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)
	mailRepo.EXPECT().GetMessageAttachments(messageID).Return([]models.AttachmentInfo{}, nil).AnyTimes()

	response, err := mailH.MarkMessageAsUnseen(userID, messageID, mockFolderResponse.LocalName)
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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SelectMessageByUserNMessage(userID, fakeMessage.MessageID).Return(fakeMessage, nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, fakeToFolder.LocalName).Return(fakeToFolder, nil)
	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(userID, fakeFromFolder.LocalName).Return(fakeFromFolder, nil)
	mailRepo.EXPECT().CheckExistingBox(userID, fakeMessage.MessageID, fakeFromFolder.FolderID).Return(true, nil)
	mailRepo.EXPECT().CheckExistingBox(userID, fakeMessage.MessageID, fakeToFolder.FolderID).Return(false, nil)
	mailRepo.EXPECT().UpdateMessageFolder(userID, fakeMessage.MessageID, fakeFromFolder.FolderID, fakeToFolder.FolderID).Return(nil)

	err := mailH.MoveMessageToFolder(userID, fakeMessage.MessageID, fakeFromFolder.LocalName, fakeToFolder.LocalName)
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
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SelectCustomFoldersByUser(userID, gomock.Any()).Return([]models.Folder{}, nil)

	response, err := mailH.GetCustomFolders(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, []models.Folder{}, response)
	}
}

func TestUseCase_SearchMessages(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var mockMessagesResponse []models.MessageInfo
	generateFakeData(&mockMessagesResponse)
	mockMessagesResponse = mockMessagesResponse[:1]
	var mockFromUserResponse *models.UserInfo
	var mockRecipientResponse []models.UserInfo
	generateFakeData(&mockFromUserResponse)
	generateFakeData(&mockRecipientResponse)
	mockRecipientResponse = mockRecipientResponse[:1]
	mockMessagesResponse[0].FromUser = *mockFromUserResponse
	mockMessagesResponse[0].Recipients = mockRecipientResponse

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SearchMessages(userID, mockFromUserResponse.Email, mockRecipientResponse[0].Email, "folder", "filter").
		Return(mockMessagesResponse, nil)
	userUC.EXPECT().GetInfo(mockFromUserResponse.UserID).Return(mockFromUserResponse, nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(mockMessagesResponse[0].MessageID, mockFromUserResponse.UserID).Return([]uint64{mockRecipientResponse[0].UserID}, nil)
	userUC.EXPECT().GetInfo(mockRecipientResponse[0].UserID).Return(&mockRecipientResponse[0], nil)

	response, err := mailH.SearchMessages(userID, mockFromUserResponse.Email, mockRecipientResponse[0].Email, "folder", "filter")
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, mockMessagesResponse, response)
	}
}

func TestUseCase_SearchRecipients(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var mockRecipientsResponse []models.UserInfo
	generateFakeData(&mockRecipientsResponse)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().SearchRecipients(userID).Return(mockRecipientsResponse, nil)

	response, err := mailH.SearchRecipients(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, mockRecipientsResponse, response)
	}
}

func TestUseCase_GetAttachInfo(t *testing.T) {
	cfg := createConfig()

	attachID := uint64(1)
	userID := uint64(1)
	var mockAttachResponse *models.AttachmentInfo
	generateFakeData(&mockAttachResponse)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().GetAttach(attachID, userID).Return(mockAttachResponse, nil)

	response, err := mailH.GetAttachInfo(attachID, userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, mockAttachResponse, response)
	}
}

func TestUseCase_GetAttach(t *testing.T) {
	cfg := createConfig()

	attachID := uint64(1)
	userID := uint64(1)
	var mockAttachResponse *models.AttachmentInfo
	var fileS3Response *models.S3File
	generateFakeData(&mockAttachResponse)
	generateFakeData(&fileS3Response)
	mockAttachResponse.FileData = fileS3Response.Data

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	fileUC := mockFileUC.NewMockUseCaseI(ctrl)
	mailH := New(cfg, mailRepo, userUC, fileUC)

	mailRepo.EXPECT().GetAttach(attachID, userID).Return(mockAttachResponse, nil)
	fileUC.EXPECT().Get(cfg.S3.S3AttachBucket, mockAttachResponse.S3FName).Return(fileS3Response, nil)

	response, err := mailH.GetAttach(attachID, userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, mockAttachResponse, response)
	}
}
