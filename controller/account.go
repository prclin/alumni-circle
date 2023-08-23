package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	. "github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
	"github.com/prclin/alumni-circle/util"
)

func init() {
	account := core.ContextRouter.Group("/account")
	account.GET("/info", GetAccountInfo)
}

/*
GetAccountInfo 获取当前登录账户的信息

参数： token
*/
func GetAccountInfo(c *gin.Context) {
	//获取token
	token, err := c.Cookie("token")
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}
	//解析token
	claims, err := util.ParseToken(token)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}

	//获取账户信息
	info, err := service.GetAccountInfo(claims.Id)
	if err != nil {
		Logger.Debug(err)
		Server(c)
		return
	}
	Ok(c, info)
}
