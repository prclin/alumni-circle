package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
Response 基础响应
*/
type Response[T any] struct {
	//状态代码
	Code int32 `json:"code"`
	//状态信息
	Message string `json:"message"`
	//数据
	Data T `json:"data"`
}

// Ok 200+
func Ok[T any](c *gin.Context, data T) {
	Write(c, Response[T]{Code: http.StatusOK, Message: "OK", Data: data})
}

// Client 400+
func Client[T any](c *gin.Context, data T) {
	Write(c, Response[T]{Code: http.StatusBadRequest, Message: "Bad Request", Data: data})
}

// Server 500+
func Server[T any](c *gin.Context, data T) {
	Write(c, Response[T]{Code: http.StatusInternalServerError, Message: "Internal Server Error", Data: data})
}

// Write 自定义写入
func Write[T any](c *gin.Context, response Response[T]) {
	c.JSON(http.StatusOK, response)
}
