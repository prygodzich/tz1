package handler

import (
	"context"
	"targetads/internal/logger"
	"targetads/internal/metrics"
	"time"
)

const (
	workTimeout = 100 * time.Millisecond
)

type getFileParams struct {
	label      string
	labelError string
	fileName   FileName
}

func (h *Handlers) fileFromLocalOrStore(ctxRequest context.Context, params getFileParams) []byte {
	log := h.log
	ctx := logger.ContextWithLogger(ctxRequest, log)

	log.Debug("get name ", params.fileName.Name)

	data, ok := h.getLocalFile(ctx, params.fileName)
	if ok {
		metrics.FileRequestsTotal.WithLabelValues(params.label).Inc()
		return data
	}

	ctx, cancel := context.WithTimeout(ctx, workTimeout)
	defer cancel()

	dataChan := make(chan []byte, 1)

	go func() {
		h.getFile(logger.ContextWithLogger(context.Background(), log), params.fileName, dataChan)
	}()

	select {
	case data, ok := <-dataChan:
		if !ok {
			log.Error("result channel closed")
			metrics.FileRequestsTotal.WithLabelValues(params.labelError).Inc()
			return h.localCache.GetDefault(ctx, params.fileName.Type)
		}
		if data == nil {
			metrics.FileRequestsTotal.WithLabelValues(params.labelError).Inc()
			return h.localCache.GetDefault(ctx, params.fileName.Type)
		}
		metrics.FileRequestsTotal.WithLabelValues(params.label).Inc()
		return data
	case <-ctx.Done():
		log.Debug("done and get default ")
		metrics.FileRequestsTotal.WithLabelValues(params.labelError).Inc()
		return h.localCache.GetDefault(ctx, params.fileName.Type)
	}
}
