package application

import (
	"context"
	"targetads/internal/storage/local"
)

type localDomain struct {
	service *local.ServiceImpl
}

func buildLocalStorageDomain(ctx context.Context) localDomain {

	repository := local.NewRepository()
	service := local.NewService(repository)
	service.Init(ctx)

	return localDomain{
		service: service,
	}
}
