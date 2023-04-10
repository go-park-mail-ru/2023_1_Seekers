package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	mockFileRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	mockUserRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/golang/mock/gomock"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestUseCase_Create(t *testing.T) {
	var fakeUser *models.User
	generateFakeData(&fakeUser)
	fakeUser.Email = "valera03@mailbox.ru"

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(userRepo, fileUC)

	userRepo.EXPECT().GetByEmail(fakeUser.Email).Return(nil, errors.ErrUserNotFound)
	userRepo.EXPECT().Create(fakeUser).Return(fakeUser, nil)

	response, err := userUC.Create(fakeUser)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestUseCase_GetInfo(t *testing.T) {
	var fakeUser *models.User
	generateFakeData(&fakeUser)
	output := &models.UserInfo{
		UserID:    fakeUser.UserID,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
		Email:     fakeUser.Email,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil)
	response, err := userUC.GetInfo(fakeUser.UserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, output, response)
	}
}

func TestUseCase_Delete(t *testing.T) {
	userID := uint64(1)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(userRepo, fileUC)

	userRepo.EXPECT().Delete(userID).Return(nil)
	err := userUC.Delete(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestUseCase_GetByID(t *testing.T) {
	var fakeUser *models.User
	generateFakeData(&fakeUser)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil)
	response, err := userUC.GetByID(fakeUser.UserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestUseCase_GetByEmail(t *testing.T) {
	var fakeUser *models.User
	generateFakeData(&fakeUser)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(userRepo, fileUC)

	userRepo.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
	response, err := userUC.GetByEmail(fakeUser.Email)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestUseCase_EditInfo(t *testing.T) {
	var fakeUser *models.User
	var request models.UserInfo
	var fakeS3 *models.S3File
	generateFakeData(&fakeUser)
	generateFakeData(&request)
	generateFakeData(&fakeS3)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil).AnyTimes()
	userRepo.EXPECT().EditInfo(fakeUser.UserID, request).Return(nil)
	userRepo.EXPECT().IsCustomAvatar(fakeUser.UserID).Return(false, nil)
	userRepo.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
	fileUC.EXPECT().Get(config.S3AvatarBucket, fakeUser.Avatar).Return(fakeS3, nil)
	fileUC.EXPECT().Upload(gomock.Any()).Return(nil)
	userRepo.EXPECT().SetAvatar(fakeUser.UserID, gomock.Any()).Return(nil)

	response, err := userUC.EditInfo(fakeUser.UserID, request)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, request, *response)
	}
}

func TestUseCase_EditAvatar(t *testing.T) {
	var fakeImage *models.Image
	var fakeUser *models.User
	generateFakeData(&fakeImage)
	generateFakeData(&fakeUser)
	f := models.S3File{
		Bucket: config.S3AvatarBucket,
		Name:   fakeUser.Email + filepath.Ext(fakeImage.Name),
		Data:   fakeImage.Data,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil)
	fileUC.EXPECT().Upload(&f).Return(nil)
	userRepo.EXPECT().SetAvatar(fakeUser.UserID, f.Name).Return(nil)
	userRepo.EXPECT().SetCustomAvatar(fakeUser.UserID).Return(nil)

	err := userUC.EditAvatar(fakeUser.UserID, fakeImage, true)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestUseCase_GetAvatar(t *testing.T) {
	var fakeUser *models.User
	var fakeS3 *models.S3File
	generateFakeData(&fakeUser)
	generateFakeData(&fakeS3)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(userRepo, fileUC)

	userRepo.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
	fileUC.EXPECT().Get(config.S3AvatarBucket, fakeUser.Avatar).Return(nil, errors.ErrNoBucket)
	fileUC.EXPECT().Get(config.S3AvatarBucket, config.DefaultAvatar).Return(fakeS3, nil)

	response, err := userUC.GetAvatar(fakeUser.Email)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, &models.Image{
			Name: fakeS3.Name,
			Data: fakeS3.Data,
		}, response)
	}
}

func TestUseCase_EditPw(t *testing.T) {
	config.PasswordSaltLen = 0
	var fakeUser *models.User
	var fakeForm models.EditPasswordRequest
	generateFakeData(&fakeUser)
	generateFakeData(&fakeForm)
	fakeForm.RepeatPw = fakeForm.Password
	var err error
	fakeUser.Password, err = pkg.HashPw(fakeForm.PasswordOld)
	if err != nil {
		t.Fatalf("error while hashing pw")
	}
	newPw, err := pkg.HashPw(fakeForm.Password)
	if err != nil {
		t.Fatalf("error while hashing pw")
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil)
	userRepo.EXPECT().EditPw(fakeUser.UserID, newPw).Return(nil)

	err = userUC.EditPw(fakeUser.UserID, fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}
