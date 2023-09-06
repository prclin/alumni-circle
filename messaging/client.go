package messaging

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

func (c *Client) Listen() {
	for {
		frame, err := c.ReadFrame()
		if err != nil {
			c.Conn.Close()
			break
		}
		fmt.Printf("%v", frame)
	}
}

func (c *Client) ReadFrame() (*Frame, error) {
	//读取消息
	messageType, message, err := c.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	//目前只支持text message
	if messageType != websocket.TextMessage {
		return nil, errors.New("unsupported message type")
	}

	//解析frame
	return Resolve(message)
}
