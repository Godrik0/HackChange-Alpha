package interfaces

import (
	"context"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type MLService interface {
	Predict(ctx context.Context, features map[string]interface{}) (*models.ScoringResult, error)

	SendTrainingData(ctx context.Context, data interface{}) error

	HealthCheck(ctx context.Context) error
}
