package getenv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	if err := InitConfig("../../.env"); err != nil {
		assert.Fail(t, fmt.Sprintf("failure to init env config for test purpose: %v", err.Error()))
		return
	}

	assert.GreaterOrEqual(t, AppConfig.IPRateLimiterMaxRequest, 1)
	assert.GreaterOrEqual(t, AppConfig.TokenRateLimiterMaxRequest, 1)
	assert.NotEmpty(t, AppConfig.IPRateLimiterBlockTime)
	assert.NotEmpty(t, AppConfig.TokenRateLimiterBlockTime)
	assert.NotEmpty(t, AppConfig.RedisAddr)
}
