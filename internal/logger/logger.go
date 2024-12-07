package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey string

const (
	json                     = "json"
	console                  = "console"
	LoggerValueKey loggerKey = "logger"
)

type Config struct {
	Level  string
	Format string
}

type Logger = *zap.SugaredLogger

// var Log Logger

func Init(ctx context.Context, config *Config) (Logger, error) {
	logLevel, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		logLevel = zapcore.InfoLevel
	}

	core := zapcore.NewCore(
		getEncoder(config),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(logLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	// ContextWithLogger(ctx, logger.Sugar())
	// Log = logger.Sugar()
	return logger.Sugar(), nil
}

func getEncoder(config *Config) zapcore.Encoder {
	switch config.Format {
	case json:
		return jsonEncoder()
	case console:
		return consoleEncoder()
	}
	return consoleEncoder()
}

func jsonEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(encoderConfig())
}
func consoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(encoderConfig())
}

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, LoggerValueKey, logger)
}

func FromContext(ctx context.Context) Logger {
	return ctx.Value(LoggerValueKey).(Logger)
}
