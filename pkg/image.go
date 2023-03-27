package pkg

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
)

func ReadImage(file multipart.File, header *multipart.FileHeader) (*models.Image, error) {
	img := models.Image{
		Data: make([]byte, header.Size),
		Name: header.Filename,
	}
	_, err := io.ReadFull(file, img.Data)
	if err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInvalidForm, err.Error())
	}

	if ok := CheckImageContentType(http.DetectContentType(img.Data)); !ok {
		return nil, errors.ErrWrongContentType
	}
	return &img, nil
}
