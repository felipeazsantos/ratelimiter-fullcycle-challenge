package app

import (
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/getenv"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/redisdb"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/entrypoint"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/repository"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/usecase/limiter"
	"github.com/go-chi/chi/v5"
)

type Application struct {
	Config *getenv.AppConfig
}

func StartDependencies(router chi.Router) error {
	rdb := redisdb.NewRedisClient()
	redisRepo := repository.NewRedisRepository(rdb)
	limiterUseCase := limiter.NewRateLimiter(redisRepo)
	limiterHandle := entrypoint.NewRateLimiterHandle(limiterUseCase)

	router.Get("/", limiterHandle.Handle)

	return nil
}
