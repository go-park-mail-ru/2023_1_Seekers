package errors

type JsonErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *JsonErr) Error() string {
	return err.Message
}

func New(code int, err error) *JsonErr {
	return &JsonErr{
		Code:    code,
		Message: err.Error(),
	}
}

func NewWrappedErr(code int, message string, err error) *JsonErr {
	return &JsonErr{
		Code:    code,
		Message: message + " : " + err.Error(),
	}
}
