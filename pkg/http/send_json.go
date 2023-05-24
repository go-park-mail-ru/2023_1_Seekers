package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	pkgJson "github.com/go-park-mail-ru/2023_1_Seekers/pkg/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, r *http.Request, status int, dataStruct any) {
	dataJSON, err := pkgJson.MarshalEasyJSON(dataStruct)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to marshal : %w", err))
		return
	}

	w.Header().Set("Content-Type", common.ContentTypeJSON)
	w.WriteHeader(status)

	_, err = w.Write(dataJSON)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}

func SendImage(w http.ResponseWriter, r *http.Request, status int, data []byte) {
	w.Header().Set("Content-Type", http.DetectContentType(data))
	imageType := http.DetectContentType(data)
	if imageType == "text/xml; charset=utf-8" {
		imageType = "image/svg+xml"
	}

	w.Header().Set("Content-Type", imageType)
	w.WriteHeader(status)

	_, err := w.Write(data)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}
