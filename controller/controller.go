package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/model"
)

/*
Init 为了显示执行副作用引入
*/
func Init() {
	//just empty
}

func init() {
	core.ContextRouter.GET("/health", func(context *gin.Context) {
		context.JSON(200, model.Response[any]{Code: 200, Message: "ok"})
	})
}
