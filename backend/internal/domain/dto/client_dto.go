package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
	"gorm.io/datatypes"
)

const DateFormat = "02-01-2006"

type CreateClientRequest struct {
	FirstName  string                 `json:"first_name" validate:"required,min=1,max=100"`
	LastName   string                 `json:"last_name" validate:"required,min=1,max=100"`
	MiddleName string                 `json:"middle_name" validate:"omitempty,max=100"`
	BirthDate  string                 `json:"birth_date" validate:"required"`
	Features   map[string]interface{} `json:"features,omitempty"`
}

type UpdateClientRequest struct {
	FirstName  string                 `json:"first_name" validate:"omitempty,min=1,max=100"`
	LastName   string                 `json:"last_name" validate:"omitempty,min=1,max=100"`
	MiddleName string                 `json:"middle_name" validate:"omitempty,max=100"`
	BirthDate  string                 `json:"birth_date" validate:"omitempty,datetime=02-01-2006"`
	Features   map[string]interface{} `json:"features,omitempty"`
}

type ClientResponse struct {
	ID         int64                  `json:"id"`
	FirstName  string                 `json:"first_name"`
	LastName   string                 `json:"last_name"`
	MiddleName string                 `json:"middle_name,omitempty"`
	BirthDate  string                 `json:"birth_date"`
	Features   map[string]interface{} `json:"features,omitempty"`
}

type ClientSearchResponse struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
	BirthDate  string `json:"birth_date"`
}

type SearchParams struct {
	FirstName  string `json:"first_name" form:"first_name"`
	LastName   string `json:"last_name" form:"last_name"`
	MiddleName string `json:"middle_name" form:"middle_name"`
	BirthDate  string `json:"birth_date" form:"birth_date"`
}

func (s SearchParams) IsEmpty() bool {
	return s.FirstName == "" && s.LastName == "" && s.MiddleName == "" && s.BirthDate == ""
}

func (r *CreateClientRequest) ToModel() (*models.Client, error) {
	parsedBirthDate, err := time.Parse(DateFormat, r.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("invalid birth_date format (expected %s): %w", DateFormat, err)
	}

	client := &models.Client{
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		MiddleName: r.MiddleName,
		BirthDate:  parsedBirthDate,
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
		ID:         client.ID,
		FirstName:  client.FirstName,
		LastName:   client.LastName,
		MiddleName: client.MiddleName,
		BirthDate:  client.BirthDate.Format(DateFormat),
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

func FromModelToSearchResponse(client *models.Client) *ClientSearchResponse {
	return &ClientSearchResponse{
		ID:         client.ID,
		FirstName:  client.FirstName,
		LastName:   client.LastName,
		MiddleName: client.MiddleName,
		BirthDate:  client.BirthDate.Format(DateFormat),
	}
}

func FromModelsToSearchResponse(clients []models.Client) []*ClientSearchResponse {
	responses := make([]*ClientSearchResponse, 0, len(clients))
	for _, client := range clients {
		responses = append(responses, FromModelToSearchResponse(&client))
	}
	return responses
}
