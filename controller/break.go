package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/model/entity"
	"github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
	"strconv"
)

// 注册路由
func init() {
	breaks := core.ContextRouter.Group("/break")
	breaks.POST("/create", CreateBreak)
}

// 从Cookie获取用户id(应放在controller/account中,暂存在此)
func GetIdFromCookie(c *gin.Context) (id int, err error) {
	cookie, err := c.Cookie("id")
	if err != nil {
		return 0, err
	}
	id, err = strconv.Atoi(cookie)
	return
}
func CreateBreak(c *gin.Context) {
	aBreak := new(entity.Break)
	if err := c.BindJSON(aBreak); err != nil {
		response.Client(c, err)
		return
	}
	if aBreak.Title == "" || aBreak.Content == "" {
		response.Client(c, "参数不完整")
		return
	}
	accountId, err := GetIdFromCookie(c)
	if err != nil {
		response.NLI(c)
		return
	}
	aBreak.AccountId = accountId
	if err := service.CreateBreak(aBreak); err != nil {
		response.Server(c, err)
	}
	response.Ok(c, aBreak)
}
