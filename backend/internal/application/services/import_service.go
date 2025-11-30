package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type ImportService struct {
	clientRepo interfaces.ClientRepository
	logger     interfaces.Logger
	batchSize  int
}

func NewImportService(clientRepo interfaces.ClientRepository, logger interfaces.Logger) interfaces.ImportService {
	return &ImportService{
		clientRepo: clientRepo,
		logger:     logger.With("component", "ImportService"),
		batchSize:  500,
	}
}

func (s *ImportService) ImportClientsCSV(ctx context.Context, reader io.Reader) (*interfaces.ImportStats, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = -1
	csvReader.LazyQuotes = true

	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	stats := &interfaces.ImportStats{}
	var clientsBatch []*models.Client
	lineNum := 1

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		lineNum++

		if err != nil {
			s.logger.Warn("Failed to read CSV line", "line", lineNum, "error", err)
			stats.AddError(lineNum, err)
			continue
		}

		rowData := s.makeRowMap(headers, record)

		client, err := s.parseClientFromCSVRow(rowData, headers)
		if err != nil {
			stats.AddError(lineNum, err)
			continue
		}

		clientsBatch = append(clientsBatch, client)

		if len(clientsBatch) >= s.batchSize {
			if err := s.insertBatch(ctx, clientsBatch, stats); err != nil {
				s.logger.Error("Failed to insert batch", "error", err)
				return nil, fmt.Errorf("failed to insert batch: %w", err)
			}
			clientsBatch = clientsBatch[:0]
		}
	}

	if len(clientsBatch) > 0 {
		if err := s.insertBatch(ctx, clientsBatch, stats); err != nil {
			s.logger.Error("Failed to insert final batch", "error", err)
			return nil, fmt.Errorf("failed to insert final batch: %w", err)
		}
	}

	stats.Total = stats.SuccessCount + stats.FailureCount
	s.logger.Info("CSV import completed", "success", stats.SuccessCount, "failures", stats.FailureCount)

	return stats, nil
}

func (s *ImportService) insertBatch(ctx context.Context, clients []*models.Client, stats *interfaces.ImportStats) error {
	created, err := s.clientRepo.BatchCreate(ctx, clients)
	if err != nil {
		s.logger.Warn("Batch insert failed, falling back to individual inserts", "error", err)
		for _, client := range clients {
			if err := s.clientRepo.Create(ctx, client); err != nil {
				stats.AddError(0, fmt.Errorf("failed to create client %s %s: %w", client.FirstName, client.LastName, err))
			} else {
				stats.SuccessCount++
			}
		}
		return nil
	}

	stats.SuccessCount += created
	return nil
}

func (s *ImportService) makeRowMap(headers []string, record []string) map[string]string {
	rowData := make(map[string]string)
	for i, header := range headers {
		if i < len(record) {
			rowData[header] = strings.TrimSpace(record[i])
		} else {
			rowData[header] = ""
		}
	}
	return rowData
}

func (s *ImportService) parseClientFromCSVRow(row map[string]string, headers []string) (*models.Client, error) {
	firstName, ok := row["first_name"]
	if !ok || firstName == "" {
		return nil, errors.New("first_name is required")
	}

	lastName, ok := row["last_name"]
	if !ok || lastName == "" {
		return nil, errors.New("last_name is required")
	}

	birthDateStr, ok := row["birth_date"]
	if !ok || birthDateStr == "" {
		return nil, errors.New("birth_date is required")
	}

	birthDate, err := s.parseBirthDate(birthDateStr)
	if err != nil {
		return nil, err
	}

	client := &models.Client{
		FirstName:  firstName,
		LastName:   lastName,
		MiddleName: row["middle_name"],
		BirthDate:  birthDate,
	}

	features := s.extractFeatures(row, headers)
	if len(features) > 0 {
		featuresJSON, err := json.Marshal(features)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal features: %w", err)
		}
		client.Features = featuresJSON
	}

	return client, nil
}

func (s *ImportService) parseBirthDate(dateStr string) (time.Time, error) {
	formats := []string{"02-01-2006", "2006-01-02", "02/01/2006", "2006/01/02"}

	for _, format := range formats {
		if birthDate, err := time.Parse(format, dateStr); err == nil {
			return birthDate, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid birth_date format: %s (expected DD-MM-YYYY or YYYY-MM-DD)", dateStr)
}

func (s *ImportService) extractFeatures(row map[string]string, headers []string) map[string]interface{} {
	features := make(map[string]interface{})

	baseFields := map[string]bool{
		"first_name":  true,
		"last_name":   true,
		"middle_name": true,
		"birth_date":  true,
		"phone":       true,
		"email":       true,
		"address":     true,
		"user_id":     true,
	}

	for _, header := range headers {
		if baseFields[header] {
			continue
		}

		value := row[header]

		if value == "" || value == "nan" || value == "NaN" || value == "null" || value == "None" {
			features[header] = 0.0
			continue
		}

		value = strings.ReplaceAll(value, ",", ".")

		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			features[header] = floatVal
		} else {
			features[header] = value
		}
	}

	return features
}
