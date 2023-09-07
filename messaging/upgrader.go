package messaging

import (
	"errors"
	"github.com/gorilla/websocket"
	"strings"
)

type Upgrader struct {
}

func (u *Upgrader) Upgrade(conn *websocket.Conn) (*Client, error) {
	client := &Client{
		Conn: conn,
	}

	//读取CONNECT/STOMP帧
	frame, err := client.ReadFrame()
	if err != nil {
		return nil, err
	}
	//不是CONNECT帧
	if frame.Command != CONNECT && frame.Command != STOMP {
		ef := NewFrame(ERROR, map[string]string{"message": "client sent other frame before CONNECT"}, []byte(frame.String()))
		conn.WriteMessage(websocket.TextMessage, []byte(ef.String()))
		return nil, errors.New("not allowed frame before CONNECT")
	}

	v, ok := frame.Headers["accept-version"]

	if !ok {
		ef := NewFrame(ERROR, map[string]string{"message": "CONNECT frame must has an accept-version header"}, []byte(frame.String()))
		conn.WriteMessage(websocket.TextMessage, []byte(ef.String()))
		return nil, errors.New("CONNECT frame must has an accept-version header")
	}

	//协议版本协商，目前只支持1.1
	cVersions := strings.Split(v, ",")
	var support bool
	for _, version := range cVersions {
		if version == "1.1" {
			support = true
		}
	}

	if !support {
		ef := NewFrame(ERROR, map[string]string{"message": "unsupported protocol version"}, []byte(frame.String()))
		conn.WriteMessage(websocket.TextMessage, []byte(ef.String()))
		return nil, errors.New("unsupported protocol version")
	}

	//连接成功
	rf := NewFrame(CONNECTED, map[string]string{"version": "1.1", "heart-beat": "0,0"}, nil)
	conn.WriteMessage(websocket.TextMessage, []byte(rf.String()))

	return client, nil
}
