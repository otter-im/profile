package config

import (
	"github.com/JeremyLoy/config"
	"log"
	"os"
	"sync"
	"time"
)

var (
	once sync.Once
	c    *AppConfig
)

type AppConfig struct {
	ServiceEnv              string        `config:"SERVICE_ENV"`
	ServiceAddress          string        `config:"SERVICE_ADDRESS"`
	DatabaseURL             string        `config:"DATABASE_URL"`
	RedisNodes              string        `config:"REDIS_NODES"`
	RedisDB                 int           `config:"REDIS_DB"`
	RedisUsername           string        `config:"REDIS_USERNAME"`
	RedisPassword           string        `config:"REDIS_PASSWORD"`
	RedisConnectTimeout     time.Duration `config:"REDIS_CONNECT_TIMEOUT"`
	RedisReadTimeout        time.Duration `config:"REDIS_READ_TIMEOUT"`
	RedisWriteTimeout       time.Duration `config:"REDIS_WRITE_TIMEOUT"`
	RedisIdleTimeout        time.Duration `config:"REDIS_IDLE_TIMEOUT"`
	RedisPoolTimeout        time.Duration `config:"REDIS_POOL_TIMEOUT"`
	RedisPoolFIFO           bool          `config:"REDIS_POOL_FIFO"`
	RedisPoolSIze           int           `config:"REDIS_POOL_SIZE"`
	RedisMinIdleConnections int           `config:"REDIS_MIN_IDLE_CONNECTIONS"`
	RedisTinyLFUSize        int           `config:"REDIS_TINYLFU_SIZE"`
	RedisTinyLFUTtl         time.Duration `config:"REDIS_TINYLFU_TTL"`
	RedisTinyLFUStats       bool          `config:"REDIS_TINYLFU_STATS"`
}

func Config() *AppConfig {
	once.Do(func() {
		conf := &AppConfig{}
		err := config.FromEnv().From("configs/local.env").To(conf)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		c = conf
	})
	return c
}
