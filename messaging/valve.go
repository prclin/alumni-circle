package messaging

import "strings"

// ChannelValve 管道阀
type ChannelValve interface {
	//Valve 处理函数
	Valve(*Context) error
	//SetNextValve 设置下一个valve
	SetNextValve(nextValve ChannelValve)
	//Next 返回下一个处理器
	Next() ChannelValve
}

type StandardValve struct {
	nextValve ChannelValve
}

func (sv *StandardValve) SetNextValve(nextValve ChannelValve) {
	sv.nextValve = nextValve
}

func (sv *StandardValve) Next() ChannelValve {
	return sv.nextValve
}

var _ ChannelValve = &SendValve{}

// SendValve 处理Send帧
type SendValve struct {
	StandardValve
	//应用消息处理器
	methodHandler MessageHandler
	//代理消息处理器
	brokerHandler MessageHandler
}

func (sv *SendValve) Valve(context *Context) error {
	if context.Frame.Command != SEND {
		return sv.nextValve.Valve(context)
	}
	if strings.HasPrefix(context.Frame.Destination(), context.broker.AppDestinationPrefix) {
		sv.methodHandler.HandleMessage(context)
	} else if strings.HasPrefix(context.Frame.Destination(), context.broker.BrokerDestinationPrefix) {
		sv.brokerHandler.HandleMessage(context)
	}
	return nil
}

type SubscribeValve struct {
	StandardValve
	//应用消息处理器
	methodHandler MessageHandler
}

func (sv SubscribeValve) Valve(context *Context) error {
	if context.Frame.Command != SUBSCRIBE {
		return sv.nextValve.Valve(context)
	}
	destination := context.Frame.Destination()
	//订阅
	if !strings.HasPrefix(destination, context.broker.BrokerDestinationPrefix) {
		return nil
	}
	//添加订阅者
	subscription := &Subscription{Id: context.Frame.Headers["id"], Destination: context.Frame.Destination(), Conn: context.Conn}
	context.broker.addSubscription(subscription)
	//调用subscribe handler
	sv.methodHandler.HandleMessage(context)
	return nil
}
