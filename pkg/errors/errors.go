package errors

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
