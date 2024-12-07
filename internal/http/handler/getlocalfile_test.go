package handler

import (
	"context"
	"net/http"
	"testing"

	"targetads/internal/apperrs"
	"targetads/internal/logger"
	"targetads/internal/storage/local"

	"github.com/stretchr/testify/assert"
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

func TestGetLocalFile(t *testing.T) {
	log, _ := logger.Init(context.Background(), &logger.Config{Level: "error", Format: "console"})

	tests := []struct {
		name           string
		typeContent    local.TypeContent
		url            string
		fileName       string
		mockData       []byte
		mockError      error
		expectedStatus bool
		expectedBody   []byte
	}{
		{
			name:           "Successful get JS file",
			typeContent:    local.Js,
			url:            "/js?name=test.js",
			fileName:       "test.js",
			mockData:       []byte("console.log('test');"),
			mockError:      nil,
			expectedStatus: true,
			expectedBody:   []byte("console.log('test');"),
		},
		{
			name:           "File not found",
			typeContent:    local.Js,
			url:            "/js?name=nonexistent.js",
			fileName:       "nonexistent.js",
			mockData:       nil,
			mockError:      apperrs.ErrNotFound,
			expectedStatus: false,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockLocalService)
			mockStorage.On("Get", mock.Anything, tt.typeContent, tt.fileName).Return(tt.mockData, tt.mockError)

			handler := &Handlers{
				localCache: mockStorage,
				log:        log,
				store:      nil,
				fastCache:  nil,
			}

			req, err := http.NewRequest("GET", tt.url, nil)
			assert.NoError(t, err)

			fname := FileName{Name: tt.fileName, Type: tt.typeContent}
			res, ok := handler.getLocalFile(req.Context(), fname)

			assert.Equal(t, ok, tt.expectedStatus)
			assert.Equal(t, tt.expectedBody, res)

			mockStorage.AssertExpectations(t)
		})
	}
}
