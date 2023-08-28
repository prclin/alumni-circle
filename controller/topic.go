package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"github.com/prclin/alumni-circle/util"
)

func init() {
	topic := core.ContextRouter.Group("/topic")
	topic.POST("", PostTopic)
}

// PostTopic 创建话题
func PostTopic(c *gin.Context) {
	//获取token
	token, err := c.Cookie("token")
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//解析token
	_, err = util.ParseToken(token)
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取话题
	var body struct {
		Name  string  `json:"name" binding:"required"`
		Extra *string `json:"extra" binding:"omitempty,json"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//创建话题
	topic := model.TTopic{Name: body.Name, Extra: body.Extra}
	tTopic, err := service.CreateTopic(topic)
	if err != nil {
		global.Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, tTopic)
}
