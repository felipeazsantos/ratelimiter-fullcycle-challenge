package domain

import (
	"context"
)

type RateLimiter struct {
	IPRateLimiterMaxRequest    int
	IPRateLimiterBlockTime     string
	TokenRateLimiterMaxRequest int
	TokenRateLimiterBlockTime  string
	Context                    context.Context
}

func NewRateLimiter(ipRateLimiterMaxRequest int,
	ipRateLimiterBlockTime string,
	tokenRateLimiterMaxRequest int,
	tokenRateLimiterBlockTime string) *RateLimiter {
	return &RateLimiter{
		IPRateLimiterMaxRequest:    ipRateLimiterMaxRequest,
		IPRateLimiterBlockTime:     ipRateLimiterBlockTime,
		TokenRateLimiterMaxRequest: tokenRateLimiterMaxRequest,
		TokenRateLimiterBlockTime:  tokenRateLimiterBlockTime,
		Context:                    context.Background(),
	}
}
