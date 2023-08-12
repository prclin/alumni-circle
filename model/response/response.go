package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
Response 响应函数
*/

// Info 自定义响应体
func Info(c *gin.Context, httpStatus int, code int, message string, data interface{}) {
	c.JSON(httpStatus, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    SUCCESS,
		"message": "success",
		"data":    data,
	})
}

// Failed 失败响应
func Failed(c *gin.Context, err string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    FAILED,
		"message": "failed",
		"data":    err,
	})
}

// Error 错误响应
func Error(c *gin.Context, err string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    Server_Error,
		"message": "server error",
		"data":    err,
	})
}
