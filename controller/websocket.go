package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"net/http"
)

func init() {
	ws := core.ContextRouter.Group("/ws")
	ws.GET("", ServeWebSocket)
}

/*
ServeWebSocket 创建websocket连接
*/
func ServeWebSocket(c *gin.Context) {
	// 将http连接升级为websocket连接
	ws, err := global.WebsocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		global.Logger.Debug(err)
		model.Write(c, model.Response[any]{Code: http.StatusBadRequest, Message: "ws: Not a websocket handshake"})
		return
	} else if err != nil {
		global.Logger.Debug(err)
		model.Write(c, model.Response[any]{Code: http.StatusBadRequest, Message: "ws: failed to Upgrade"})
		return
	}
	// 创建会话逻辑
	session, err := service.WebsocketNewSession(ws, "")
	fmt.Println(session)
	// TODO: 将会话加入websocket连接池

	return
}
