package config

import (
	"targetads/internal/http/server"
	"targetads/internal/logger"
	"targetads/internal/storage/aws"
	"targetads/internal/storage/redis"
	"time"

	"github.com/spf13/viper"
)

const (
	logLevelEnv  = "LOG_LEVEL"
	logFormatEnv = "LOG_FORMAT"
	redisUriEnv  = "REDIS_URI"
	awsRegionEnv = "AWS_REGION"
	awsBucketEnv = "AWS_BUCKET_NAME"
	httpPortEnv  = "PORT"
	httpHostEnv  = "HOST"
	httpClearEnv = "CLEAR_LOCAL_CACHE_PERIOD"
)

type Config struct {
	Redis  redis.Config
	AWS    aws.Config
	HTTP   server.Config
	Logger logger.Config
}

func setDefaultConfig() {
	viper.SetDefault(logLevelEnv, "info")
	viper.SetDefault(logFormatEnv, "json")
	viper.SetDefault(httpPortEnv, 8089)
	viper.SetDefault(httpHostEnv, "localhost")
	viper.SetDefault(httpClearEnv, 10*time.Minute)
}

func Parse() (*Config, error) {
	var cfg Config

	setDefaultConfig()
	viper.AutomaticEnv()
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Level = viper.GetString(logLevelEnv)
	cfg.Logger.Format = viper.GetString(logFormatEnv)
	cfg.Redis.URI = viper.GetString(redisUriEnv)
	cfg.AWS.Region = viper.GetString(awsRegionEnv)
	cfg.AWS.Bucket = viper.GetString(awsBucketEnv)
	cfg.HTTP.Port = viper.GetInt(httpPortEnv)
	cfg.HTTP.Host = viper.GetString(httpHostEnv)
	cfg.HTTP.ClearLocalCachePeriod = viper.GetDuration(httpClearEnv)

	return &cfg, nil
}
