package messaging

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var (
	_ Resolver = &TextMessageResolver{}
)

// Resolver 帧解析器
type Resolver interface {
	//Resolve 解析stomp帧
	Resolve(buf []byte) (*Frame, error)
}

// TextMessageResolver 简单消息解析器
type TextMessageResolver struct {
}

func (s *TextMessageResolver) Resolve(buf []byte) (*Frame, error) {
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
