package controller

import (
	"context"
	"github.com/prclin/alumni-circle/global"
	"github.com/redis/go-redis/v9"
)

var engine = NewKeyEventEngine()

type KeyEventHandler func(msg *redis.Message)

type KeyEventEngine struct {
	handlers map[string]KeyEventHandler
}

func NewKeyEventEngine() *KeyEventEngine {
	return &KeyEventEngine{handlers: make(map[string]KeyEventHandler)}
}

// Handle 添加处理key过期事件函数
func (e *KeyEventEngine) Handle(key string, handler KeyEventHandler) {
	e.handlers[key] = handler
}

// handle 处理事件
func (e *KeyEventEngine) handle(msg *redis.Message) {
	e.handlers[msg.Payload](msg)
}

func SubscribeKeyEvent(engine *KeyEventEngine) {
	//监听键空间事件
	subscribe := global.RedisClient.Subscribe(context.Background(), "__keyevent@0__:expired")
	defer subscribe.Close()
	channel := subscribe.Channel()
	for msg := range channel {
		go engine.handle(msg)
	}
}

func init() {
	engine.Handle("flush_likes", HandleFlushLikes)
	SubscribeKeyEvent(engine)
}

func HandleFlushLikes(msg *redis.Message) {

}
