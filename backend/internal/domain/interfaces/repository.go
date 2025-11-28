package interfaces

import (
	"context"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) error

	GetByID(ctx context.Context, id int64) (*models.Client, error)

	Search(ctx context.Context, params dto.SearchParams) ([]models.Client, error)

	Update(ctx context.Context, client *models.Client) error

	Delete(ctx context.Context, id int64) error

	List(ctx context.Context, limit, offset int) ([]models.Client, error)
}
