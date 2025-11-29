package mocks

import (
	"context"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type MockMLService struct {
	PredictFunc          func(ctx context.Context, features map[string]interface{}) (*models.ScoringResult, error)
	SendTrainingDataFunc func(ctx context.Context, data interface{}) error
	HealthCheckFunc      func(ctx context.Context) error
}

func (m *MockMLService) Predict(ctx context.Context, features map[string]interface{}) (*models.ScoringResult, error) {
	if m.PredictFunc != nil {
		return m.PredictFunc(ctx, features)
	}
	return &models.ScoringResult{Score: 0.75}, nil
}

func (m *MockMLService) SendTrainingData(ctx context.Context, data interface{}) error {
	if m.SendTrainingDataFunc != nil {
		return m.SendTrainingDataFunc(ctx, data)
	}
	return nil
}

func (m *MockMLService) HealthCheck(ctx context.Context) error {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc(ctx)
	}
	return nil
}
