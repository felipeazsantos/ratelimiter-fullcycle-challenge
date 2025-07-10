package repository

import "github.com/redis/go-redis/v9"

type redisLimiter struct{}

func NewRedisRepository(rdb *redis.Client) IRepositoryLimiter {
	return &redisLimiter{}
}

func (r *redisLimiter) Save() error {
	return nil
}

func (r *redisLimiter) GetByIp() error {
	return nil
}

func (r *redisLimiter) GetByToken() error {
	return nil
}
