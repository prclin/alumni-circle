package messaging

import (
	"errors"
	"fmt"
)

// command常量,不能打乱顺序
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
	// 确保Command实现了Stringer
	_ fmt.Stringer = Command(0)
)

// Command stomp帧的指令
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

func (f *Frame) String() string {
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

func (f *Frame) Destination() string {
	return f.Headers["destination"]
}

func IsClientCommand(command string) bool {
	for _, clientCommand := range clientCommands {
		if command == clientCommand {
			return true
		}
	}
	return false
}
