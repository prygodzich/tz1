package handler

import (
	"net/http"
	"targetads/internal/metrics"
	"targetads/internal/storage/local"
	"time"
)

const (
	labelJs         = "js"
	labelVideo      = "video"
	labelJsError    = "error_js"
	labelVideoError = "error_video"
)

func (h *Handlers) js(w http.ResponseWriter, r *http.Request) {
	log := h.log
	log.Debug("js")
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.FileRequestDuration.WithLabelValues(labelJs).Observe(duration)
	}()

	name := r.URL.Query().Get(nameParam)
	if name == "" {
		metrics.FileRequestsTotal.WithLabelValues(labelJsError).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := getFileParams{
		label:      labelJs,
		labelError: labelJsError,
		fileName:   FileName{Name: name, Type: local.Js},
	}

	w.Header().Set("Content-Type", "text/javascript")
	_, _ = w.Write(h.fileFromLocalOrStore(r.Context(), params))

}

func (h *Handlers) video(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.FileRequestDuration.WithLabelValues(labelVideo).Observe(duration)
	}()

	name := r.URL.Query().Get(nameParam)
	if name == "" {
		metrics.FileRequestsTotal.WithLabelValues(labelVideoError).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := getFileParams{
		label:      labelVideo,
		labelError: labelVideoError,
		fileName:   FileName{Name: name, Type: local.Video},
	}

	w.Header().Set("Content-Type", "video/mp4")
	_, _ = w.Write(h.fileFromLocalOrStore(r.Context(), params))
}
