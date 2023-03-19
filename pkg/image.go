package pkg

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
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
		return nil, user.ErrInvalidForm
	}

	if ok := CheckImageContentType(http.DetectContentType(img.Data)); !ok {
		return nil, user.ErrWrongContentType
	}
	return &img, nil
}
