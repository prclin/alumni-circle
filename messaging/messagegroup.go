package messaging

type MessageGroup struct {
	prefix string
	broker *Broker
}

type Handler func(ctx *Context)

func (mg *MessageGroup) Handle(destination string, handlers ...Handler) {

}

func (mg *MessageGroup) Subscribe(destination string, handlers ...Handler) {

}

func (mg *MessageGroup) Group(prefix string) *MessageGroup {
	return &MessageGroup{
		prefix: mg.prefix + prefix,
		broker: mg.broker,
	}
}
