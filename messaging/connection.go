package messaging

import (
	"errors"
	"github.com/gorilla/websocket"
)

// Conn stomp连接
type Conn struct {
	conn *websocket.Conn
}

// ReadFrame 读取stomp帧
//
// 此方法会阻塞，直到有可读的stomp帧
//
// 当读取消息错误、消息类型不支持或解析帧失败时返回error
func (c *Conn) ReadFrame() (*Frame, error) {
	//读取消息
	messageType, message, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	//目前只支持text message
	if messageType != websocket.TextMessage {
		return nil, errors.New("unsupported message type")
	}

	//解析frame
	resolver := &TextMessageResolver{}
	return resolver.Resolve(message)
}

// WriteFrame 写stomp帧
//
// 目前只支持写文本消息类型
func (c *Conn) WriteFrame(frame *Frame) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(frame.String()))
}

// Close 关闭stomp连接
func (c *Conn) Close() error {
	return c.conn.Close()
}
