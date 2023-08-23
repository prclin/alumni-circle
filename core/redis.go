package core

import (
	"context"
	"github.com/prclin/alumni-circle/global"
	"github.com/redis/go-redis/v9"
)

// initRedis 初始化redis连接
func initRedis() {
	//创建client
	client := redis.NewClient(global.Configuration.Redis.Options)
	//测试ping pong
	err := client.Ping(context.Background()).Err()
	if err != nil {
		global.Logger.Fatal(err)
	}
	global.RedisClient = client
}
