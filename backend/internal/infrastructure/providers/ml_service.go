package providers

import (
	"github.com/Godrik0/HackChange-Alpha/backend/internal/config"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/ml"
)

type MLServiceProvider interface {
	ProvideMLService(cfg *config.Config, logger interfaces.Logger) interfaces.MLService
}

type DefaultMLServiceProvider struct{}

func (p *DefaultMLServiceProvider) ProvideMLService(cfg *config.Config, logger interfaces.Logger) interfaces.MLService {
	return ml.NewMLClient(cfg.ML.BaseURL, cfg.ML.Timeout, logger)
}
