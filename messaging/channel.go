package messaging

// MessageChannel 消息管道
//
// 对接收到的消息作前置处理
type MessageChannel interface {
	//Send 向管道发送信息
	Send(*Context)
}

var _ MessageChannel = &ClientInboundChannel{}

// ClientInboundChannel 接收消息
type ClientInboundChannel struct {
	//处理阀
	rootValve ChannelValve
}

// AddValves 添加阀处理器
func (cic *ClientInboundChannel) AddValves(valves ...ChannelValve) {
	for _, valve := range valves {
		valve.SetNextValve(cic.rootValve)
		cic.rootValve = valve
	}
}

func (cic *ClientInboundChannel) Send(context *Context) {
	var err error
	//阀处理
	if cic.rootValve != nil {
		err = cic.rootValve.Valve(context)
	}
	if err != nil {
		return
	}
}

// ClientOutBoundChannel 发送消息
type ClientOutBoundChannel struct {
}

func (coc *ClientOutBoundChannel) Send(context *Context) {
	for _, subscription := range context.broker.subscriptions {
		if subscription.Destination == context.Frame.Destination() {
			context.Frame.Headers["subscription"] = subscription.Id
			subscription.Conn.WriteFrame(context.Frame)
		}
	}
}
