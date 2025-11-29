package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type scoringService struct {
	clientRepo interfaces.ClientRepository
	mlService  interfaces.MLService
	logger     interfaces.Logger
}

func NewScoringService(
	clientRepo interfaces.ClientRepository,
	mlService interfaces.MLService,
	logger interfaces.Logger,
) interfaces.ScoringService {
	return &scoringService{
		clientRepo: clientRepo,
		mlService:  mlService,
		logger:     logger.With("component", "ScoringService"),
	}
}

func (s *scoringService) CalculateScoring(ctx context.Context, id int64) (*models.ScoringResult, error) {
	s.logger.Debug("Calculating scoring", "client_id", id)

	client, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get client for scoring", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get client for scoring: %w", err)
	}

	features, err := s.extractFeatures(client)
	s.logger.Debug("Extracted features", "features", features)
	if err != nil {
		s.logger.Error("Failed to extract features", "client_id", id, "error", err)
		return nil, fmt.Errorf("failed to extract features: %w", err)
	}

	result, err := s.mlService.Predict(ctx, features)
	if err != nil {
		s.logger.Error("Failed to predict scoring", "client_id", id, "error", err)
		return nil, fmt.Errorf("failed to predict scoring: %w", err)
	}

	s.logger.Info("Scoring calculated successfully", "client_id", id, "score", result.Score)
	return result, nil
}

func (s *scoringService) extractFeatures(client *models.Client) (map[string]interface{}, error) {
	features := make(map[string]interface{})

	features["first_name"] = client.FirstName
	features["last_name"] = client.LastName
	features["birth_date"] = client.BirthDate.Format(dto.DateFormat)

	if len(client.Features) > 0 {
		var clientFeatures map[string]interface{}
		if err := json.Unmarshal(client.Features, &clientFeatures); err == nil {
			for k, v := range clientFeatures {
				features[k] = v
			}
		}
	}

	if len(features) == 0 {
		return nil, fmt.Errorf("no features available for client %d", client.ID)
	}

	return features, nil
}
