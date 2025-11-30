package interfaces

import (
	"context"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type MLService interface {
	Predict(ctx context.Context, features map[string]interface{}) (*models.ScoringResult, error)
	PredictWithExplanation(ctx context.Context, features map[string]interface{}) (*dto.MLScoringResponse, error)
	SendTrainingData(ctx context.Context, data interface{}) error
	HealthCheck(ctx context.Context) error
}
