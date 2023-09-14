package messaging

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// StompBroker 消息代理
//
// 数据流向
//
/*
InboundChannel	---->  MethodMessageHandler

		|       ---->  BrokerMessageHandler	<----       |

OutboundChannel	<----         |
*/
type StompBroker struct {
	MessageGroup
	Upgrader
	inbound       MessageChannel
	outbound      MessageChannel
	brokerHandler MessageHandler
	methodHandler MessageHandler
	//应用消息前缀
	AppDestinationPrefix string
	//代理消息前缀
	BrokerDestinationPrefix string
	//锁
	lock sync.Mutex
	//订阅
	subscriptions []*Subscription
	//send帧处理函数
	sendMap map[string]MethodHandler
	//subscribe帧处理函数
	subscribeMap map[string]MethodHandler
}

func NewStompBroker() *StompBroker {
	outbound := &ClientOutBoundChannel{}
	brokerHandler := &BrokerMessageHandler{outboundChannel: outbound}
	methodHandler := &MethodMessageHandler{}
	inbound := &ClientInboundChannel{}
	inbound.AddValves(&SubscribeValve{methodHandler: methodHandler}, &SendValve{methodHandler: methodHandler, brokerHandler: brokerHandler})
	broker := &StompBroker{
		MessageGroup: MessageGroup{
			prefix: "",
		},
		AppDestinationPrefix:    "/app",
		BrokerDestinationPrefix: "/topic",
		Upgrader: Upgrader{
			wsUpgrader: websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
		outbound:      outbound,
		brokerHandler: brokerHandler,
		methodHandler: methodHandler,
		inbound:       inbound,
		subscriptions: make([]*Subscription, 100),
		sendMap:       make(map[string]MethodHandler),
		subscribeMap:  make(map[string]MethodHandler),
	}
	broker.MessageGroup.broker = broker
	return broker
}

// ServeOverHttp 与client建立连接，建立成功后会阻塞在这，直到发生错误或者连接中断
func (sb *StompBroker) ServeOverHttp(w http.ResponseWriter, r *http.Request) error {
	//升级为stomp协议
	conn, err := sb.Upgrade(w, r)
	if err != nil {
		return err
	}
	//监听消息
	for {
		frame, err := conn.ReadFrame()
		if err != nil {
			conn.Close()
			return err
		}
		sb.inbound.Send(&Context{broker: sb, Frame: frame, Conn: conn, Params: make(map[string]string, 2)})
	}
}
func (sb *StompBroker) addSubscription(subscription *Subscription) {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	sb.subscriptions = append(sb.subscriptions, subscription)
}

type Subscription struct {
	Id          string
	Destination string
	Conn        *Conn
}

type MessageGroup struct {
	prefix string
	broker *StompBroker
}

func (mg *MessageGroup) Handle(destination string, handler MethodHandler) {
	mg.broker.sendMap[mg.broker.AppDestinationPrefix+mg.prefix+destination] = handler
}

func (mg *MessageGroup) Subscribe(destination string, handler MethodHandler) {
	mg.broker.subscribeMap[mg.broker.BrokerDestinationPrefix+mg.prefix+destination] = handler
}

func (mg *MessageGroup) Group(prefix string) *MessageGroup {
	return &MessageGroup{
		prefix: mg.prefix + prefix,
		broker: mg.broker,
	}
}
