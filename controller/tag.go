package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/po"
	. "github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
	"strconv"
)

func init() {
	tag := core.ContextRouter.Group("/tag")
	tag.POST("", PostTag)
	tag.PUT("/:id", PutTag)
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

// PutTag 修改兴趣标签
func PutTag(c *gin.Context) {
	//获取参数
	param := c.Param("id")
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}
	var body struct {
		Name  string  `json:"name" binding:"required"`
		Extra *string `json:"extra" binding:"omitempty,json"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}
	//修改
	tag, err := service.UpdateTag(po.TTag{
		Id:    uint32(id),
		Name:  body.Name,
		Extra: body.Extra,
	})
	if err != nil {
		Logger.Debug(err)
		Server(c)
		return
	}
	Ok(c, tag)
}
