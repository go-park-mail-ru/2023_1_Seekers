package repository

import "github.com/go-park-mail-ru/2023_1_Seekers/internal/models"

//go:generate mockgen -destination=./mocks/mockrepo.go -source=./interface.go -package=mocks

type FileStorageRepoI interface {
	Get(bName, fName string) (*models.S3File, error)
	Upload(file *models.S3File) error
}
