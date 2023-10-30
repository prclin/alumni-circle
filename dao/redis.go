package dao

import (
	"context"
	. "github.com/prclin/alumni-circle/global"
	"time"
)

func HGetAll(key string) (map[string]string, error) {
	return RedisClient.HGetAll(context.Background(), key).Result()
}

func HIncrBy(key string, field string, increment int64) (int64, error) {
	return RedisClient.HIncrBy(context.Background(), key, field, increment).Result()
}

// HSet 设置hash
func HSet(key string, values ...any) (int64, error) {
	return RedisClient.HSet(context.Background(), key, values).Result()
}

// SetString 设置字符串
func SetString(key string, value string, expiration time.Duration) error {
	return RedisClient.Set(context.Background(), key, value, expiration).Err()
}

// DeleteKey 删除key
func DeleteKey(key string) (int64, error) {
	return RedisClient.Del(context.Background(), key).Result()
}

// GetString 获取字符串
func GetString(key string) (string, error) {
	cmd := RedisClient.Get(context.Background(), key)
	return cmd.Result()
}

// SetBit 设置bitmap指定bit
func SetBit(key string, offset int64, value int) (int64, error) {
	cmd := RedisClient.SetBit(context.Background(), key, offset, value)
	return cmd.Result()
}

// GetBit 获取bitmap指定bit
func GetBit(key string, offset int64) (int64, error) {
	cmd := RedisClient.GetBit(context.Background(), key, offset)
	return cmd.Result()
}
