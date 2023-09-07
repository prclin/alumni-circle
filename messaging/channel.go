package messaging

import "fmt"

type InboundChannel struct {
	frames chan *Frame
}

func (ic *InboundChannel) Process() {
	for {
		frame := <-ic.frames
		fmt.Println(frame)
	}
}

type OutboundChannel struct {
	frames chan *Frame
}

func (oc *OutboundChannel) Process() {
	for {
		frame := <-oc.frames
		fmt.Println(frame)
	}
}
