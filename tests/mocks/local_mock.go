package mocks

import (
	"context"
	"targetads/internal/storage/local"

	"github.com/stretchr/testify/mock"
)

// MockLocalService - мок для локального хранилища
type MockLocalService struct {
	mock.Mock
}

func (m *MockLocalService) Init(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockLocalService) SetDefault(ctx context.Context, typeContent local.TypeContent, data []byte) {
	m.Called(ctx, typeContent, data)
}

func (m *MockLocalService) Set(ctx context.Context, typeContent local.TypeContent, name string, data []byte) {
	m.Called(ctx, typeContent, name, data)
}

func (m *MockLocalService) GetDefault(ctx context.Context, typeContent local.TypeContent) []byte {
	args := m.Called(ctx, typeContent)
	return args.Get(0).([]byte)
}

func (m *MockLocalService) Get(ctx context.Context, typeContent local.TypeContent, name string) ([]byte, error) {
	args := m.Called(ctx, typeContent, name)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockLocalService) Clear(ctx context.Context) {
	m.Called(ctx)
}
