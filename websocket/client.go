package websocket

import "github.com/gorilla/websocket"

type Client struct {
	Id   uint64
	Conn *websocket.Conn
}

func NewClient(id uint64, conn *websocket.Conn) *Client {
	return &Client{Id: id, Conn: conn}
}
