package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"github.com/prclin/alumni-circle/util"
	"strconv"
	"strings"
	"time"
)

func init() {
	file := core.ContextRouter.Group("/file")
	file.POST("/image", PostImage)
}

// PostImage 上传图片
func PostImage(c *gin.Context) {
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
	//获取文件
	file, err := c.FormFile("file")
	if err != nil {
		global.Logger.Debug(err)
		model.Client(c)
		return
	}
	//校验文件类型
	cType := file.Header.Get("Content-Type")
	if strings.Split(cType, "/")[0] != "image" {
		global.Logger.Debug("文件类型错误!")
		model.Client(c)
		return
	}

	//重命名文件
	file.Filename = strconv.FormatInt(time.Now().UnixMilli(), 10) /*当前时间戳*/ + "-" + strconv.FormatUint(claims.Id, 10) /*账户id*/ + "." + strings.Split(cType, "/")[1] /*扩展名*/

	//上传文件到oss
	image, err := service.UploadFileToOSS("alumni-circle", "development", file) //开发阶段，此处目录为development
	if err != nil {
		global.Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, image)
}
