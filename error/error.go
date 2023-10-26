package _error

import "net/http"

type ErrorInterface interface {
	error
	businessError() string
}

// BusinessError 业务异常, 400<=code<500
type BusinessError struct {
	Code    int
	Message string
}

func (err *BusinessError) Error() string {
	return err.Message
}

func (err *BusinessError) businessError() string {
	return err.Error()
}

func New(code int, message string) error {
	return &BusinessError{Code: code, Message: message}
}

func NewWithCode(code int) error {
	return &BusinessError{Code: code, Message: http.StatusText(code)}
}
