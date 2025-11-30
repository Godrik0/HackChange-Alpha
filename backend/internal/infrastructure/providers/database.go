package providers

import (
	"github.com/Godrik0/HackChange-Alpha/backend/internal/config"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	storage "github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type DatabaseProvider interface {
	ProvideDatabase(cfg *config.Config, logger interfaces.Logger) (*gorm.DB, error)
	RunMigrations(db *gorm.DB, logger interfaces.Logger) error
}

type PostgresProvider struct{}

func (p *PostgresProvider) ProvideDatabase(cfg *config.Config, logger interfaces.Logger) (*gorm.DB, error) {
	return storage.NewPostgresRepository(cfg, logger)
}

func (p *PostgresProvider) RunMigrations(db *gorm.DB, logger interfaces.Logger) error {
	return storage.RunMigrations(db, logger)
}
