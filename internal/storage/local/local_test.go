package local

import (
	"context"
	"targetads/internal/apperrs"
	"targetads/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	// Создаем тестовые данные
	defaultData := map[TypeContent][]byte{
		Video: []byte("default video"),
		Js:    []byte("default js"),
	}

	// Инициализируем репозиторий
	repo := NewRepository()
	log, _ := logger.Init(context.Background(), &logger.Config{
		Level:  "info",
		Format: "json",
	})
	ctx := logger.ContextWithLogger(context.Background(), log)

	// Тестируем GetDefault
	t.Run("GetDefault", func(t *testing.T) {
		repo.SetDefault(ctx, Video, defaultData[Video])
		repo.SetDefault(ctx, Js, defaultData[Js])

		// Проверяем получение дефолтного видео
		videoData := repo.GetDefault(ctx, Video)
		assert.Equal(t, defaultData[Video], videoData)

		// Проверяем получение дефолтного JS
		jsData := repo.GetDefault(ctx, Js)
		assert.Equal(t, defaultData[Js], jsData)

		// Проверяем получение несуществующего типа контента
		unknownData := repo.GetDefault(ctx, Unknown)
		assert.Nil(t, unknownData)
	})

	// Тестируем Get
	t.Run("Get", func(t *testing.T) {
		// Добавляем тестовые данные
		testVideo := []byte("test video")
		testJS := []byte("test js")
		repo.Set(ctx, Video, "test_video", testVideo)
		repo.Set(ctx, Js, "test_js", testJS)

		// Проверяем получение существующих данных
		retrievedVideo, err := repo.Get(ctx, Video, "test_video")
		require.NoError(t, err)
		assert.Equal(t, testVideo, retrievedVideo)

		retrievedJS, err := repo.Get(ctx, Js, "test_js")
		require.NoError(t, err)
		assert.Equal(t, testJS, retrievedJS)

		// Проверяем получение несуществующих данных
		_, err = repo.Get(ctx, Video, "non_existent")
		assert.Error(t, err)
		assert.ErrorIs(t, err, apperrs.ErrNotFound)
	})

	// Тестируем Set
	t.Run("Set", func(t *testing.T) {
		// Устанавливаем новые данные
		newVideo := []byte("new video")
		repo.Set(ctx, Video, "new_video", newVideo)

		// Проверяем, что данные успешно установлены
		retrievedVideo, err := repo.Get(ctx, Video, "new_video")
		require.NoError(t, err)
		assert.Equal(t, newVideo, retrievedVideo)

		// Проверяем перезапись существующих данных
		updatedVideo := []byte("updated video")
		repo.Set(ctx, Video, "new_video", updatedVideo)

		retrievedUpdatedVideo, err := repo.Get(ctx, Video, "new_video")
		require.NoError(t, err)
		assert.Equal(t, updatedVideo, retrievedUpdatedVideo)
	})

	// Тестируем Delete
	t.Run("Clear", func(t *testing.T) {
		repo.SetDefault(ctx, Video, defaultData[Video])
		repo.SetDefault(ctx, Js, defaultData[Js])

		// Добавляем данные для удаления
		dataToDelete := []byte("data to delete")
		repo.Set(ctx, Video, "delete_me", dataToDelete)

		// Проверяем, что данные существуют
		_, err := repo.Get(ctx, Video, "delete_me")
		require.NoError(t, err)

		// Удаляем данные
		repo.Clear(ctx)

		// Проверяем, что данные удалены
		_, err = repo.Get(ctx, Video, "delete_me")
		assert.Error(t, err)
		assert.ErrorIs(t, err, apperrs.ErrNotFound)
		// assert.Equal(t, apperrs.ErrNotFound, err)

		// Проверяем что дефолтные остались
		videoData := repo.GetDefault(ctx, Video)
		assert.Equal(t, defaultData[Video], videoData)
		jsData := repo.GetDefault(ctx, Js)
		assert.Equal(t, defaultData[Js], jsData)
	})
}
