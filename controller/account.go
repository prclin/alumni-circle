package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	_error "github.com/prclin/alumni-circle/error"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"github.com/prclin/alumni-circle/util"
	"io"
	"net/http"
	"strconv"
	"time"
)

func init() {
	account := core.ContextRouter.Group("/account")
	account.GET("/info/:id", GetAccountInfo)
	account.PUT("/info", PutAccountInfo)
	account.GET("/photo/:id", GetAccountPhoto)
	account.PUT("/photo", PutAccountPhoto)
	account.PUT("/tag", PutAccountTag)
	account.GET("/tag/:id", GetAccountTag)
	account.POST("/follow", PostFollow)
	account.DELETE("/follow", DeleteFollow)
}

// DeleteFollow 取关
func DeleteFollow(c *gin.Context) {
	//获取token
	token, err := c.Cookie("token")
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//解析token
	claims, err := util.ParseToken(token)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取参数
	var body struct {
		FolloweeId uint64 `json:"followee_id" binding:"required"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//取关
	err = service.RevokeFollow(model.TFollow{
		FollowerId: claims.Id,
		FolloweeId: body.FolloweeId,
	})
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "取关成功"})
}

// PostFollow 关注
func PostFollow(c *gin.Context) {
	//获取token
	token, err := c.Cookie("token")
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//解析token
	claims, err := util.ParseToken(token)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取参数
	var body struct {
		FolloweeId uint64  `json:"followee_id" binding:"required"`
		Extra      *string `json:"extra" binding:"omitempty,json"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//关注
	err = service.FollowAccount(model.TFollow{
		FollowerId: claims.Id,
		FolloweeId: body.FolloweeId,
		Extra:      body.Extra,
	})
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "关注成功"})
}

// GetAccountTag 获取账户兴趣标签
func GetAccountTag(c *gin.Context) {
	//获取id
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}

	tags, err := service.GetAccountTag(id)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, tags)
}

// PutAccountTag 修改兴趣标签
func PutAccountTag(c *gin.Context) {
	//获取token
	token, err := c.Cookie("token")
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//解析token
	claims, err := util.ParseToken(token)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取标签列表
	var tags []uint32
	err = c.ShouldBindJSON(&tags)
	if err != nil {
		if err == io.EOF {
			tags = make([]uint32, 0, 0)
		} else {
			Logger.Debug(err)
			model.Client(c)
			return
		}
	}
	//修改
	tag, err := service.UpdateAccountTag(claims.Id, tags)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, tag)
}

// GetAccountInfo 获取账户信息
func GetAccountInfo(c *gin.Context) {
	//获取id
	param := c.Param("id")
	acquiree, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		Logger.Debug(err)
		util.Error(c, _error.PathParamFormatError)
		return
	}
	//获取token
	token, err := c.Cookie("token")
	if err != nil {
		Logger.Debug(err)
		util.Error(c, _error.TokenNotFoundError)
		return
	}
	//解析token
	claims, err := util.ParseToken(token)
	if err != nil { //token错误
		Logger.Debug(err)
		util.Error(c, _error.InvalidTokenError)
		return
	}

	//获取账户信息
	accountInfo, err := service.GetAccountInfo(claims.Id, acquiree)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Ok(c, "获取成功", accountInfo)
}

// PutAccountInfo 修改账户信息
func PutAccountInfo(c *gin.Context) {
	//获取cookie
	token, err := c.Cookie("token")
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//解析token
	claims, err := util.ParseToken(token)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//参数绑定
	var info struct {
		Nickname      string  `json:"nickname" binding:"required"`
		AvatarURL     string  `json:"avatar_url" binding:"required,url"`
		BackgroundURL string  `json:"background_url" binding:"required,url"`
		Sex           *uint8  `json:"sex" binding:"required,min=0,max=1"`
		Brief         *string `json:"brief" binding:"required"`
		Birthday      string  `json:"birthday"`
		Extra         *string `json:"extra" binding:"omitempty,json"`
	}
	err = c.ShouldBindJSON(&info)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}

	var birthday *time.Time
	if info.Birthday != "" {
		pTime, err := time.Parse(time.DateTime, info.Birthday)
		birthday = &pTime
		if err != nil {
			Logger.Debug(err)
			util.Error(c, _error.NewClientError("生日格式错误"))
			return
		}
	}

	//修改信息
	infoR, err := service.UpdateAccountInfo(model.TAccountInfo{
		Id:            claims.Id,
		AvatarURL:     info.AvatarURL,
		BackgroundURL: info.BackgroundURL,
		Nickname:      info.Nickname,
		Sex:           *info.Sex,
		Brief:         *info.Brief,
		Birthday:      birthday,
		Extra:         info.Extra,
	})
	if err != nil {
		util.Error(c, err)
		return
	}
	Logger.Debug(infoR)
	util.Ok(c, "修改成功", infoR)
}

// GetAccountPhoto 获取照片墙
func GetAccountPhoto(c *gin.Context) {
	//获取参数
	id := c.Param("id")
	accountId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}

	//获取照片墙
	wall, err := service.GetPhotoWall(accountId)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, util.Ternary(len(wall) == 0, []model.Photo{}, wall))
}

// PutAccountPhoto 修改照片墙
func PutAccountPhoto(c *gin.Context) {
	//获取token
	token, err := c.Cookie("token")
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//解析token
	claims, err := util.ParseToken(token)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取照片墙
	var photos []model.TPhotoBinding
	err = c.ShouldBindJSON(&photos)
	if err != nil && !errors.Is(err, io.EOF) {
		Logger.Debug(err)
		model.Client(c)
		return
	}

	//修改
	err = service.UpdateAccountPhoto(claims.Id, photos)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "修改成功"})
}
