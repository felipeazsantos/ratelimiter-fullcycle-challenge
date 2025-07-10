package entrypoint

import (
	"net/http"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/usecase/limiter"
)


type RateLimiterHandle struct {
	UseCase *limiter.RateLimiterUseCase
}


func NewRateLimiterHandle(usecase *limiter.RateLimiterUseCase) *RateLimiterHandle {
	return &RateLimiterHandle{
		UseCase: usecase,
	}
}

func (handle *RateLimiterHandle) Handle(w http.ResponseWriter, r *http.Request) {
	
}