package aws

import "context"

type repository interface {
	DownloadFileFromBucket(ctx context.Context, filename string) ([]byte, error)
}

type ServiceImpl struct {
	repo repository
}

type Service interface {
	Get(ctx context.Context, name string) ([]byte, error)
}

func NewService(ctx context.Context, repo repository) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) Get(ctx context.Context, name string) ([]byte, error) {
	return s.repo.DownloadFileFromBucket(ctx, name)
}
