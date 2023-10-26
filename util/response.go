package util

import (
	"github.com/gin-gonic/gin"
	_error "github.com/prclin/alumni-circle/error"
	"github.com/prclin/alumni-circle/model"
	"net/http"
)

func Error(c *gin.Context, err error) {
	businessError, ok := err.(*_error.BusinessError)
	if !ok {
		Write(c, &model.Response[any]{Code: http.StatusInternalServerError, Message: err.Error(), Data: nil})
		return
	}
	Write(c, &model.Response[any]{Code: businessError.Code, Message: businessError.Message, Data: nil})
}

func Ok[T any](c *gin.Context, message string, data T) {
	Write(c, &model.Response[any]{Code: http.StatusOK, Message: message, Data: data})
}
func Write[T any](c *gin.Context, response *model.Response[T]) {
	c.JSON(http.StatusOK, response)
}
