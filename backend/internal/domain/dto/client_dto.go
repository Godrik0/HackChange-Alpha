package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
	"gorm.io/datatypes"
)

const DateFormat = "2006-01-02"

type CreateClientRequest struct {
	FirstName string                 `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string                 `json:"last_name" validate:"required,min=1,max=100"`
	BirthDate string                 `json:"birth_date" validate:"required"`
	Features  map[string]interface{} `json:"features,omitempty"`
}

type UpdateClientRequest struct {
	FirstName string                 `json:"first_name" validate:"omitempty,min=1,max=100"`
	LastName  string                 `json:"last_name" validate:"omitempty,min=1,max=100"`
	BirthDate string                 `json:"birth_date" validate:"required,datetime=2006-01-02"`
	Features  map[string]interface{} `json:"features,omitempty"`
}

type ClientResponse struct {
	ID        int64                  `json:"id"`
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	BirthDate string                 `json:"birth_date"`
	Features  map[string]interface{} `json:"features,omitempty"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

type SearchParams struct {
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	BirthDate string `json:"birth_date" form:"birth_date"`
}

func (s SearchParams) IsEmpty() bool {
	return s.FirstName == "" && s.LastName == "" && s.BirthDate == ""
}

func (r *CreateClientRequest) ToModel() (*models.Client, error) {
	parsedBirthDate, err := time.Parse(DateFormat, r.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("invalid birth_date format (expected %s): %w", DateFormat, err)
	}

	client := &models.Client{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		BirthDate: parsedBirthDate,
	}

	if r.Features != nil {
		featuresJSON, err := json.Marshal(r.Features)
		if err != nil {
			return nil, err
		}
		client.Features = datatypes.JSON(featuresJSON)
	}

	return client, nil
}

func FromModel(client *models.Client) (*ClientResponse, error) {
	response := &ClientResponse{
		ID:        client.ID,
		FirstName: client.FirstName,
		LastName:  client.LastName,
		BirthDate: client.BirthDate.Format(DateFormat),
		CreatedAt: client.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: client.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if len(client.Features) > 0 {
		var features map[string]interface{}
		if err := json.Unmarshal(client.Features, &features); err == nil {
			response.Features = features
		}
	}

	return response, nil
}

func FromModels(clients []models.Client) ([]*ClientResponse, error) {
	responses := make([]*ClientResponse, 0, len(clients))
	for _, client := range clients {
		response, err := FromModel(&client)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}
