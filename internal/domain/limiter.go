package domain

import (
	"context"
)

type RateLimiter struct {
	IPRateLimiterMaxRequest    int
	IPRateLimiterBlockTime     int
	TokenRateLimiterMaxRequest int
	TokenRateLimiterBlockTime  int
	Context                    context.Context
}

func NewRateLimiter(ipRateLimiterMaxRequest,
	ipRateLimiterBlockTime,
	tokenRateLimiterMaxRequest,
	tokenRateLimiterBlockTime int) *RateLimiter {
	return &RateLimiter{
		IPRateLimiterMaxRequest:    ipRateLimiterMaxRequest,
		IPRateLimiterBlockTime:     ipRateLimiterBlockTime,
		TokenRateLimiterMaxRequest: tokenRateLimiterMaxRequest,
		TokenRateLimiterBlockTime:  tokenRateLimiterBlockTime,
		Context:                    context.Background(),
	}
}
