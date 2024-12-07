package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"targetads/internal/logger"
	"time"
)

type Config struct {
	Port                  int           `env:"PORT" default:"8089"`
	Host                  string        `env:"HOST" default:"localhost"`
	ClearLocalCachePeriod time.Duration `env:"CLEAR_LOCAL_CACHE_PERIOD" default:"10m"`
}

type (
	closeFn func(ctx context.Context) error
	Routes  struct{}
)

func ServeHTTP(cfg *Config, handler http.Handler) (closeFn, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	s := http.Server{Handler: handler}
	// NOTE::
	// выкинем панику и остановим все приложение в случае ошибки старта сервера.
	// это редкий случай, т.к port мы слушать уже начали.
	go func() {
		err := s.Serve(l)
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	closer := func(ctx context.Context) error {
		err := s.Shutdown(ctx)
		log := ctx.Value(logger.LoggerValueKey).(logger.Logger)
		log.Debug("stop http server")
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return closer, nil
}
