package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	mockFileUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/usecase/mocks"
	mockMailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/repository/mocks"
	mockUserUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
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
