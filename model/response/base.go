package response

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/prclin/alumni-circle/model/po"
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
func Client(c *gin.Context) {
	Write(c, Response[any]{Code: http.StatusBadRequest, Message: "Bad Request"})
}

// Server 500+
func Server(c *gin.Context) {
	Write(c, Response[any]{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
}

// Write 自定义写入
func Write[T any](c *gin.Context, response Response[T]) {
	c.JSON(http.StatusOK, response)
}

type TokenClaims struct {
	jwt.RegisteredClaims
	po.TAccountInfo
	RoleIds []uint32 `json:"role_ids"`
}
