package errors

import (
	"errors"
	"fmt"
)

type JSONError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *JSONError) Error() string {
	return err.Message
}

func New(code int, err error) *JSONError {
	return &JSONError{
		Code:    code,
		Message: err.Error(),
	}
}

func NewWrappedErr(code int, message string, err error) *JSONError {
	return &JSONError{
		Code:    code,
		Message: message + " : " + err.Error(),
	}
}

func UnwrapError(err error) error {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}
	fmt.Println(err)

	return err
}
