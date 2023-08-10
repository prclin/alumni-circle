package config

import "github.com/redis/go-redis/v9"

// Redis redis配置，包含单客户端，集群客户端等配置
type Redis struct {
	//单客户端配置，使用go-redis的Options
	Options *redis.Options
}

var DefaultRedis = &Redis{
	Options: &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	},
}
