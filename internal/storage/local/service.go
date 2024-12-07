package local

import (
	"context"
	"targetads/internal/logger"
)

type repository interface {
	Init(ctx context.Context)
	SetDefault(ctx context.Context, typeContent TypeContent, data []byte)
	Set(ctx context.Context, typeContent TypeContent, name string, data []byte)
	GetDefault(ctx context.Context, typeContent TypeContent) []byte
	Get(ctx context.Context, typeContent TypeContent, name string) ([]byte, error)
	Clear(ctx context.Context)
}

type Service interface {
	Init(ctx context.Context)
	SetDefault(ctx context.Context, typeContent TypeContent, data []byte)
	Set(ctx context.Context, typeContent TypeContent, name string, data []byte)
	GetDefault(ctx context.Context, typeContent TypeContent) []byte
	Get(ctx context.Context, typeContent TypeContent, name string) ([]byte, error)
	Clear(ctx context.Context)
}

type ServiceImpl struct {
	repo repository
}

func NewService(repo repository) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) Init(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Debug("init local storage")
	s.repo.Init(ctx)
}

func (s *ServiceImpl) SetDefault(ctx context.Context, typeContent TypeContent, data []byte) {
	s.repo.SetDefault(ctx, typeContent, data)
}

func (s *ServiceImpl) Set(ctx context.Context, typeContent TypeContent, name string, data []byte) {
	s.repo.Set(ctx, typeContent, name, data)
}

func (s *ServiceImpl) GetDefault(ctx context.Context, typeContent TypeContent) []byte {
	return s.repo.GetDefault(ctx, typeContent)
}

func (s *ServiceImpl) Get(ctx context.Context, typeContent TypeContent, name string) ([]byte, error) {
	return s.repo.Get(ctx, typeContent, name)
}

func (s *ServiceImpl) Clear(ctx context.Context) {
	s.repo.Clear(ctx)
}
