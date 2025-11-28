package logging

import (
	"log/slog"
	"os"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
)

type SlogConfig struct {
	Level      string
	Format     string
	OutputPath string
}

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger(cfg SlogConfig) (interfaces.Logger, error) {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	output := os.Stdout

	if cfg.OutputPath != "" && cfg.OutputPath != "stdout" {
		file, err := os.OpenFile(cfg.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		output = file
	}

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)

	return &SlogLogger{logger: logger}, nil
}

func NewDefaultLogger() interfaces.Logger {
	logger, _ := NewSlogLogger(SlogConfig{
		Level:      "info",
		Format:     "json",
		OutputPath: "stdout",
	})
	return logger
}

func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *SlogLogger) With(args ...any) interfaces.Logger {
	return &SlogLogger{
		logger: l.logger.With(args...),
	}
}
