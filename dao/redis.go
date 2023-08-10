package dao

import (
	"context"
	. "github.com/prclin/alumni-circle/global"
	"time"
)

// SetString 设置字符串
func SetString(key string, value string, expiration time.Duration) error {
	return RedisClient.Set(context.Background(), key, value, expiration).Err()
}

// DeleteKey 删除key
func DeleteKey(key string) {
	RedisClient.Del(context.Background(), key)
}
