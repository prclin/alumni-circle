package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"github.com/prclin/alumni-circle/util"
	"net/http"
	"strconv"
)

func init() {
	_break := core.ContextRouter.Group("/break")
	_break.POST("", PostBreak)
	_break.PUT("/:id", PutBreak)
	_break.DELETE("/:id", DeleteBreak)
}

// DeleteBreak 删除课间
//
// 只能删除账户自己发布的课间，如果删除其他人的课间，不会报错，但是删除不成功
func DeleteBreak(c *gin.Context) {
	//获取并解析token
	claims, err := util.ParseToken(util.IgnoreError(c.Cookie("token")))
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取id
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//删除
	err = service.DeleteBreak(model.TBreak{Id: id, AccountId: claims.Id})
	if err != nil {
		global.Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "删除成功"})
}

// PutBreak 更新课间
//
// 目前暂时只支持更新可见性
//
// 只能更新自己帖子的可见性，如果更新的是别人的帖子不会返回错误，但更新不会成功
func PutBreak(c *gin.Context) {
	//获取id
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取并解析token
	claims, err := util.ParseToken(util.IgnoreError(c.Cookie("token")))
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取参数
	var body struct {
		Visibility *uint8 `json:"visibility" binding:"required,min=0,max=3"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//更新
	err = service.UpdateBreakVisibility(model.TBreak{Id: id, AccountId: claims.Id, Visibility: *body.Visibility})
	if err != nil {
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "更新成功"})
}

// PostBreak 发布课间
func PostBreak(c *gin.Context) {
	//获取token
	token, err := c.Cookie("token")
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//解析token
	claims, err := util.ParseToken(token)
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//参数绑定
	var body struct {
		Content    string   `json:"content" binding:"required,max=2000"`
		Visibility *uint8   `json:"visibility" binding:"required,min=0,max=3"`
		State      *uint8   `json:"state" binding:"required,eq=1"` //发布时状态只能为1
		Extra      *string  `json:"extra"`
		ShotIds    []uint64 `json:"shot_ids" binding:"required,max=9"`
		TopicIds   []uint64 `json:"topic_ids" binding:"required"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//发布课间
	tBreak := model.TBreak{
		AccountId:  claims.Id,
		Content:    body.Content,
		Visibility: *body.Visibility,
		State:      *body.State,
		Extra:      body.Extra,
	}
	_break, err := service.PublishBreak(tBreak, body.ShotIds, body.TopicIds)
	if err != nil {
		global.Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, _break)
}
