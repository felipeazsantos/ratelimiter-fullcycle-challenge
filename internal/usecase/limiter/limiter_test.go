package usecase

import (
	"errors"
	"net"
	"net/http/httptest"
	"testing"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/domain"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/mocks"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/statics"
	"github.com/stretchr/testify/assert"
)

func TestExecute_WithTokenAllowed(t *testing.T) {
	mockRepo := new(mocks.MockRateLimiterRepo)
	rl := domain.NewRateLimiter(1, "20s", 1, "20s")
	usecase := NewRateLimiterUseCase(mockRepo, rl)

	token := "abc123"

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(statics.API_KEY, token)
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
	rl := domain.NewRateLimiter(1, "20s", 1, "20s")
	usecase := NewRateLimiterUseCase(mockRepo, rl)

	clientIP := "1.2.3.4:55555"
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = clientIP
	w := httptest.NewRecorder()

	ip, _, _ := net.SplitHostPort(clientIP)
	mockRepo.On("AllowIP", rl.Context, ip, rl.IPRateLimiterMaxRequest, rl.IPRateLimiterBlockTime).
		Return(true, nil)

	allowed, err := usecase.Execute(w, req)
	assert.NoError(t, err)
	assert.True(t, allowed)
	mockRepo.AssertExpectations(t)
}

func TestExecute_withTokenError(t *testing.T) {
	mockRepo := new(mocks.MockRateLimiterRepo)
	rl := domain.NewRateLimiter(1, "20s", 1, "20s")

	usecase := NewRateLimiterUseCase(mockRepo, rl)

	token := "abc123"
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(statics.API_KEY, token)
	w := httptest.NewRecorder()

	mockRepo.On("AllowToken", rl.Context, token, rl.TokenRateLimiterMaxRequest, rl.TokenRateLimiterBlockTime).
		Return(false, errors.New("some error on redis"))

	allowed, err := usecase.Execute(w, req)
	assert.Error(t, err)
	assert.False(t, allowed)

	mockRepo.AssertExpectations(t)
}

func TestExecute_withIPError(t *testing.T) {
	mockRepo := new(mocks.MockRateLimiterRepo)
	rl := domain.NewRateLimiter(1, "20s", 1, "20s")

	usecase := NewRateLimiterUseCase(mockRepo, rl)

	clientIP := "1.2.3.4:55555"
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = clientIP
	w := httptest.NewRecorder()

	ip, _, _ := net.SplitHostPort(clientIP)
	mockRepo.On("AllowIP", rl.Context, ip, rl.IPRateLimiterMaxRequest, rl.IPRateLimiterBlockTime).
		Return(false, errors.New("some error on redis"))

	allowed, err := usecase.Execute(w, req)
	assert.Error(t, err)
	assert.False(t, allowed)

	mockRepo.AssertExpectations(t)
}

func TestExecute_withIPParseError(t *testing.T) {
	mockRepo := new(mocks.MockRateLimiterRepo)
	rl := domain.NewRateLimiter(1, "20s", 1, "20s")

	usecase := NewRateLimiterUseCase(mockRepo, rl)

	clientIP := "1.2.3.4"
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = clientIP
	w := httptest.NewRecorder()

	allowed, err := usecase.Execute(w, req)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error parse ip")
	assert.False(t, allowed)

	mockRepo.AssertExpectations(t)
}
