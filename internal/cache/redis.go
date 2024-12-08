package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"jianghai-hu/wallet-service/internal/config"
)

var globalRedisClient *redis.Client

func InitRedis(ctx context.Context) {
	globalRedisClient = redis.NewClient(config.DefaultRedisConfig)
	_, err := globalRedisClient.Ping(ctx).Result()
	if err != nil {
		glog.FatalContextf(ctx, "redis init failed: %v", err)
	}
}

func RedisClient(ctx context.Context) *redis.Client {
	if globalRedisClient == nil {
		glog.FatalContext(ctx, "redis client is nil")
	}
	return globalRedisClient
}
