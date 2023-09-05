package messaging

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Frame struct {
	Command string
	Headers map[string]string
	Payload []byte
}

// commands
const (
	Message     = "MESSAGE"
	Receipt     = "RECEIPT"
	Error       = "ERROR"
	Connect     = "CONNECT"
	Connected   = "CONNECTED"
	Send        = "SEND"
	Subscribe   = "SUBSCRIBE"
	Unsubscribe = "UNSUBSCRIBE"
	Ack         = "ACK"
	Nack        = "NACK"
	Begin       = "BEGIN"
	Commit      = "COMMIT"
	Abort       = "ABORT"
	Disconnect  = "DISCONNECT"
)

var (
	ClientCommands = map[string]struct{}{
		Connect:     {},
		Send:        {},
		Subscribe:   {},
		Unsubscribe: {},
		Ack:         {},
		Nack:        {},
		Begin:       {},
		Commit:      {},
		Abort:       {},
		Disconnect:  {},
	}
	ServerCommands = map[string]struct{}{
		Connected: {},
		Message:   {},
		Receipt:   {},
		Error:     {},
	}
)

func IsClientCommand(command string) bool {
	_, ok := ClientCommands[command]
	return ok
}

// Resolve 解析frame
func Resolve(buf []byte) (*Frame, error) {
	//获取reader
	reader := bufio.NewReader(bytes.NewReader(buf))

	f := &Frame{}
	//读取command
	cs, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	cs = strings.TrimSuffix(cs, "\n")
	//是否是client frame
	if !IsClientCommand(cs) {
		fmt.Println(cs)
		return nil, errors.New("not client command")
	}
	f.Command = cs

	f.Headers = make(map[string]string, 2)
	//读取消息头
	for {
		hs, err1 := reader.ReadString('\n')
		if err1 != nil {
			return nil, err1
		}
		hs = strings.TrimSuffix(hs, "\n")
		//读到空行结束
		if hs == "" {
			break
		}
		kv := strings.SplitN(hs, ":", 2)
		//header格式是否合法
		if len(kv) != 2 {
			return nil, errors.New("wrong header format")
		}
		//如果有重复的header，保留第一个
		if _, ok := f.Headers[kv[0]]; !ok {
			f.Headers[kv[0]] = kv[1]
		}
	}

	//读取payload,目前body只支持json
	payload, err := reader.ReadBytes(0x00)
	if err != nil {
		return nil, err
	}

	f.Payload = payload[:len(payload)-1]
	return f, nil
}
