package messaging

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

// Upgrader 协议升级器
//
// 从http连接升级到stomp连接
type Upgrader struct {
	wsUpgrader websocket.Upgrader
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	//建立websocket连接
	wsConn, err := u.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	conn := &Conn{
		conn: wsConn,
	}

	//读取CONNECT/STOMP帧
	frame, err := conn.ReadFrame()
	if err != nil {
		return nil, err
	}
	//不是CONNECT帧
	if frame.Command != CONNECT && frame.Command != STOMP {
		ef := NewFrame(ERROR, map[string]string{"message": "conn sent other frame before CONNECT"}, []byte(frame.String()))
		conn.WriteFrame(ef)
		return nil, errors.New("not allowed frame before CONNECT")
	}

	v, ok := frame.Headers["accept-version"]

	if !ok {
		ef := NewFrame(ERROR, map[string]string{"message": "CONNECT frame must has an accept-version header"}, []byte(frame.String()))
		conn.WriteFrame(ef)
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
		conn.WriteFrame(ef)
		return nil, errors.New("unsupported protocol version")
	}

	//连接成功
	rf := NewFrame(CONNECTED, map[string]string{"version": "1.1", "heart-beat": "0,0"}, nil)
	conn.WriteFrame(rf)

	return conn, nil
}
