package messaging

import (
	"errors"
	"github.com/gorilla/websocket"
	"strings"
)

// OverWebsocket 在websocket上运行stomp子协议
func OverWebsocket(conn *websocket.Conn) error {
	//读取CONNECT/STOMP帧
	messageType, buf, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	//目前只支持text message
	if messageType != websocket.TextMessage {
		return errors.New("unsupported message type")
	}
	//解析frame
	frame, err := Resolve(buf)
	if err != nil {
		return err
	}
	//不是CONNECT帧
	if frame.Command != Connect {
		data := make([]byte, 0)
		data = append(data, Error+"\n"...)
		data = append(data, "message:client sent other frame before CONNECT\n"...)
		data = append(data, "\n"...)
		data = append(data, 0x00)
		conn.WriteMessage(websocket.TextMessage, data)
		return errors.New("not allowed frame before CONNECT")
	}

	v, ok := frame.Headers["accept-version"]

	if !ok {
		data := make([]byte, 0)
		data = append(data, Error+"\n"...)
		data = append(data, "CONNECT frame must has an accept-version header\n"...)
		data = append(data, "\n"...)
		data = append(data, 0x00)
		conn.WriteMessage(websocket.TextMessage, data)
		return errors.New("CONNECT frame must has an accept-version header")
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
		data := make([]byte, 0)
		data = append(data, Error+"\n"...)
		data = append(data, "unsupported protocol version\n"...)
		data = append(data, "\n"...)
		data = append(data, 0x00)
		conn.WriteMessage(websocket.TextMessage, data)
		return errors.New("unsupported protocol version")
	}

	//连接成功
	data := make([]byte, 0)
	data = append(data, Connected+"\n"...)
	data = append(data, "version:1.1\n"...)
	data = append(data, "heart-beat:0,0\n"...)
	data = append(data, "\n"...)
	data = append(data, 0x00)
	conn.WriteMessage(websocket.TextMessage, data)
	return nil
}
