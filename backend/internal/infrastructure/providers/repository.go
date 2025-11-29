package providers

import (
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	storage "github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type RepositoryProvider interface {
	ProvideClientRepository(db *gorm.DB, logger interfaces.Logger) interfaces.ClientRepository
}

type DefaultRepositoryProvider struct{}

func (p *DefaultRepositoryProvider) ProvideClientRepository(db *gorm.DB, logger interfaces.Logger) interfaces.ClientRepository {
	return storage.NewClientRepository(db, logger)
}
