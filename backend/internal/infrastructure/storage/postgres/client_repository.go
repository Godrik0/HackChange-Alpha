package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
	domainerrors "github.com/Godrik0/HackChange-Alpha/backend/internal/domain/errors"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
	"gorm.io/gorm"
)

type clientRepository struct {
	db     *gorm.DB
	logger interfaces.Logger
}

func NewClientRepository(db *gorm.DB, logger interfaces.Logger) interfaces.ClientRepository {
	return &clientRepository{
		db:     db,
		logger: logger.With("component", "ClientRepository"),
	}
}

func (r *clientRepository) Create(ctx context.Context, client *models.Client) error {
	if client == nil {
		return fmt.Errorf("client cannot be nil")
	}

	r.logger.Debug("Creating new client", "first_name", client.FirstName, "last_name", client.LastName)

	result := r.db.WithContext(ctx).Create(client)
	if result.Error != nil {
		r.logger.Error("Failed to create client", "error", result.Error)
		return fmt.Errorf("failed to create client: %w", result.Error)
	}

	r.logger.Info("Client created successfully", "id", client.ID)
	return nil
}

func (r *clientRepository) BatchCreate(ctx context.Context, clients []*models.Client) (int, error) {
	if len(clients) == 0 {
		return 0, nil
	}

	r.logger.Debug("Batch creating clients", "count", len(clients))

	result := r.db.WithContext(ctx).CreateInBatches(clients, 100)
	if result.Error != nil {
		r.logger.Error("Failed to batch create clients", "error", result.Error)
		return 0, fmt.Errorf("failed to batch create clients: %w", result.Error)
	}

	createdCount := int(result.RowsAffected)
	r.logger.Info("Clients batch created successfully", "count", createdCount)
	return createdCount, nil
}

func (r *clientRepository) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	if id <= 0 {
		return nil, domainerrors.ErrInvalidClientID
	}

	r.logger.Debug("Getting client by ID", "id", id)

	var client models.Client
	result := r.db.WithContext(ctx).First(&client, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Warn("Client not found", "id", id)
			return nil, domainerrors.ErrClientNotFound
		}
		r.logger.Error("Failed to get client", "id", id, "error", result.Error)
		return nil, fmt.Errorf("failed to get client: %w", result.Error)
	}

	r.logger.Debug("Client retrieved successfully", "id", id)
	return &client, nil
}

func (r *clientRepository) Search(ctx context.Context, params dto.SearchParams) ([]models.Client, error) {
	r.logger.Debug("Searching clients", "params", params)

	var clients []models.Client
	query := r.db.WithContext(ctx).Model(&models.Client{})

	if params.FirstName != "" {
		escaped := EscapeLikePattern(params.FirstName)
		query = query.Where("LOWER(first_name) LIKE LOWER(?)", "%"+escaped+"%")
	}

	if params.LastName != "" {
		escaped := EscapeLikePattern(params.LastName)
		query = query.Where("LOWER(last_name) LIKE LOWER(?)", "%"+escaped+"%")
	}

	if params.MiddleName != "" {
		escaped := EscapeLikePattern(params.MiddleName)
		query = query.Where("LOWER(middle_name) LIKE LOWER(?)", "%"+escaped+"%")
	}

	if params.BirthDate != "" {
		parsedDate, err := time.Parse(dto.DateFormat, params.BirthDate)
		if err == nil {
			query = query.Where("birth_date = ?", parsedDate)
		}
	}

	result := query.Find(&clients)
	if result.Error != nil {
		r.logger.Error("Failed to search clients", "error", result.Error)
		return nil, fmt.Errorf("failed to search clients: %w", result.Error)
	}

	r.logger.Info("Clients search completed", "count", len(clients))
	return clients, nil
}

func (r *clientRepository) Update(ctx context.Context, client *models.Client) error {
	if client == nil {
		return fmt.Errorf("client cannot be nil")
	}

	if client.ID <= 0 {
		return domainerrors.ErrInvalidClientID
	}

	r.logger.Debug("Updating client", "id", client.ID)

	result := r.db.WithContext(ctx).Save(client)
	if result.Error != nil {
		r.logger.Error("Failed to update client", "id", client.ID, "error", result.Error)
		return fmt.Errorf("failed to update client: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Client not found for update", "id", client.ID)
		return domainerrors.ErrClientNotFound
	}

	r.logger.Info("Client updated successfully", "id", client.ID)
	return nil
}

func (r *clientRepository) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domainerrors.ErrInvalidClientID
	}

	r.logger.Debug("Deleting client", "id", id)

	result := r.db.WithContext(ctx).Delete(&models.Client{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete client", "id", id, "error", result.Error)
		return fmt.Errorf("failed to delete client: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Client not found for deletion", "id", id)
		return domainerrors.ErrClientNotFound
	}

	r.logger.Info("Client deleted successfully", "id", id)
	return nil
}

func (r *clientRepository) List(ctx context.Context, limit, offset int) ([]models.Client, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	r.logger.Debug("Listing clients", "limit", limit, "offset", offset)

	var clients []models.Client
	result := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&clients)

	if result.Error != nil {
		r.logger.Error("Failed to list clients", "error", result.Error)
		return nil, fmt.Errorf("failed to list clients: %w", result.Error)
	}

	r.logger.Info("Clients listed successfully", "count", len(clients))
	return clients, nil
}
