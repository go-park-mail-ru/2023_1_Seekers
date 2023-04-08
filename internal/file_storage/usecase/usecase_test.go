package usecase

import (
	"github.com/go-faker/faker/v4"
	mockS3Repo "github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage/repository/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
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

func TestUseCase_Get(t *testing.T) {
	var fakeS3File *models.S3File
	generateFakeData(&fakeS3File)
	var bName, fName string
	generateFakeData(&bName)
	generateFakeData(&fName)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3Repo := mockS3Repo.NewMockRepoI(ctrl)
	s3UC := New(s3Repo)

	s3Repo.EXPECT().Get(bName, fName).Return(fakeS3File, nil)
	file, err := s3UC.Get(bName, fName)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, file, fakeS3File)
	}
}

func TestUseCase_Upload(t *testing.T) {
	var fakeS3File *models.S3File
	generateFakeData(&fakeS3File)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3Repo := mockS3Repo.NewMockRepoI(ctrl)
	s3UC := New(s3Repo)

	s3Repo.EXPECT().Upload(fakeS3File).Return(nil)
	err := s3UC.Upload(fakeS3File)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}
