package application

import (
	"context"
	"targetads/internal/logger"
	"targetads/internal/storage/redis"
)

type redisDomain struct {
	service *redis.Service
}

func buildRedisDomain(ctx context.Context, config *redis.Config) (*redisDomain, error) {
	log := logger.FromContext(ctx)
	repository, err := redis.NewRepository(ctx, config)
	if err != nil {
		log.Errorf("build redis repository: %v", err)
		return nil, err
	}

	service := redis.NewService(ctx, repository)
	return &redisDomain{service: service}, nil
}
