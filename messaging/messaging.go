package messaging

type Broker struct {
	MessageGroup
	Upgrader Upgrader
	inbound  InboundChannel
	outbound OutboundChannel
}

func (b *Broker) Run() error {
	go b.outbound.Process()
	go b.inbound.Process()
	return nil
}
