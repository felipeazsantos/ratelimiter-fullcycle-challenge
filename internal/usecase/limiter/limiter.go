package limiter

import (
	"context"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/repository"
)

type RateLimiterUseCase struct {
	RedisRepo repository.IRepositoryLimiter
}

func NewRateLimiter(redisRepo repository.IRepositoryLimiter) *RateLimiterUseCase {
	return &RateLimiterUseCase{
		RedisRepo: redisRepo,
	}
}

func (rl *RateLimiterUseCase) Execute(ctx context.Context) {

}
