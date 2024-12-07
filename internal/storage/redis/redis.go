package redis

import (
	"context"
	"targetads/internal/apperrs"
	"targetads/internal/logger"
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	URI string `env:"REDIS_URI"`
}

const (
	defaultKey = "default"
)

type Repository struct {
	config *Config
	client *redis.Client
}

func NewRepository(ctx context.Context, config *Config) (*Repository, error) {
	opts, err := redis.ParseURL(config.URI)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	// Проверка соединения
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Repository{
		config: config, client: client,
	}, nil
}

func (r *Repository) Get(ctx context.Context, key string) ([]byte, error) {
	if r.client == nil {
		return nil, apperrs.ErrNotInitialized
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, apperrs.ErrNotFound
		}
		return nil, err
	}
	return val, nil
}

func (r *Repository) GetDefault(ctx context.Context) ([]byte, error) {
	return r.Get(ctx, defaultKey)
}

func (r *Repository) Set(ctx context.Context, key string, data []byte) error {
	log := logger.FromContext(ctx)
	log.Debug("set key", key)
	if r.client == nil {
		log.Error("not initialized")
		return apperrs.ErrNotInitialized
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := r.client.Set(ctx, key, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetDefault(ctx context.Context, value []byte) error {
	return r.Set(ctx, defaultKey, value)
}

func (r *Repository) Close(ctx context.Context) error {
	if r.client == nil {
		return nil
	}
	err := r.client.Close()
	if err != nil {
		return err
	}
	r.client = nil
	return nil
}
