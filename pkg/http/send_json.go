package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"net/http"
)

func SendJSON(w http.ResponseWriter, r *http.Request, status int, dataStruct any) {
	dataJSON, err := json.Marshal(dataStruct)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to marshal : %w", err))
		return
	}

	w.Header().Set("Content-Type", pkg.ContentTypeJSON)
	w.WriteHeader(status)

	_, err = w.Write(dataJSON)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}

func SendImage(w http.ResponseWriter, r *http.Request, status int, data []byte) {
	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.WriteHeader(status)

	_, err := w.Write(data)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}
