package aws

import "context"

type repository interface {
	DownloadFileFromBucket(ctx context.Context, filename string) ([]byte, error)
}

type Service struct {
	repo repository
}

func NewService(ctx context.Context, repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, name string) ([]byte, error) {
	return s.repo.DownloadFileFromBucket(ctx, name)
}
