package _error

import "net/http"

var (
	TokenNotFoundError   = New(http.StatusUnauthorized, "token未提供")
	PathParamFormatError = NewClientError("路径参数格式错误")
	InvalidTokenError    = New(http.StatusUnauthorized, "无效token")
	InternalServerError  = New(http.StatusInternalServerError, "服务器内部错误，请稍后重试!")
	QueryParamError      = New(http.StatusBadRequest, "query参数错误")
)

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

func NewClientError(message string) error {
	return &BusinessError{Code: http.StatusBadRequest, Message: message}
}

func NewServerError(message string) error {
	return &BusinessError{Code: http.StatusInternalServerError, Message: message}
}
