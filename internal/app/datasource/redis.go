package datasource

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/otter-im/profile/internal/config"
	"sync"
	"time"
)

var (
	redisRingOnce  sync.Once
	redisRing      *redis.Ring
	redisCacheOnce sync.Once
	redisCache     *cache.Cache
)

func RedisRing() *redis.Ring {
	redisRingOnce.Do(func() {
		options := &redis.RingOptions{
			Addrs:    config.RedisNodes(),
			Password: config.RedisPassword(),
			DB:       config.RedisDB(),
		}

		redisRing = redis.NewRing(options)
	})
	return redisRing
}

func RedisCache() *cache.Cache {
	redisCacheOnce.Do(func() {
		options := &cache.Options{
			Redis: RedisRing(),
		}
		redisCache = cache.New(options)
	})
	return redisCache
}

func CheckRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if cmd := RedisRing().Ping(ctx); cmd.Err() != nil {
		return fmt.Errorf("redis connection failure: %v", cmd.Err())
	}
	return nil
}
