package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/po"
	"github.com/prclin/alumni-circle/model/request"
	. "github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
	"github.com/prclin/alumni-circle/util"
	"net/http"
	"strconv"
)

func init() {
	tag := core.ContextRouter.Group("/tag")
	tag.POST("", PostTag)
	tag.PUT("/:id", PutTag)
	tag.DELETE("/:id", DeleteTag)
	tag.GET("/list", GetTagList)
}

// GetTagList 获取兴趣标签列表
func GetTagList(c *gin.Context) {
	//获取参数
	var query struct {
		request.Pagination
		State *uint8 `form:"state" binding:"max=1"`
	}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}
	//获取
	tags, err := service.GetTagList(query.Pagination, query.State)
	if err != nil {
		Logger.Debug(err)
		Server(c)
		return
	}
	Ok(c, util.Ternary(len(tags) == 0, make([]po.TTag, 0, 0), tags))
}

// PostTag 创建兴趣标签
func PostTag(c *gin.Context) {
	//获取参数
	var body struct {
		Name  string  `json:"name" binding:"required"`
		State *uint8  `json:"state" binding:"required,min=0,max=1"`
		Extra *string `json:"extra" binding:"omitempty,json"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}
	//创建
	tag, err := service.CreateTag(po.TTag{
		Name:  body.Name,
		State: *body.State,
		Extra: body.Extra,
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
		State *uint8  `json:"state" binding:"required,min=0,max=1"`
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
		State: *body.State,
		Extra: body.Extra,
	})
	if err != nil {
		Logger.Debug(err)
		Server(c)
		return
	}
	Ok(c, tag)
}

// DeleteTag 删除tag
func DeleteTag(c *gin.Context) {
	//获取参数
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}
	//删除
	err = service.DeleteTag(uint32(id))
	if err != nil {
		Logger.Debug(err)
		Server(c)
		return
	}
	Write(c, Response[any]{Code: http.StatusOK, Message: "删除成功"})
}
