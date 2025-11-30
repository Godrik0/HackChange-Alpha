package interfaces

import (
	"context"
	"fmt"
	"io"

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
	CalculateScoring(ctx context.Context, id int64) (*dto.ScoringResponse, error)
}

type ImportStats struct {
	SuccessCount int      `json:"success_count"`
	FailureCount int      `json:"failure_count"`
	Total        int      `json:"total"`
	Errors       []string `json:"errors,omitempty"`
}

func (stats *ImportStats) AddError(lineNum int, err error) {
	stats.FailureCount++
	if lineNum > 0 {
		stats.Errors = append(stats.Errors, fmt.Sprintf("Line %d: %v", lineNum, err))
	} else {
		stats.Errors = append(stats.Errors, err.Error())
	}
}

type ImportService interface {
	ImportClientsCSV(ctx context.Context, reader io.Reader) (*ImportStats, error)
}
