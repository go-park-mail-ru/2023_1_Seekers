package pkg

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"net/http"
)

func SendError(w http.ResponseWriter, error *errors.JsonErr) {
	SendJson(w, error.Code, error)
}
