package getenv

import (
	"github.com/spf13/viper"
)

type appConfig struct {
	IPRateLimiterMaxRequest    int    `mapstructure:"RATE_LIMITER_IP_MAX_REQUESTS"`
	IPRateLimiterBlockTime     string    `mapstructure:"RATE_LIMITER_IP_BLOCK_TIME"`
	TokenRateLimiterMaxRequest int    `mapstructure:"RATE_LIMITER_TOKEN_MAX_REQUESTS"`
	TokenRateLimiterBlockTime  string    `mapstructure:"RATE_LIMITER_TOKEN_BLOCK_TIME"`
	RedisAddr                  string `mapstructure:"REDIS_ADDR"`
}

var AppConfig *appConfig

func NewAppConfig(ipMax, tokenMax int, ipBlock, tokenBlock, redisAddr string) *appConfig {
	return &appConfig{
		IPRateLimiterMaxRequest:    ipMax,
		IPRateLimiterBlockTime:     ipBlock,
		TokenRateLimiterMaxRequest: tokenMax,
		TokenRateLimiterBlockTime:  tokenBlock,
		RedisAddr:                  redisAddr,
	}
}

func InitConfig(path string) error {
	config := &appConfig{}

	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	AppConfig = config

	return nil
}
