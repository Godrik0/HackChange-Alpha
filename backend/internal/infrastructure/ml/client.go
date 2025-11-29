package ml

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type mlClient struct {
	baseURL    string
	httpClient *http.Client
	logger     interfaces.Logger
}

func NewMLClient(baseURL string, timeoutSeconds int, logger interfaces.Logger) interfaces.MLService {
	return &mlClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
		logger: logger.With("component", "MLClient"),
	}
}

type predictRequest struct {
	Features map[string]interface{} `json:"features"`
}

type predictResponse struct {
	Score           float64            `json:"score"`
	Recommendations []string           `json:"recommendations"`
	Factors         map[string]float64 `json:"factors"`
}

func (c *mlClient) Predict(ctx context.Context, features map[string]interface{}) (*models.ScoringResult, error) {
	if features == nil || len(features) == 0 {
		return nil, fmt.Errorf("features cannot be empty")
	}

	c.logger.Debug("Requesting ML prediction", "features_count", len(features))

	reqBody := predictRequest{Features: features}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Error("Failed to marshal prediction request", "error", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/predict", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create HTTP request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to execute ML prediction request", "error", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("ML service returned error", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("ML service returned status %d: %s", resp.StatusCode, string(body))
	}

	var response predictResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode ML response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := &models.ScoringResult{
		Score:           response.Score,
		Recommendations: response.Recommendations,
		Factors:         response.Factors,
	}

	c.logger.Info("ML prediction completed", "score", result.Score)
	return result, nil
}

func (c *mlClient) SendTrainingData(ctx context.Context, data interface{}) error {
	if data == nil {
		return fmt.Errorf("training data cannot be nil")
	}

	c.logger.Debug("Sending training data to ML service")

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("Failed to marshal training data", "error", err)
		return fmt.Errorf("failed to marshal training data: %w", err)
	}

	url := fmt.Sprintf("%s/api/train", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create training request", "error", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to send training data", "error", err)
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("ML service rejected training data", "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("ML service returned status %d: %s", resp.StatusCode, string(body))
	}

	c.logger.Info("Training data sent successfully")
	return nil
}

func (c *mlClient) HealthCheck(ctx context.Context) error {
	c.logger.Debug("Checking ML service health")

	url := fmt.Sprintf("%s/health", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Warn("ML service health check failed", "error", err)
		return fmt.Errorf("ML service health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Warn("ML service unhealthy", "status", resp.StatusCode)
		return fmt.Errorf("ML service returned status %d", resp.StatusCode)
	}

	c.logger.Debug("ML service is healthy")
	return nil
}
