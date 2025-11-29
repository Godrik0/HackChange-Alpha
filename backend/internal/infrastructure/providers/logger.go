package providers

import (
	"github.com/Godrik0/HackChange-Alpha/backend/internal/config"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/logging"
)

type LoggerProvider interface {
	ProvideLogger(cfg *config.Config) (interfaces.Logger, error)
}

type DefaultLoggerProvider struct{}

func (p *DefaultLoggerProvider) ProvideLogger(cfg *config.Config) (interfaces.Logger, error) {
	return logging.NewSlogLogger(logging.SlogConfig{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		OutputPath: "stdout",
	})
}
