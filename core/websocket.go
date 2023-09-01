package core

import (
	"github.com/gorilla/websocket"
	"github.com/prclin/alumni-circle/global"
	"net/http"
)

func initWebsocket() {
	wsConfig := global.Configuration.Websocket
	// 创建upGrader实例
	upgrader := &websocket.Upgrader{
		ReadBufferSize:    wsConfig.Upgrader.ReadBufferSize,
		WriteBufferSize:   wsConfig.Upgrader.WriteBufferSize,
		EnableCompression: wsConfig.Upgrader.EnableCompression,
		CheckOrigin:       func(r *http.Request) bool { return true },
	}
	global.WebsocketUpgrader = upgrader
	// TODO: 初始化连接池
}
