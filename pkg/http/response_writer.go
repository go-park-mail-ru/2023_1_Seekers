package http

import "net/http"

type ResponseWriterCode struct {
	http.ResponseWriter
	StatusCode int
}

func NewResponseWriterCode(w http.ResponseWriter) *ResponseWriterCode {
	return &ResponseWriterCode{w, http.StatusOK}
}

func (rw *ResponseWriterCode) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
