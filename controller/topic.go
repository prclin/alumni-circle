package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/service"
)

// 注册路由
func init() {
	topic := core.ContextRouter.Group("/topic")
	topic.POST("/create", CreateTopic)
}

// TODO
func CreateTopic(c *gin.Context) {
	service.CreateTopic()
}
