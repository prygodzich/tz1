package application

import (
	"context"
	"fmt"
	"targetads/internal/config"
	"targetads/internal/http/handler"
	"targetads/internal/logger"
	"targetads/internal/storage/local"
	"time"

	"github.com/go-chi/chi"
)

const (
	DefaultVideoName = "default.avi"
	DefaultJsName    = "default.js"
)

func initHandler(ctx context.Context, config *config.Config) (*chi.Mux, error) {
	log := logger.FromContext(ctx)
	localDomain := buildLocalStorageDomain(ctx)
	awsDomain, err := buildAwsDomain(ctx, &config.AWS)
	if err != nil {
		return nil, fmt.Errorf("build aws domain: %w", err)
	}
	redisDomain, err := buildRedisDomain(ctx, &config.Redis)
	if err != nil {
		return nil, fmt.Errorf("build redis domain: %w", err)
	}

	defaultVideoFile, err := awsDomain.service.Get(ctx, DefaultVideoName)
	if err != nil {
		return nil, fmt.Errorf("download file: %w", err)
	}
	localDomain.service.SetDefault(ctx, local.Video, defaultVideoFile)

	defaultJsFile, err := awsDomain.service.Get(ctx, DefaultJsName)
	if err != nil {
		return nil, fmt.Errorf("download file: %w", err)
	}
	localDomain.service.SetDefault(ctx, local.Js, defaultJsFile)

	go func() {
		ticker := time.NewTicker(config.HTTP.ClearLocalCachePeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Debug("clear local cache")
				localDomain.service.Clear(ctx)
			case <-ctx.Done():
				log.Debug("stop clear local cache")
				return
			}
		}
	}()

	hs := handler.NewHandlers(localDomain.service, redisDomain.service, awsDomain.service, log)
	return hs, nil
}
