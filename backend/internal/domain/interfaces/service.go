package interfaces

import (
	"context"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type ClientService interface {
	GetClient(ctx context.Context, id int64) (*models.Client, error)

	SearchClients(ctx context.Context, params dto.SearchParams) ([]models.Client, error)

	CreateClient(ctx context.Context, req *dto.CreateClientRequest) (*models.Client, error)

	UpdateClient(ctx context.Context, id int64, req *dto.UpdateClientRequest) (*models.Client, error)

	DeleteClient(ctx context.Context, id int64) error

	ListClients(ctx context.Context, limit, offset int) ([]models.Client, error)
}

type ScoringService interface {
	CalculateScoring(ctx context.Context, id int64) (*models.ScoringResult, error)
}
