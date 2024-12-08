package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockLocalService - мок для локального хранилища
type MockFastCacheService struct {
	mock.Mock
}

func (m *MockFastCacheService) Get(ctx context.Context, name string) ([]byte, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFastCacheService) Set(ctx context.Context, name string, data []byte) error {
	args := m.Called(ctx, name, data)
	return args.Error(0)
}

func (m *MockFastCacheService) GetDefault(ctx context.Context) ([]byte, error) {
	args := m.Called(ctx)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFastCacheService) SetDefault(ctx context.Context, data []byte) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockFastCacheService) Close(ctx context.Context) error {
	m.Called(ctx)
	return nil
}
