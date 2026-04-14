package logger

import (
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	Debug *zap.Logger
	Info  *zap.Logger
	Warn  *zap.Logger
	Error *zap.Logger
}

func NewLogger() *Logger {

	level := os.Getenv("LOG_LEVEL")

	cfg := zap.NewProductionConfig()

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	base, _ := cfg.Build()

	return &Logger{
		Debug: base.Named("debug"),
		Info:  base.Named("info"),
		Warn:  base.Named("warn"),
		Error: base.Named("error"),
	}
}
