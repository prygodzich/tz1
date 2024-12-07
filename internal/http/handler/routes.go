package handler

import (
	"targetads/internal/logger"
	"targetads/internal/storage/aws"
	"targetads/internal/storage/local"
	"targetads/internal/storage/redis"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	Handlers struct {
		localCache local.Service
		fastCache  *redis.Service
		store      *aws.Service
		log        logger.Logger
	}
)

func NewHandlers(
	localCache local.Service,
	fastCache *redis.Service,
	store *aws.Service,
	log logger.Logger,
) *chi.Mux {
	r := chi.NewMux()

	h := Handlers{
		localCache: localCache,
		fastCache:  fastCache,
		store:      store,
		log:        log,
	}
	h.build(r)

	return r
}

func (h *Handlers) build(r chi.Router) {
	r.Use(middleware.Recoverer)
	r.Get("/js", h.js)
	r.Get("/video", h.video)
	r.Handle("/metrics", promhttp.Handler())
}
