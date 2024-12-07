package handler

import (
	"context"
	"fmt"
	"targetads/internal/apperrs"
	"targetads/internal/logger"
	"targetads/internal/metrics"
	"targetads/internal/storage/local"
	"time"
)

const (
	jsPrefix    = "js_"
	videoPrefix = "video_"
	nameParam   = "name"
	labelLocal  = "local"
	labelAws    = "aws"
	labelFast   = "fast"

	fileDownloadTimeout = 10 * time.Second
)

type FileName struct {
	Name string
	Type local.TypeContent
}

func (h *Handlers) getFile(ctx context.Context, fileName FileName) []byte {
	log := h.log
	ctx2, cancel := context.WithTimeout(context.Background(), fileDownloadTimeout)
	ctx2 = logger.ContextWithLogger(ctx2, log)

	resultChan := make(chan []byte, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer cancel()
		defer func() {
			log.Debug("close channels")
			close(resultChan)
			close(errorChan)
		}()

		data, err := h.fastCache.Get(ctx2, fileName.Name)
		if err == nil {
			log.Debug("get file from fast cache ", fileName.Name)
			resultChan <- data
			h.localCache.Set(ctx2, fileName.Type, fileName.Name, data)
			metrics.CacheHits.WithLabelValues(labelFast).Inc()
			return
		}
		metrics.CacheMisses.WithLabelValues(labelFast).Inc()

		data, err = h.store.Get(ctx2, fileName.Name)
		if err != nil {
			// metrics.FileRequestsTotal.WithLabelValues(labelAws).Inc()
			errorChan <- apperrs.ErrNotFound
			return
		}
		log.Debug("get file from aws ", fileName.Name)
		resultChan <- data
		h.localCache.Set(ctx2, fileName.Type, fileName.Name, data)
		log.Debug("try set file to fast cache ", fileName.Name)
		err = h.fastCache.Set(ctx2, fileName.Name, data)
		if err != nil {
			log.Error("set file to fast cache ", fileName.Name, err)
		}
		log.Debug("set file to fast cache ", fileName.Name)
	}()

	select {
	case result := <-resultChan:
		log.Debug("file found ", fileName.Name)
		return result
	case err := <-errorChan:
		if err != nil {
			log.Error(fmt.Sprintf("file not found %s error: %s", fileName.Name, err))
		}
		log.Debug("error and get default ")
		return h.localCache.GetDefault(ctx2, fileName.Type)
	case <-ctx2.Done():
		log.Debug("done and get default ")
		return h.localCache.GetDefault(ctx2, fileName.Type)
	}
}
