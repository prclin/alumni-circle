package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
)

// 注册路由
func init() {
	file := core.ContextRouter.Group("/file")
	image := file.Group("/image")
	image.POST("/upload", ImageUpload)
}

func ImageUpload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Client(c, err)
		return
	}
	filename := fileHeader.Filename
	file, err := fileHeader.Open()
	if err != nil {
		response.Client(c, err)
		return
	}
	image, err := service.ImageUpload(filename, file)
	if err != nil {
		response.Server(c, err)
		return
	}
	response.Ok(c, image)
}
