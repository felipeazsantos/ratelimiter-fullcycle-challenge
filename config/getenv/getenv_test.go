package getenv

import (
	"fmt"
	"os"
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

func TestInitConfig_FilePathError(t *testing.T) {
	err := InitConfig("some-invalid-path")
	assert.Error(t, err)
}

func TestInitConfig_UnmarshallError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "invalidenv")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("RATE_LIMITER_IP_MAX_REQUESTS=notanumber\n")
	assert.NoError(t, err)
	tmpFile.Close()

	err = InitConfig(tmpFile.Name())
	assert.Error(t, err)
}