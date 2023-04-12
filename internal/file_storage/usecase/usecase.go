package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkgErrors "github.com/pkg/errors"
)

type useCase struct {
	s3Repo file_storage.RepoI
}

func New(s3R file_storage.RepoI) file_storage.UseCaseI {
	return &useCase{
		s3Repo: s3R,
	}
}

func (uc *useCase) Get(bName, fName string) (*models.S3File, error) {
	file, err := uc.s3Repo.Get(bName, fName)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed get file : %v")
	}
	return file, nil
}
func (uc *useCase) Upload(file *models.S3File) error {
	err := uc.s3Repo.Upload(file)
	if err != nil {
		return pkgErrors.Wrap(err, "failed upload file")
	}
	return nil
}
