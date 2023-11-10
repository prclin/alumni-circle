package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	_error "github.com/prclin/alumni-circle/error"
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
	file.POST("/images", PostImages)
}

//PostImages 多图片上传，迭代中
func PostImages(context *gin.Context) {
	//解析token
	claims, err := util.GetTokenClaims(context)
	if err != nil {
		util.Error(context, err)
		return
	}
	// Multipart form
	form, err := context.MultipartForm()
	if err != nil {
		model.Client(context)
		return
	}

	//上传
	files := form.File["files"]
	images := make([]model.TImage, 0, len(files))
	for _, file := range files {
		//校验文件类型
		cType := file.Header.Get("Content-Type")
		if strings.Split(cType, "/")[0] != "image" {
			global.Logger.Debug("文件类型错误!")
			util.Error(context, _error.NewClientError("不支持的文件类型"))
			return
		}
		//重命名文件
		file.Filename = strconv.FormatInt(time.Now().UnixNano(), 10) /*当前时间戳*/ + "-" + strconv.FormatUint(claims.Id, 10) /*账户id*/ + "." + strings.Split(cType, "/")[1] /*扩展名*/
		//上传文件到oss
		image, err := service.UploadFileToOSS("alumni-circle", "development", file) //开发阶段，此处目录为development
		if err != nil {
			global.Logger.Debug(err)
			util.Error(context, _error.InternalServerError)
			return
		}
		images = append(images, image)
	}

	util.Ok(context, "上传成功", images)
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
