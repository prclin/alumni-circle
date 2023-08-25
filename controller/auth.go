package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"net/http"
)

func init() {
	auth := core.ContextRouter.Group("/auth")
	auth.POST("/sign_up", EmailSignUp) //邮箱注册
	auth.PUT("/sign_in", EmailSignIn)  //邮箱登录
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
		model.Client(c)
		return
	}

	//注册逻辑
	account := model.TAccount{
		Email:    body.Email,
		Password: body.Password,
	}
	res := service.EmailSignUp(account, body.Captcha)
	model.Write(c, res)
}

/*
EmailSignIn 账户登录（邮箱）

参数：邮箱 密码
*/
func EmailSignIn(c *gin.Context) {
	var body struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	//参数绑定
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}

	//登录逻辑
	res := service.EmailSignIn(body.Email, body.Password)
	//登录成功回写cookie
	if res.Code == http.StatusOK {
		c.SetCookie("token", *res.Data, -1, "/", "*", false, false)
	}
	model.Write(c, res)
}
