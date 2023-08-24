package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/po"
	. "github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
)

func init() {
	tag := core.ContextRouter.Group("/tag")
	tag.POST("", PostTag)
}

// PostTag 创建兴趣标签
func PostTag(c *gin.Context) {
	//获取参数
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}
	//创建
	tag, err := service.CreateTag(po.TTag{
		Name: body.Name,
	})
	if err != nil {
		Logger.Debug(err)
		Server(c)
		return
	}
	Ok(c, tag)
}
