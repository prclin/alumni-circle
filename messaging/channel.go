package messaging

type InboundChannel struct {
	frames chan *Frame
}

func (ic *InboundChannel) Process() {
	for {
		frame := <-ic.frames
		switch frame.Command {
		case SEND:
		}
	}
}

type OutboundChannel struct {
	frames chan *Frame
}

func (oc *OutboundChannel) Process() {
	for {
		frame := <-oc.frames
		switch frame.Command {
		case SEND:
		}
	}
}
