package application

import (
	"context"
	"targetads/internal/storage/aws"
)

type awsDomain struct {
	service *aws.Service
}

func buildAwsDomain(ctx context.Context, config *aws.Config) (awsDomain, error) {
	repository, err := aws.NewRepository(config)
	if err != nil {
		return awsDomain{}, err
	}
	service := aws.NewService(ctx, repository)
	return awsDomain{
		service: service,
	}, nil
}
