package repository

import "context"

type IRateLimiterRepository interface {
	AllowToken(ctx context.Context, key string, tokenMaxRequest, tokenBlockTime int) (bool, error)
	AllowIP(ctx context.Context, key string, ipMaxRequest, ipBlockTime int)  (bool, error)
}
