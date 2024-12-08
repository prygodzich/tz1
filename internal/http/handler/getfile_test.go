package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"targetads/internal/logger"
	"targetads/internal/storage/local"
	"targetads/tests/mocks"
)

const (
	logLevel = "info"
)

func TestHandlers_getFile(t *testing.T) {
	log, _ := logger.Init(context.Background(), &logger.Config{Level: logLevel, Format: "console"})

	tests := []struct {
		name           string
		fileName       FileName
		setupMocks     func(*mocks.MockLocalService, *mocks.MockFastCacheService, *mocks.MockStoreService)
		expectedResult []byte
		expectedError  bool
	}{
		{
			name:     "Fast cache hit",
			fileName: FileName{Name: "test.js", Type: local.Js},
			setupMocks: func(localCache *mocks.MockLocalService, fastCache *mocks.MockFastCacheService, store *mocks.MockStoreService) {
				localCache.On("Set", mock.Anything, local.Js, "test.js", []byte("fast cache data"))
				fastCache.On("Get", mock.Anything, "test.js").Return([]byte("fast cache data"), nil)
			},
			expectedResult: []byte("fast cache data"),
			expectedError:  false,
		},
		{
			name:     "Fast cache miss, store hit",
			fileName: FileName{Name: "test.js", Type: local.Js},
			setupMocks: func(localCache *mocks.MockLocalService, fastCache *mocks.MockFastCacheService, store *mocks.MockStoreService) {
				fastCache.On("Get", mock.Anything, "test.js").Return(nil, errors.New("not found"))
				store.On("Get", mock.Anything, "test.js").Return([]byte("store data"), nil)
				localCache.On("Set", mock.Anything, local.Js, "test.js", []byte("store data"))
				fastCache.On("Set", mock.Anything, "test.js", []byte("store data")).Return(nil)
			},
			expectedResult: []byte("store data"),
			expectedError:  false,
		},
		{
			name:     "Fast cache miss, store miss",
			fileName: FileName{Name: "test.js", Type: local.Js},
			setupMocks: func(localCache *mocks.MockLocalService, fastCache *mocks.MockFastCacheService, store *mocks.MockStoreService) {
				fastCache.On("Get", mock.Anything, "test.js").Return(nil, errors.New("not found"))
				store.On("Get", mock.Anything, "test.js").Return(nil, errors.New("not found"))
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	ctx := logger.ContextWithLogger(context.Background(), log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fastCache := new(mocks.MockFastCacheService)
			localCache := new(mocks.MockLocalService)
			store := new(mocks.MockStoreService)

			tt.setupMocks(localCache, fastCache, store)

			h := &Handlers{
				log:        log,
				fastCache:  fastCache,
				localCache: localCache,
				store:      store,
			}

			dataChan := make(chan []byte, 1)
			go h.getFile(ctx, tt.fileName, dataChan)

			var result []byte
			select {
			case data, ok := <-dataChan:
				if !ok {
					log.Debug("result channel closed")
				}

				result = data
			case <-time.After(time.Second):
				if !tt.expectedError {
					t.Fatal("Test timed out")
				}
			}

			if tt.expectedError {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expectedResult, result)
			}

			fastCache.AssertExpectations(t)
			localCache.AssertExpectations(t)
			store.AssertExpectations(t)
		})
	}
}
