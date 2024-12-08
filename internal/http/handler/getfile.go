package handler

import (
	"context"
	"targetads/internal/metrics"
	"targetads/internal/storage/local"
)

const (
	jsPrefix    = "js_"
	videoPrefix = "video_"
	nameParam   = "name"
	labelLocal  = "local"
	labelAws    = "aws"
	labelFast   = "fast"
)

type FileName struct {
	Name string
	Type local.TypeContent
}

func (h *Handlers) getFile(ctx context.Context, fileName FileName, dataChan chan<- []byte) {
	defer close(dataChan)
	log := h.log

	data, err := h.fastCache.Get(ctx, fileName.Name)
	if err == nil {
		log.Debug("get file from fast cache ", fileName.Name)
		dataChan <- data
		h.localCache.Set(ctx, fileName.Type, fileName.Name, data)
		metrics.CacheHits.WithLabelValues(labelFast).Inc()
		return
	}
	metrics.CacheMisses.WithLabelValues(labelFast).Inc()

	data, err = h.store.Get(ctx, fileName.Name)
	if err != nil {
		// metrics.FileRequestsTotal.WithLabelValues(labelAws).Inc()
		log.Error("not found file in store ", fileName.Name)
		return
	}
	log.Debug("get file from aws ", fileName.Name)
	dataChan <- data
	h.localCache.Set(ctx, fileName.Type, fileName.Name, data)
	log.Debug("try set file to fast cache ", fileName.Name)
	err = h.fastCache.Set(ctx, fileName.Name, data)
	if err != nil {
		log.Error("set file to fast cache ", fileName.Name, err)
	}
}
