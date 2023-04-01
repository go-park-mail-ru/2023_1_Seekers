package usecase

import (
	fileStorageUC "github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage/usecase/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	userRepoMock "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUseCase_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := userRepoMock.NewMockRepoI(ctrl)
	fStorageUC := fileStorageUC.NewMockUseCaseI(ctrl)

	input := &models.User{
		Email:     "test@test.com",
		Password:  "123456",
		FirstName: "test",
		LastName:  "test",
		Avatar:    "default_avatar",
	}

	output := &models.User{
		UserID:    123,
		Email:     "test@test.com",
		Password:  "123456",
		FirstName: "test",
		LastName:  "test",
		Avatar:    "default_avatar",
	}

	userRepo.EXPECT().Create(input).Return(output, nil)
	userRepo.EXPECT().GetByEmail(input.Email).Return(nil, errors.ErrUserExists)
	userUC := New(userRepo, fStorageUC)
	res, err := userUC.Create(input)

	require.Nil(t, err)
	require.Equal(t, output, res)
}

func TestUseCase_Delete(t *testing.T) {

}

func TestUseCase_EditAvatar(t *testing.T) {

}

func TestUseCase_EditInfo(t *testing.T) {

}

func TestUseCase_GetAvatar(t *testing.T) {

}

func TestUseCase_EditPw(t *testing.T) {

}
