package usecase

import (
	"fmt"
	"net"
	"net/http"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/domain"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/repository"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/statics"
)

type RateLimiterUseCase struct {
	Repo        repository.IRateLimiterRepository
	RateLimiter *domain.RateLimiter
}

func NewRateLimiterUseCase(repo repository.IRateLimiterRepository, rl *domain.RateLimiter) *RateLimiterUseCase {
	return &RateLimiterUseCase{
		Repo:        repo,
		RateLimiter: rl,
	}
}

func (useCase *RateLimiterUseCase) Execute(w http.ResponseWriter, r *http.Request) (bool, error) {
	token := r.Header.Get(statics.API_KEY)
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

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false, fmt.Errorf("error parse ip: %v", err)
	}

	isAllowedIP, err := useCase.Repo.AllowIP(
		useCase.RateLimiter.Context,
		clientIP,
		useCase.RateLimiter.IPRateLimiterMaxRequest,
		useCase.RateLimiter.IPRateLimiterBlockTime)

	if err != nil {
		return false, fmt.Errorf("error validating ip: %v", err)
	}

	return isAllowedIP, nil
}
