package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisLimiterRepository struct {
	RClient *redis.Client
}

func NewRedisRepository(rdb *redis.Client) IRateLimiterRepository {
	return &redisLimiterRepository{}
}

func (r *redisLimiterRepository) AllowToken(ctx context.Context, key string, tokenMaxRequest, tokenBlockTime int) (bool, error) {
	pipe := r.RClient.TxPipeline()

	incr := pipe.Incr(ctx, key)

	pipe.Expire(ctx, key, time.Duration(tokenBlockTime))
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}

	return incr.Val() <= int64(tokenMaxRequest), nil
}

func (r *redisLimiterRepository) AllowIP(ctx context.Context, key string, ipMaxRequest, ipBlockTime int) (bool, error) {
	pipe := r.RClient.TxPipeline()

	incr := pipe.Incr(ctx, key)

	pipe.Expire(ctx, key, time.Duration(ipBlockTime))
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}

	return incr.Val() <= int64(ipMaxRequest), nil
}
