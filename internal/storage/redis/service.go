package redis

import "context"

type repository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte) error
	GetDefault(ctx context.Context) ([]byte, error)
	SetDefault(ctx context.Context, data []byte) error
	Close(ctx context.Context) error
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
	return s.repo.Get(ctx, name)
}

func (s *Service) Set(ctx context.Context, name string, data []byte) error {
	return s.repo.Set(ctx, name, data)
}

func (s *Service) GetDefault(ctx context.Context) ([]byte, error) {
	return s.repo.GetDefault(ctx)
}

func (s *Service) SetDefault(ctx context.Context, data []byte) error {
	return s.repo.SetDefault(ctx, data)
}

func (s *Service) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}
