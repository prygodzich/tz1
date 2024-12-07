package handler

import (
	"context"
	"targetads/internal/metrics"
)

func (h *Handlers) getLocalFile(ctx context.Context, fileName FileName) ([]byte, bool) {
	log := h.log
	if data, err := h.localCache.Get(ctx, fileName.Type, fileName.Name); err == nil {
		log.Debug("found file local cache ", fileName.Name)
		metrics.CacheHits.WithLabelValues(labelLocal).Inc()
		return data, true
	}
	metrics.CacheMisses.WithLabelValues(labelLocal).Inc()
	return nil, false
}
