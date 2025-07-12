package app

import (
	"net/http"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/getenv"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/domain"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/entrypoint"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/middleware"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/repository"
	usecase "github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/usecase/limiter"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func StartDependencies(router chi.Router, rdb *redis.Client) error {
	config := getenv.AppConfig
	rl := domain.NewRateLimiter(
		config.IPRateLimiterMaxRequest,
		config.IPRateLimiterBlockTime,
		config.TokenRateLimiterMaxRequest,
		config.TokenRateLimiterBlockTime,
	)

	redisRepo := repository.NewRedisRepository(rdb)
	limiterUseCase := usecase.NewRateLimiter(redisRepo, rl)
	limiterHandle := entrypoint.NewRateLimiterHandle()

	router.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimiterMiddleware(limiterUseCase, next)
	})
	router.Get("/", limiterHandle.Handle)

	return nil
}
