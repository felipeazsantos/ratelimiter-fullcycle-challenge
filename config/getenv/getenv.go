package getenv

import "github.com/spf13/viper"


type AppConfig struct {
	IPRateLimiterMaxRequest int `mapstructure:"RATE_LIMITER_IP_MAX_REQUESTS"`
	IPRateLimiterBlockTime int `mapstructure:"RATE_LIMITER_IP_BLOCK_TIME"`
	TokenRateLimiterMaxRequest int `mapstructure:"RATE_LIMITER_TOKEN_MAX_REQUESTS"`
	TokenRateLimiterBlockTime int `mapstructure:"RATE_LIMITER_TOKEN_BLOCK_TIME"`
	RedisAddr string `mapstructure:"REDIS_ADDR"`
}

var (
	IPRateLimiterMaxRequest, IPRateLimiterBlockTime, TokenRateLimiterMaxRequest, TokenRateLimiterBlockTime int
	RedisAddr string
)


func LoadConfig(path string) (*AppConfig, error) {
	config := &AppConfig{}

	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}