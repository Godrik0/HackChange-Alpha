package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
	"gorm.io/datatypes"
)

type clientService struct {
	clientRepo interfaces.ClientRepository
	mlService  interfaces.MLService
	logger     interfaces.Logger
}

func NewClientService(
	clientRepo interfaces.ClientRepository,
	mlService interfaces.MLService,
	logger interfaces.Logger,
) interfaces.ClientService {
	return &clientService{
		clientRepo: clientRepo,
		mlService:  mlService,
		logger:     logger.With("component", "ClientService"),
	}
}

func (s *clientService) GetClient(ctx context.Context, id int64) (*models.Client, error) {
	s.logger.Debug("Getting client", "id", id)

	client, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get client", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return client, nil
}

func (s *clientService) SearchClients(ctx context.Context, params dto.SearchParams) ([]models.Client, error) {
	if params.IsEmpty() {
		return nil, fmt.Errorf("at least one search parameter must be provided")
	}

	s.logger.Debug("Searching clients", "params", params)

	clients, err := s.clientRepo.Search(ctx, params)
	if err != nil {
		s.logger.Error("Failed to search clients", "error", err)
		return nil, fmt.Errorf("failed to search clients: %w", err)
	}

	return clients, nil
}

func (s *clientService) CalculateScoring(ctx context.Context, id int64) (*models.ScoringResult, error) {
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

func (s *clientService) CreateClient(ctx context.Context, req *dto.CreateClientRequest) (*models.Client, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	s.logger.Debug("Creating client", "first_name", req.FirstName, "last_name", req.LastName)

	client, err := req.ToModel()
	if err != nil {
		s.logger.Error("Failed to convert DTO to model", "error", err)
		return nil, fmt.Errorf("failed to convert DTO to model: %w", err)
	}

	if err := s.validateClient(client); err != nil {
		s.logger.Warn("Client validation failed", "error", err)
		return nil, fmt.Errorf("client validation failed: %w", err)
	}

	if err := s.clientRepo.Create(ctx, client); err != nil {
		s.logger.Error("Failed to create client", "error", err)
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	s.logger.Info("Client created successfully", "id", client.ID)
	return client, nil
}

func (s *clientService) UpdateClient(ctx context.Context, id int64, req *dto.UpdateClientRequest) (*models.Client, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	s.logger.Debug("Updating client", "id", id)

	client, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get client for update", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	if req.FirstName != "" {
		client.FirstName = req.FirstName
	}
	if req.LastName != "" {
		client.LastName = req.LastName
	}
	if req.BirthDate != "" {
		parsedDate, err := time.Parse(dto.DateFormat, req.BirthDate)
		if err != nil {
			s.logger.Warn("Invalid birth date format", "input", req.BirthDate, "error", err)
			return nil, fmt.Errorf("invalid birth date format (expected YYYY-MM-DD): %w", err)
		}
		client.BirthDate = parsedDate
	}
	if req.CoreData != nil {
		coreDataJSON, err := json.Marshal(req.CoreData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal core data: %w", err)
		}
		client.CoreData = datatypes.JSON(coreDataJSON)
	}
	if req.Features != nil {
		featuresJSON, err := json.Marshal(req.Features)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal features: %w", err)
		}
		client.Features = datatypes.JSON(featuresJSON)
	}

	if err := s.validateClient(client); err != nil {
		s.logger.Warn("Client validation failed", "id", id, "error", err)
		return nil, fmt.Errorf("client validation failed: %w", err)
	}

	if err := s.clientRepo.Update(ctx, client); err != nil {
		s.logger.Error("Failed to update client", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	s.logger.Info("Client updated successfully", "id", id)
	return client, nil
}

func (s *clientService) DeleteClient(ctx context.Context, id int64) error {
	s.logger.Debug("Deleting client", "id", id)

	if err := s.clientRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete client", "id", id, "error", err)
		return fmt.Errorf("failed to delete client: %w", err)
	}

	s.logger.Info("Client deleted successfully", "id", id)
	return nil
}

func (s *clientService) ListClients(ctx context.Context, limit, offset int) ([]models.Client, error) {
	s.logger.Debug("Listing clients", "limit", limit, "offset", offset)

	clients, err := s.clientRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list clients", "error", err)
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}

	return clients, nil
}

func (s *clientService) extractFeatures(client *models.Client) (map[string]interface{}, error) {
	features := make(map[string]interface{})

	features["first_name"] = client.FirstName
	features["last_name"] = client.LastName
	features["birth_date"] = client.BirthDate

	if len(client.CoreData) > 0 {
		var coreData map[string]interface{}
		if err := json.Unmarshal(client.CoreData, &coreData); err == nil {
			for k, v := range coreData {
				features[k] = v
			}
		}
	}

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

func (s *clientService) validateClient(client *models.Client) error {
	if !client.IsValid() {
		return fmt.Errorf("client has invalid fields")
	}
	return nil
}
