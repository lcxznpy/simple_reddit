package redis

import (
	"context"
	"fmt"
	"goweb/settings"

	"github.com/redis/go-redis/v9"
)

// 声明一个全局的redis db变量
var (
	client *redis.Client
	Nil    = redis.Nil
)

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host,
			cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize, // use default DB
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

// 关闭redis
func Close() {
	_ = client.Close()
}
