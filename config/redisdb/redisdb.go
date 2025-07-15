package redisdb

import (
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/getenv"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     getenv.AppConfig.RedisAddr,
		Password: "",
		DB:       0,
	})

	return rdb
}
