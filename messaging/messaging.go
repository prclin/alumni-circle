package messaging

import (
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
)

type StompBroker struct {
	MessageGroup
	Upgrader
	inbound  InboundChannel
	outbound OutboundChannel
	//最大连接数
	MaxConnections *atomic.Uint32
	MaxWorkers     *atomic.Uint32
}

// NewStompBroker 构造器
func NewStompBroker() *StompBroker {
	mc := &atomic.Uint32{}
	mc.Store(8192)
	return &StompBroker{
		MaxConnections: mc,
		inbound: InboundChannel{
			frames: make(chan *Frame, 12),
		},
	}
}

// Run 启动消息代理
func (sb *StompBroker) Run() error {
	go sb.inbound.Process()
	go sb.outbound.Process()
	return nil
}

// ServeOver 与client建立连接，建立成功后会阻塞在这，知道发生错误或者连接中断
func (sb *StompBroker) ServeOver(conn *websocket.Conn) error {
	//减少连接数
	for maxConns := sb.MaxConnections.Load(); maxConns > 0 && !sb.MaxConnections.CompareAndSwap(maxConns, maxConns-1); maxConns = sb.MaxConnections.Load() {
		time.Sleep(time.Millisecond)
	}
	//增加连接数
	defer func() {
		for maxConns := sb.MaxConnections.Load(); !sb.MaxConnections.CompareAndSwap(maxConns, maxConns+1); maxConns = sb.MaxConnections.Load() {
			time.Sleep(time.Millisecond)
		}
	}()
	//升级为stomp协议
	client, err := sb.Upgrade(conn)
	if err != nil {
		return err
	}
	//读取消息
	for {
		frame, err1 := client.ReadFrame()
		if err1 != nil {
			//写回错误
			client.WriteFrame(NewFrame(ERROR, map[string]string{"message": err1.Error()}, nil))
			//关闭连接
			client.Conn.Close()
			return err1
		}
		//传递到inbound channel
		sb.inbound.frames <- frame
	}
}
