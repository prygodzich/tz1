package redis

import (
	"context"
	"targetads/internal/logger"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	// Создаем мини-Redis сервер для тестирования
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	// Создаем клиент Redis, подключенный к мини-серверу
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	// Создаем репозиторий
	repo := Repository{client: client}
	ctx := context.Background()
	log, _ := logger.Init(ctx, &logger.Config{Level: "error", Format: "json"})
	ctx = logger.ContextWithLogger(ctx, log)

	// Тестируем Set и Get
	t.Run("Set and Get", func(t *testing.T) {

		key := "testKey"
		value := []byte("testValue")

		// Set
		err := repo.Set(ctx, key, value)
		assert.NoError(t, err)

		// Get
		retrievedValue, err := repo.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	// Тестируем GetDefault
	t.Run("GetDefault", func(t *testing.T) {
		defaultValue := []byte("defaultValue")

		// Set default value
		err := repo.Set(ctx, defaultKey, defaultValue)
		assert.NoError(t, err)

		// Get default
		retrievedDefault, err := repo.GetDefault(ctx)
		assert.NoError(t, err)
		assert.Equal(t, defaultValue, retrievedDefault)
	})

}
