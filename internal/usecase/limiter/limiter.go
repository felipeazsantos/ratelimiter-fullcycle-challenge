package usecase

import (
	"fmt"
	"net/http"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/domain"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/repository"
)

type RateLimiterUseCase struct {
	Repo repository.IRateLimiterRepository
	RateLimiter *domain.RateLimiter
}

func NewRateLimiter(redisRepo repository.IRateLimiterRepository, rl *domain.RateLimiter) *RateLimiterUseCase {
	return &RateLimiterUseCase{
		Repo: redisRepo,
		RateLimiter: rl,
	}
}

func (useCase *RateLimiterUseCase) Execute(w http.ResponseWriter, r *http.Request) (bool, error) {
	token := r.Header.Get("token")
	if token != "" {
		isAllowedToken, err := useCase.Repo.AllowToken(
			useCase.RateLimiter.Context, 
			token, 
			useCase.RateLimiter.TokenRateLimiterMaxRequest, 
			useCase.RateLimiter.TokenRateLimiterBlockTime)
		
		if err != nil {
			return false, fmt.Errorf("error validating token: %v", err)
		}

		return isAllowedToken, nil
	}

	clientIP := r.RemoteAddr
	isAllowedIP, err := useCase.Repo.AllowIP(
		useCase.RateLimiter.Context, 
		clientIP, 
		useCase.RateLimiter.IPRateLimiterMaxRequest, 
		useCase.RateLimiter.IPRateLimiterBlockTime)

	if err != nil {
		return false, fmt.Errorf("error validating token: %v", err)
	}

	return isAllowedIP, nil
}
