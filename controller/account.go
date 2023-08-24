package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/entity"
	"github.com/prclin/alumni-circle/model/po"
	. "github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
	"github.com/prclin/alumni-circle/util"
	"net/http"
	"strconv"
	"time"
)

func init() {
	account := core.ContextRouter.Group("/account")
	account.GET("/info", GetAccountInfo)
	account.PUT("/info", PutAccountInfo)
	account.GET("/photo/:id", GetAccountPhoto)
}

// GetAccountInfo 获取当前登录账户的信息
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

// PutAccountInfo 修改账户信息
func PutAccountInfo(c *gin.Context) {
	//获取cookie
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
	//参数绑定
	var info struct {
		CampusId  uint32  `json:"campus_id" binding:"required"`
		AvatarURL string  `json:"avatar_url" binding:"required,url"`
		Nickname  string  `json:"nickname" binding:"required"`
		Sex       uint8   `json:"sex" binding:"required,min=1,max=2"`
		Birthday  string  `json:"birthday" binding:"required,datetime=2006-01-02"`
		Extra     *string `json:"extra" binding:"omitempty,json"`
	}
	err = c.ShouldBindJSON(&info)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}

	//修改信息
	err = service.UpdateAccountInfo(po.TAccountInfo{
		Id:        claims.Id,
		CampusId:  info.CampusId,
		AvatarURL: info.AvatarURL,
		Nickname:  info.Nickname,
		Sex:       info.Sex,
		Birthday:  util.IgnoreError(time.Parse(time.DateOnly, info.Birthday)),
		Extra:     info.Extra,
	})
	if err != nil {
		Logger.Debug(err)
		Server(c)
		return
	}
	Write(c, Response[any]{Code: http.StatusOK, Message: "修改成功"})
}

// GetAccountPhoto 获取照片墙
func GetAccountPhoto(c *gin.Context) {
	//获取参数
	id := c.Param("id")
	accountId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		Logger.Debug(err)
		Client(c)
		return
	}

	//获取照片墙
	wall, err := service.GetPhotoWall(accountId)
	if err != nil {
		Logger.Debug(err)
		Server(c)
		return
	}
	Ok(c, util.Ternary(len(wall) == 0, []entity.Photo{}, wall))
}
