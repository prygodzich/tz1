package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockStoreService struct {
	mock.Mock
}

func (m *MockStoreService) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
