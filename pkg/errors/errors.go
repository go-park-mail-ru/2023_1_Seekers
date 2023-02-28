package errors

type JSONErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *JSONErr) Error() string {
	return err.Message
}

func New(code int, err error) *JSONErr {
	return &JSONErr{
		Code:    code,
		Message: err.Error(),
	}
}

func NewWrappedErr(code int, message string, err error) *JSONErr {
	return &JSONErr{
		Code:    code,
		Message: message + " : " + err.Error(),
	}
}
