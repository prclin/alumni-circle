package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/messaging"
	websocket2 "github.com/prclin/alumni-circle/websocket"
	"net/http"
)

func init() {
	go broker.ProxyHandle()
	core.ContextRouter.GET("/ws", GetWebsocketConnection)
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	broker = websocket2.DefaultBroker()
)

// GetWebsocketConnection websocket连接
func GetWebsocketConnection(c *gin.Context) {
	stompBroker := messaging.NewStompBroker()
	a := func(context *messaging.Context) {
		fmt.Println(context.Frame.String())
	}
	stompBroker.Handle("/test", a)
	stompBroker.Subscribe("/test", a)
	err := stompBroker.ServeOverHttp(c.Writer, c.Request)
	if err != nil {
		global.Logger.Debug(err)
	}
	//client := websocket2.NewClient(0, connection)
	//broker.AddClient(client)
	//var msg websocket2.Message
	//_ = connection.ReadJSON(&msg)
	//broker.Channel <- msg
}
