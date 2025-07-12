package usecase

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/domain"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestExecute_WithTokenAllowed(t *testing.T) {
	mockRepo := new(mocks.MockRateLimiterRepo)
	rl := &domain.RateLimiter{
		Context:                    context.Background(),
		TokenRateLimiterMaxRequest: 10,
		TokenRateLimiterBlockTime:  60,
	}
	usecase := NewRateLimiterUseCase(mockRepo, rl)

	token := "abc123"

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("token", token)
	w := httptest.NewRecorder()

	mockRepo.On("AllowToken", rl.Context, token, rl.TokenRateLimiterMaxRequest, rl.TokenRateLimiterBlockTime).
		Return(true, nil)

	allowed, err := usecase.Execute(w, req)
	assert.NoError(t, err)
	assert.True(t, allowed)
	mockRepo.AssertExpectations(t)
}

func TestExecute_withIPAllowed(t *testing.T) {
	mockRepo := new(mocks.MockRateLimiterRepo)
	rl := &domain.RateLimiter{
		Context:                 context.Background(),
		IPRateLimiterMaxRequest: 5,
		IPRateLimiterBlockTime:  30,
	}
	usecase := NewRateLimiterUseCase(mockRepo, rl)

	clientIP := "1.2.3.4"
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = clientIP
	w := httptest.NewRecorder()

	mockRepo.On("AllowIP", rl.Context, clientIP, rl.IPRateLimiterMaxRequest, rl.IPRateLimiterBlockTime).
		Return(true, nil)

	allowed, err := usecase.Execute(w, req)
	assert.NoError(t, err)
	assert.True(t, allowed)
	mockRepo.AssertExpectations(t)
}
