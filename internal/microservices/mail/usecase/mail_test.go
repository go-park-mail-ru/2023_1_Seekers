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
			Email:     "max03@mailbox.ru",
		},
		{
			UserID:    userID,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     "valera03@mailbox.ru",
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

	for i, _ := range output {
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
			Email:     "max03@mailbox.ru",
		},
		{
			UserID:    userID,
			FirstName: "valera",
			LastName:  "vinokurshin",
			Email:     "valera03@mailbox.ru",
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

func TestUseCase_SendMessage(t *testing.T) {
	cfg := createConfig()
	userID := uint64(1)
	fromUser := models.User{
		UserID:     userID,
		Email:      "valera03@mailbox.ru",
		Password:   "",
		FirstName:  "valera",
		LastName:   "vinokurshin",
		Avatar:     "default",
		IsExternal: false,
	}

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

	userUC.EXPECT().GetByID(userID).Return(&fromUser, nil)
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

//
//func TestUseCase_SendFailedSendingMessage(t *testing.T) {
//	cfg := createConfig()
//
//	supportEmail := "support@mailbox.ru"
//	userEmail := "valera03@mailbox.ru"
//	mockUserResponse := []models.UserInfo{
//		{
//			UserID:    1,
//			FirstName: "valera",
//			LastName:  "vinokurshin",
//			Email:     userEmail,
//		},
//		{
//			UserID:    2,
//			FirstName: "Suppot",
//			LastName:  "Supportov",
//			Email:     supportEmail,
//		},
//	}
//	invalidEmails := []string{"123123123@mailbox.ru"}
//	formMessage := models.FormMessage{
//		Recipients: []string{userEmail},
//		Title:      "Ваше сообщение не доставлено",
//		Text: "Это письмо создано автоматически сервером Mailbox.ru, отвечать на него не нужно.\n\n" +
//			"К сожалению, Ваше письмо не может быть доставлено одному или нескольким получателям:\n" +
//			strings.Join(invalidEmails, "\n") + "\n\nРекомендуем Вам проверить корректность указания адресов получателей.",
//		ReplyToMessageID: nil,
//	}
//	mockFolderResponse := []models.Folder{
//		{
//			FolderID:  1,
//			UserID:    mockUserResponse[0].UserID,
//			LocalName: "inbox",
//			Name:      "Входящие",
//		},
//		{
//			FolderID:  2,
//			UserID:    mockUserResponse[1].UserID,
//			LocalName: "outbox",
//			Name:      "Исходящие",
//		},
//	}
//	user2folder := []models.User2Folder{
//		{
//			UserID:   mockUserResponse[1].UserID,
//			FolderID: mockFolderResponse[1].FolderID,
//		},
//		{
//			UserID:   mockUserResponse[0].UserID,
//			FolderID: mockFolderResponse[0].FolderID,
//		},
//	}
//	newMessage := models.MessageInfo{
//		Title:            formMessage.Title,
//		CreatedAt:        pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
//		Text:             formMessage.Text,
//		ReplyToMessageID: formMessage.ReplyToMessageID,
//	}
//	messageSelected := newMessage
//	messageSelected.FromUser = mockUserResponse[1]
//	messageSelected.ReplyToMessageID = nil
//
//	t.Parallel()
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
//	userUC := mockUserUC.NewMockUseCaseI(ctrl)
//	mailH := New(cfg, mailRepo, userUC)
//
//	userUC.EXPECT().GetInfoByEmail(userEmail).Return(&mockUserResponse[0], nil)
//	userUC.EXPECT().GetInfoByEmail(supportEmail).Return(&mockUserResponse[1], nil)
//	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[1].UserID, "outbox").Return(&mockFolderResponse[1], nil)
//	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[0].UserID, "inbox").Return(&mockFolderResponse[0], nil)
//	mailRepo.EXPECT().InsertMessage(mockUserResponse[1].UserID, &newMessage, user2folder).Return(nil).SetArg(1, messageSelected)
//	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[1].UserID, messageSelected.MessageID).Return(&messageSelected, nil)
//	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)
//	mailRepo.EXPECT().SelectRecipientsByMessage(messageSelected.MessageID, mockUserResponse[1].UserID).Return([]uint64{mockUserResponse[0].UserID}, nil)
//	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
//
//	err := mailH.SendFailedSendingMessage(userEmail, invalidEmails)
//	causeErr := pkgErr.Cause(err)
//
//	if causeErr != nil {
//		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
//	}
//}
//
//func TestUseCase_SendWelcomeMessage(t *testing.T) {
//	cfg := createConfig()
//
//	supportEmail := "support@mailbox.ru"
//	userEmail := "valera03@mailbox.ru"
//	mockUserResponse := []models.UserInfo{
//		{
//			UserID:    1,
//			FirstName: "valera",
//			LastName:  "vinokurshin",
//			Email:     userEmail,
//		},
//		{
//			UserID:    2,
//			FirstName: "Suppot",
//			LastName:  "Supportov",
//			Email:     supportEmail,
//		},
//	}
//	formMessage := models.FormMessage{
//		Recipients: []string{userEmail},
//		Title:      "Добро пожаловать в почту Mailbox",
//		Text: "Это письмо создано автоматически сервером Mailbox.ru, отвечать на него не нужно.\n" +
//			"Поздравляем Вас с присоединением к нашей почте. Уверены, что вы останетесь довольны ее использованием!",
//		ReplyToMessageID: nil,
//	}
//	mockFolderResponse := []models.Folder{
//		{
//			FolderID:  1,
//			UserID:    mockUserResponse[0].UserID,
//			LocalName: "inbox",
//			Name:      "Входящие",
//		},
//		{
//			FolderID:  2,
//			UserID:    mockUserResponse[1].UserID,
//			LocalName: "outbox",
//			Name:      "Исходящие",
//		},
//	}
//	user2folder := []models.User2Folder{
//		{
//			UserID:   mockUserResponse[1].UserID,
//			FolderID: mockFolderResponse[1].FolderID,
//		},
//		{
//			UserID:   mockUserResponse[0].UserID,
//			FolderID: mockFolderResponse[0].FolderID,
//		},
//	}
//	newMessage := models.MessageInfo{
//		Title:            formMessage.Title,
//		CreatedAt:        pkg.GetCurrentTime(cfg.Logger.LogsTimeFormat),
//		Text:             formMessage.Text,
//		ReplyToMessageID: formMessage.ReplyToMessageID,
//	}
//	messageSelected := newMessage
//	messageSelected.FromUser = mockUserResponse[1]
//	messageSelected.ReplyToMessageID = nil
//
//	t.Parallel()
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
//	userUC := mockUserUC.NewMockUseCaseI(ctrl)
//	mailH := New(cfg, mailRepo, userUC)
//
//	userUC.EXPECT().GetInfoByEmail(userEmail).Return(&mockUserResponse[0], nil)
//	userUC.EXPECT().GetInfoByEmail(supportEmail).Return(&mockUserResponse[1], nil)
//	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[1].UserID, "outbox").Return(&mockFolderResponse[1], nil)
//	mailRepo.EXPECT().SelectFolderByUserNFolderSlug(mockUserResponse[0].UserID, "inbox").Return(&mockFolderResponse[0], nil)
//	mailRepo.EXPECT().InsertMessage(mockUserResponse[1].UserID, &newMessage, user2folder).Return(nil).SetArg(1, messageSelected)
//	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[1].UserID, messageSelected.MessageID).Return(&messageSelected, nil)
//	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)
//	mailRepo.EXPECT().SelectRecipientsByMessage(messageSelected.MessageID, mockUserResponse[1].UserID).Return([]uint64{mockUserResponse[0].UserID}, nil)
//	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
//
//	err := mailH.SendWelcomeMessage(userEmail)
//	causeErr := pkgErr.Cause(err)
//
//	if causeErr != nil {
//		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
//	}
//}
//
//func TestUseCase_MarkMessageAsSeen(t *testing.T) {
//	cfg := createConfig()
//
//	userID := uint64(1)
//	messageID := uint64(1)
//	state := "seen"
//	stateValue := true
//
//	mockUserResponse := []models.UserInfo{
//		{
//			UserID:    userID,
//			FirstName: "valera",
//			LastName:  "vinokurshin",
//			Email:     "valera03@mailbox.ru",
//		},
//		{
//			UserID:    2,
//			FirstName: "max",
//			LastName:  "vlasov",
//			Email:     "max03@mailbox.ru",
//		},
//	}
//	mockMessageResponse := &models.MessageInfo{
//		FromUser:   mockUserResponse[0],
//		MessageID:  messageID,
//		Recipients: nil,
//		Title:      "test",
//		CreatedAt:  "2023-01-29",
//		Text:       "test text",
//		Seen:       stateValue,
//	}
//
//	t.Parallel()
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mailRepo := mockMailRepo.NewMockMailRepoI(ctrl)
//	userUC := mockUserUC.NewMockUseCaseI(ctrl)
//	mailH := New(cfg, mailRepo, userUC)
//
//	mailRepo.EXPECT().UpdateMessageState(userID, messageID, state, stateValue).Return(nil)
//	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[0].UserID, mockMessageResponse.MessageID).Return(mockMessageResponse, nil)
//	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
//	mailRepo.EXPECT().SelectRecipientsByMessage(mockMessageResponse.MessageID, mockUserResponse[0].UserID).Return([]uint64{mockUserResponse[1].UserID}, nil)
//	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)
//
//	response, err := mailH.MarkMessageAsSeen(userID, messageID)
//	causeErr := pkgErr.Cause(err)
//
//	if causeErr != nil {
//		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
//	} else {
//		require.Equal(t, mockMessageResponse, response)
//	}
//}

func TestUseCase_MarkMessageAsUnseen(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	messageID := uint64(1)
	state := "seen"
	stateValue := false

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

	mailRepo.EXPECT().UpdateMessageState(userID, messageID, state, stateValue).Return(nil)
	mailRepo.EXPECT().SelectMessageByUserNMessage(mockUserResponse[0].UserID, mockMessageResponse.MessageID).Return(mockMessageResponse, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[0].UserID).Return(&mockUserResponse[0], nil)
	mailRepo.EXPECT().SelectRecipientsByMessage(mockMessageResponse.MessageID, mockUserResponse[0].UserID).Return([]uint64{mockUserResponse[1].UserID}, nil)
	userUC.EXPECT().GetInfo(mockUserResponse[1].UserID).Return(&mockUserResponse[1], nil)

	response, err := mailH.MarkMessageAsUnseen(userID, messageID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, mockMessageResponse, response)
	}
}
