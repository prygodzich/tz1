package application

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"targetads/internal/config"
	"targetads/internal/http/server"
	"targetads/internal/logger"
	"time"
)

const (
	serverShutdownTimeout = 1 * time.Minute
)

func Run(ctx context.Context) error {
	config, err := config.Parse()
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	log, err := logger.Init(ctx, &config.Logger)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	ctx = logger.ContextWithLogger(ctx, log)

	defer func() {
		log.Debug("stop logger")
		err = log.Sync()
		if !errors.Is(err, syscall.EINVAL) {
			fmt.Fprintf(os.Stderr, "logger sync warn : %v\n", err)
		}
	}()

	log.Debug("start init app")

	hs, err := initHandler(ctx, config)
	if err != nil {
		return fmt.Errorf("init handler: %w", err)
	}

	stopHTTPServer, err := server.ServeHTTP(&config.HTTP, hs)
	if err != nil {
		return fmt.Errorf("start HTTP server: %w", err)
	}
	log.Debug(fmt.Sprintf("start http server %s:%d", config.HTTP.Host, config.HTTP.Port))
	// ошибки при завершении работы сервера не имеют значения
	// nolint:errcheck
	defer stopHTTPServer(func() context.Context {
		// утечка контекста при завершении работы не имеет значения
		// nolint:govet
		ctx, _ := context.WithTimeout(ctx, serverShutdownTimeout)
		return ctx
	}())

	<-ctx.Done()
	log.Info("stop app")
	return nil
}
