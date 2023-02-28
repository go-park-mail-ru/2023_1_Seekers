package pkg

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"net/http"
)

func SendError(w http.ResponseWriter, error *errors.JsonErr) error {
	err := SendJson(w, error.Code, error)
	if err != nil {
		return fmt.Errorf("cant send error %w", err)
	}
	return nil
}
