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
	_break := core.ContextRouter.Group("/break")
	_break.POST("", PostBreak)
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
