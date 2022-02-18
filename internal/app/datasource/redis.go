package datasource

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"github.com/otter-im/profile/internal/app/config"
	"strings"
	"sync"
	"time"
)

var (
	redisRing  *redis.Ring
	redisCache *cache.Cache

	redisRingOnce  sync.Once
	redisCacheOnce sync.Once
)

func RedisRing() *redis.Ring {
	redisRingOnce.Do(func() {
		conf := config.Config()

		opts := &redis.RingOptions{
			Addrs:              parseNodes(conf.RedisNodes),
			Username:           conf.RedisUsername,
			Password:           conf.RedisPassword,
			DB:                 conf.RedisDB,
			DialTimeout:        conf.RedisConnectTimeout,
			ReadTimeout:        conf.RedisReadTimeout,
			WriteTimeout:       conf.RedisWriteTimeout,
			PoolFIFO:           conf.RedisPoolFIFO,
			PoolSize:           conf.RedisPoolSIze,
			MinIdleConns:       conf.RedisMinIdleConnections,
			MaxConnAge:         0,
			PoolTimeout:        conf.RedisPoolTimeout,
			IdleTimeout:        conf.RedisIdleTimeout,
			IdleCheckFrequency: 0,
			// TODO: TLS config
		}
		redisRing = redis.NewRing(opts)

		_ = redisRing.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
			shard.AddHook(redisotel.TracingHook{})
			return nil
		})
	})
	return redisRing
}

func RedisCache() *cache.Cache {
	redisCacheOnce.Do(func() {
		conf := config.Config()
		redisCache = cache.New(&cache.Options{
			Redis:        RedisRing(),
			LocalCache:   cache.NewTinyLFU(conf.RedisTinyLFUSize, conf.RedisTinyLFUTtl),
			StatsEnabled: conf.RedisTinyLFUStats,
		})
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

// REDIS_NODES="prod=localhost:6379,dev=localhost:6379"
func parseNodes(v string) map[string]string {
	result := make(map[string]string)
	split := strings.Split(v, ",")
	for _, v := range split {
		nameSplit := strings.Split(v, "=")
		if len(nameSplit) == 1 {
			result[nameSplit[0]] = nameSplit[0]
		} else {
			result[nameSplit[0]] = nameSplit[1]
		}
	}
	return result
}
