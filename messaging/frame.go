package messaging

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// commands
const (
	MESSAGE Command = iota
	RECEIPT
	ERROR

	CONNECT
	STOMP
	CONNECTED
	SEND
	SUBSCRIBE
	UNSUBSCRIBE
	ACK
	NACK
	BEGIN
	COMMIT
	ABORT
	DISCONNECT
)

var (
	//所有支持的command
	commands       = []string{"MESSAGE", "RECEIPT", "ERROR", "CONNECT", "STOMP", "CONNECTED", "SEND", "SUBSCRIBE", "UNSUBSCRIBE", "ACK", "NACK", "BEGIN", "COMMIT", "ABORT", "DISCONNECT"}
	serverCommands = commands[:2]
	clientCommands = commands[3:]
)

// 确保Command实现了Stringer
var _ fmt.Stringer = Command(0)

type Command int8

// CommandOf 通过字符串值，获取command枚举
func CommandOf(command string) (Command, error) {
	for i, v := range commands {
		if v == command {
			return Command(i), nil
		}
	}
	return -1, errors.New("unsupported command")
}

func (c Command) String() string {
	return commands[c]
}

// Frame stomp协议帧
//
// Only the SEND, MESSAGE, and ERROR frames can have a body. All other frames MUST NOT have a body.
type Frame struct {
	Command Command
	Headers map[string]string
	Body    []byte
}

func NewFrame(command Command, headers map[string]string, payload []byte) *Frame {
	return &Frame{Command: command, Headers: headers, Body: payload}
}

func (f Frame) String() string {
	s := ""
	//command
	s += f.Command.String() + "\n"
	//headers
	for key, value := range f.Headers {
		s += key + ":" + value + "\n"
	}
	//blank line
	s += "\n"
	//body
	s += string(f.Body) + string([]byte{0x00})
	return s
}

func IsClientCommand(command string) bool {
	for _, clientCommand := range clientCommands {
		if command == clientCommand {
			return true
		}
	}
	return false
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
	f.Command, err = CommandOf(cs)
	if err != nil {
		return nil, err
	}

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

	f.Body = payload[:len(payload)-1]
	return f, nil
}
