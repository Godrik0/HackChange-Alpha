package mocks

import (
	"context"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

// MockClientRepository мок репозитория клиентов для тестов
type MockClientRepository struct {
	CreateFunc  func(ctx context.Context, client *models.Client) error
	GetByIDFunc func(ctx context.Context, id int64) (*models.Client, error)
	SearchFunc  func(ctx context.Context, params dto.SearchParams) ([]models.Client, error)
	UpdateFunc  func(ctx context.Context, client *models.Client) error
	DeleteFunc  func(ctx context.Context, id int64) error
	ListFunc    func(ctx context.Context, limit, offset int) ([]models.Client, error)
}

func (m *MockClientRepository) Create(ctx context.Context, client *models.Client) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, client)
	}
	return nil
}

func (m *MockClientRepository) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockClientRepository) Search(ctx context.Context, params dto.SearchParams) ([]models.Client, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, params)
	}
	return []models.Client{}, nil
}

func (m *MockClientRepository) Update(ctx context.Context, client *models.Client) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, client)
	}
	return nil
}

func (m *MockClientRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockClientRepository) List(ctx context.Context, limit, offset int) ([]models.Client, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return []models.Client{}, nil
}
