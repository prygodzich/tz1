package handler

import (
	"context"
	"net/http"
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

func (h *Handlers) fileFromLocalOrStore(ctxRequest context.Context, params getFileParams, w http.ResponseWriter) {
	log := h.log
	ctx := logger.ContextWithLogger(ctxRequest, log)

	log.Debug("get name ", params.fileName.Name)

	data, ok := h.getLocalFile(ctx, params.fileName)
	if ok {
		_, _ = w.Write(data)
		metrics.FileRequestsTotal.WithLabelValues(params.label).Inc()
		return
	}

	ctx, cancel := context.WithTimeout(ctx, workTimeout)
	defer cancel()

	dataChan := make(chan []byte, 1)

	go func() {
		defer close(dataChan)
		data := h.getFile(logger.ContextWithLogger(context.Background(), log), params.fileName)
		select {
		case dataChan <- data:
		default:
			log.Debug("close channel")
		}
	}()

	select {
	case data, ok := <-dataChan:
		if !ok {
			log.Error("result channel closed")
			metrics.FileRequestsTotal.WithLabelValues(params.labelError).Inc()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if data == nil {
			metrics.FileRequestsTotal.WithLabelValues(params.labelError).Inc()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
		metrics.FileRequestsTotal.WithLabelValues(params.label).Inc()
	case <-ctx.Done():
		log.Debug("done and get default ")
		data = h.localCache.GetDefault(ctx, params.fileName.Type)
		_, _ = w.Write(data)
		metrics.FileRequestsTotal.WithLabelValues(params.labelError).Inc()
	}
}
