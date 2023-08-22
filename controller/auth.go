package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/entity"
	. "github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
)

func init() {
	auth := core.ContextRouter.Group("/auth")
	auth.POST("/sign_up", EmailSignUp) //注册
}

/*
EmailSignUp 账户注册（邮箱）
参数：邮箱 验证码 密码
*/
func EmailSignUp(c *gin.Context) {
	//参数结构
	var body struct {
		Email    string `form:"email" binding:"required"`
		Captcha  string `form:"captcha" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	//参数绑定
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}

	//注册逻辑
	account := entity.Account{
		Email:    body.Email,
		Password: body.Password,
	}
	res := service.EmailSignUp(account, body.Captcha)
	Write(c, res)
}
