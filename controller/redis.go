package controller

import (
	"context"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/service"
	"github.com/redis/go-redis/v9"
	"time"
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
func (e *KeyEventEngine) Handle(key string, handler KeyEventHandler, expiration time.Duration) {
	_, err := global.RedisClient.Set(context.Background(), key, expiration, expiration).Result()
	if err != nil {
		global.Logger.Fatalln(err)
	}
	e.handlers[key] = handler
}

// handle 处理事件
func (e *KeyEventEngine) handle(msg *redis.Message) {
	handler, ok := e.handlers[msg.Payload]
	if !ok {
		global.Logger.Warn("unhandled event", msg.Payload)
		return
	}
	handler(msg)
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
	engine.Handle("expired_break_likes", HandleBreakLikesExpire, time.Hour)
	go SubscribeKeyEvent(engine)
}

func HandleBreakLikesExpire(msg *redis.Message) {
	service.FlushBreakLikes()
}
