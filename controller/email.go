package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
	"regexp"
)

// 注册路由
func init() {
	email := core.ContextRouter.Group("/email")
	email.GET("/:email", GetVerifyEmail)
}

// 邮箱校验器
var emailRegexp *regexp.Regexp

// 初始全局化变量
func init() {
	//初始化emailRegexp
	reg, err := regexp.Compile("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$")
	if err != nil {
		Logger.Fatal(err)
	}
	emailRegexp = reg
}

// GetVerifyEmail 向用户邮箱发送带有校验码的校验邮件
func GetVerifyEmail(c *gin.Context) {
	//获取参数
	to := c.Param("email")
	//校验参数
	matched := emailRegexp.MatchString(to)
	if !matched {
		c.JSON(200, response.Response[any]{Message: "邮箱格式错误", Code: 400})
		return
	}
	//发送邮件
	if err := service.SendVerifyEmail(to); err != nil {
		c.JSON(200, response.Response[any]{Message: "发送失败", Code: 500})
		return
	}
	c.JSON(200, response.Response[any]{Message: "发送成功", Code: 200})
}
