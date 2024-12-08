package redis

import "context"

type repository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte) error
	GetDefault(ctx context.Context) ([]byte, error)
	SetDefault(ctx context.Context, data []byte) error
	Close(ctx context.Context) error
}

type ServiceImpl struct {
	repo repository
}

type Service interface {
	Get(ctx context.Context, name string) ([]byte, error)
	Set(ctx context.Context, name string, data []byte) error
	GetDefault(ctx context.Context) ([]byte, error)
	SetDefault(ctx context.Context, data []byte) error
	Close(ctx context.Context) error
}

func NewService(ctx context.Context, repo repository) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) Get(ctx context.Context, name string) ([]byte, error) {
	return s.repo.Get(ctx, name)
}

func (s *ServiceImpl) Set(ctx context.Context, name string, data []byte) error {
	return s.repo.Set(ctx, name, data)
}

func (s *ServiceImpl) GetDefault(ctx context.Context) ([]byte, error) {
	return s.repo.GetDefault(ctx)
}

func (s *ServiceImpl) SetDefault(ctx context.Context, data []byte) error {
	return s.repo.SetDefault(ctx, data)
}

func (s *ServiceImpl) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}
