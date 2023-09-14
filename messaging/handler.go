package messaging

import "strings"

var (
	//确保多态
	_ MessageHandler = &MethodMessageHandler{}
	_ MessageHandler = &BrokerMessageHandler{}
)

// MessageHandler 消息处理器
type MessageHandler interface {
	HandleMessage(*Context)
}

// MethodMessageHandler 应用消息处理器
type MethodMessageHandler struct {
}

func (mmh *MethodMessageHandler) HandleMessage(context *Context) {
	methodHandler := mmh.getMethodHandler(context)
	if methodHandler == nil {
		return
	}
	methodHandler(context)
}

func (mmh *MethodMessageHandler) getMethodHandler(context *Context) MethodHandler {
	destArr := strings.Split(context.Frame.Destination(), "/")
	var handler MethodHandler
	for key, value := range context.broker.sendMap {
		keyArr := strings.Split(key, "/")
		if len(destArr) != len(keyArr) {
			continue
		}
		match := true
		for i := 0; i < len(keyArr); i++ {
			if strings.HasPrefix(keyArr[i], ":") {
				context.Params[strings.TrimPrefix(keyArr[i], ":")] = destArr[i]
				continue
			}
			if keyArr[i] != destArr[i] {
				match = false
				break
			}
		}
		if match {
			handler = value
			break
		}
	}

	for key, value := range context.broker.subscribeMap {
		keyArr := strings.Split(key, "/")
		if len(destArr) != len(keyArr) {
			continue
		}
		match := true
		for i := 0; i < len(keyArr); i++ {
			if strings.HasPrefix(keyArr[i], ":") {
				context.Params[strings.TrimPrefix(keyArr[i], ":")] = destArr[i]
				continue
			}
			if keyArr[i] != destArr[i] {
				match = false
				break
			}
		}
		if match {
			handler = value
			break
		}
	}
	return handler
}

// BrokerMessageHandler 中介消息处理器
//
// 处理ClientInboundChannel和BrokerChannel中的消息
type BrokerMessageHandler struct {
	outboundChannel MessageChannel
}

func (bmh *BrokerMessageHandler) HandleMessage(context *Context) {
	context.Frame.Command = MESSAGE
	bmh.outboundChannel.Send(context)
}

type MethodHandler func(ctx *Context)
