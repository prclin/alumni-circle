package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	websocket2 "github.com/prclin/alumni-circle/websocket"
)

func init() {
	go broker.ProxyHandle()
	core.ContextRouter.GET("/ws", GetWebsocketConnection)
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	broker = websocket2.DefaultBroker()
)

// GetWebsocketConnection websocket连接
func GetWebsocketConnection(c *gin.Context) {
	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.Logger.Debug(err)
		model.Server(c)
		return
	}
	client := websocket2.NewClient(0, connection)
	broker.AddClient(client)
	for {
		var msg websocket2.Message
		_ = connection.ReadJSON(&msg)
		broker.Channel <- msg
	}
}
