package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	pkg "github.com/go-park-mail-ru/2023_1_Seekers/pkg/crypto"
	"path/filepath"

	//"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	mockFileRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/usecase/mocks"
	mockUserRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/repository/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/golang/mock/gomock"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func createConfig() *config.Config {
	cfg := new(config.Config)
	cfg.S3.S3AvatarBucket = "avatars_mailbox_vkcloud"
	cfg.UserService.DefaultAvatar = "default_avatar.png"
	cfg.UserService.MaxImageSize = 2147483648
	cfg.UserService.UserDefaultAvatarSize = 46
	cfg.UserService.AvatarTTFPath = "../../../../cmd/config/wqy-zenhei.ttf"
	cfg.Password.PasswordSaltLen = 0

	return cfg
}

func TestUseCase_Create(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)
	fakeUser.Email = "valera03@mailbox.ru"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(nil, errors.ErrUserNotFound)
	userRepo.EXPECT().Create(fakeUser).Return(fakeUser, nil)

	response, err := userUC.Create(fakeUser)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestUseCase_GetInfo(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)
	output := &models.UserInfo{
		UserID:    fakeUser.UserID,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
		Email:     fakeUser.Email,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil)
	response, err := userUC.GetInfo(fakeUser.UserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, output, response)
	}
}

func TestUseCase_Delete(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().Delete(userID).Return(nil)
	err := userUC.Delete(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	}
}

func TestUseCase_GetByID(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil)
	response, err := userUC.GetByID(fakeUser.UserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestUseCase_GetByEmail(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
	response, err := userUC.GetByEmail(fakeUser.Email)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestUseCase_GetInfoByEmail(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)
	fakeUserInfo := &models.UserInfo{
		UserID:    fakeUser.UserID,
		FirstName: fakeUser.FirstName,
		LastName:  fakeUser.LastName,
		Email:     fakeUser.Email,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
	response, err := userUC.GetInfoByEmail(fakeUser.Email)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeUserInfo, response)
	}
}

func TestUseCase_EditInfo(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	var request *models.UserInfo
	generateFakeData(&fakeUser)
	generateFakeData(&request)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil)
	userRepo.EXPECT().EditInfo(fakeUser.UserID, request).Return(nil)
	userRepo.EXPECT().IsCustomAvatar(fakeUser.UserID).Return(true, nil)

	response, err := userUC.EditInfo(fakeUser.UserID, request)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, *request, *response)
	}
}

func TestUseCase_EditPw(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	var fakePassword *models.EditPasswordRequest
	generateFakeData(&fakeUser)
	generateFakeData(&fakePassword)
	fakePassword.RepeatPw = fakePassword.Password
	fakeUser.Password, _ = pkg.HashPw(fakePassword.PasswordOld, cfg.Password.PasswordSaltLen)
	newHashPw, _ := pkg.HashPw(fakePassword.Password, cfg.Password.PasswordSaltLen)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil)
	userRepo.EXPECT().EditPw(fakeUser.UserID, newHashPw).Return(nil)
	err := userUC.EditPw(fakeUser.UserID, fakePassword)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	}
}

func TestUseCase_EditAvatar(t *testing.T) {
	cfg := createConfig()

	var fakeImage *models.Image
	var fakeUser *models.User
	generateFakeData(&fakeImage)
	generateFakeData(&fakeUser)
	f := models.S3File{
		Bucket: cfg.S3.S3AvatarBucket,
		Name:   fakeUser.Email + filepath.Ext(fakeImage.Name),
		Data:   fakeImage.Data,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().GetByID(fakeUser.UserID).Return(fakeUser, nil)
	fileUC.EXPECT().Upload(&f).Return(nil)
	userRepo.EXPECT().SetAvatar(fakeUser.UserID, f.Name).Return(nil)
	userRepo.EXPECT().SetCustomAvatar(fakeUser.UserID).Return(nil)

	err := userUC.EditAvatar(fakeUser.UserID, fakeImage, true)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	}
}

func TestUseCase_GetAvatar(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	var fakeS3 *models.S3File
	generateFakeData(&fakeUser)
	generateFakeData(&fakeS3)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockUserRepoI(ctrl)
	fileUC := mockFileRepo.NewMockUseCaseI(ctrl)
	userUC := New(cfg, userRepo, fileUC)

	userRepo.EXPECT().GetByEmail(fakeUser.Email).Return(fakeUser, nil)
	fileUC.EXPECT().Get(cfg.S3.S3AvatarBucket, fakeUser.Avatar).Return(nil, errors.ErrNoBucket)
	fileUC.EXPECT().Get(cfg.S3.S3AvatarBucket, cfg.UserService.DefaultAvatar).Return(fakeS3, nil)

	response, err := userUC.GetAvatar(fakeUser.Email)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, &models.Image{
			Name: fakeS3.Name,
			Data: fakeS3.Data,
		}, response)
	}
}
