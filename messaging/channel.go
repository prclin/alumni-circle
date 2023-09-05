package messaging

type InboundChannel struct {
	//消息帧
	fChan chan *Frame
}

type OutboundChannel struct {
	//消息帧
}
