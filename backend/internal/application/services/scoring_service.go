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
	clientRepo    interfaces.ClientRepository
	mlService     interfaces.MLService
	creditCalc    *CreditLimitCalculator
	promoProvider interfaces.PromoProvider
	logger        interfaces.Logger
}

func NewScoringService(
	clientRepo interfaces.ClientRepository,
	mlService interfaces.MLService,
	promoProvider interfaces.PromoProvider,
	logger interfaces.Logger,
) interfaces.ScoringService {
	return &scoringService{
		clientRepo:    clientRepo,
		mlService:     mlService,
		creditCalc:    NewCreditLimitCalculator(),
		promoProvider: promoProvider,
		logger:        logger.With("component", "ScoringService"),
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

	mlResponse, err := s.mlService.PredictWithExplanation(ctx, features)
	if err != nil {
		s.logger.Error("Failed to predict scoring", "client_id", id, "error", err)
		return nil, fmt.Errorf("failed to predict scoring: %w", err)
	}

	creditLimit := s.calculateCreditLimit(features, mlResponse.Prediction)
	recommendations := s.getRecommendations(ctx, mlResponse.Prediction)
	positiveFactors, negativeFactors := s.splitFactorsBySign(mlResponse.Explanation)

	result := &models.ScoringResult{
		Score:           mlResponse.Prediction,
		Recommendations: recommendations,
		Factors: map[string]float64{
			"predicted_income": mlResponse.Prediction,
			"credit_limit":     creditLimit.RecommendationCreditLimit,
		},
		PositiveFactors: FormatPositiveFactors(positiveFactors),
		NegativeFactors: FormatNegativeFactors(negativeFactors),
		CreditLimit:     creditLimit.RecommendationCreditLimit,
		MaxCreditLimit:  creditLimit.LimitLegal,
	}

	s.logger.Info("Scoring calculated successfully", "client_id", id, "score", result.Score)
	return result, nil
}

func (s *scoringService) calculateCreditLimit(features map[string]interface{}, predictedIncome float64) dto.CreditLimitResult {
	creditLimitInput := s.extractCreditLimitInput(features, predictedIncome)
	return s.creditCalc.Calculate(creditLimitInput)
}

func (s *scoringService) getRecommendations(ctx context.Context, prediction float64) []string {
	predictedIncome := int64(prediction)
	recommendations, err := s.promoProvider.GetPromos(ctx, predictedIncome, prediction)
	if err != nil {
		s.logger.Error("Failed to get promos", "error", err)
		return []string{}
	}
	return recommendations
}

func (s *scoringService) splitFactorsBySign(explanation map[string]float64) (positive, negative map[string]float64) {
	positive = make(map[string]float64)
	negative = make(map[string]float64)

	for key, value := range explanation {
		if value > 0 {
			positive[key] = value
		} else if value < 0 {
			negative[key] = value
		}
	}

	return positive, negative
}

func (s *scoringService) extractCreditLimitInput(features map[string]interface{}, predictedIncome float64) dto.CreditLimitInput {
	getFloat := func(key string) float64 {
		if val, ok := features[key]; ok {
			switch v := val.(type) {
			case float64:
				return v
			case int:
				return float64(v)
			case int64:
				return float64(v)
			}
		}
		return 0.0
	}

	getInt := func(key string) int {
		if val, ok := features[key]; ok {
			switch v := val.(type) {
			case int:
				return v
			case float64:
				return int(v)
			}
		}
		return 0
	}

	return dto.CreditLimitInput{
		PredictedIncome:        predictedIncome,
		ActiveCCMaxLimit:       getFloat("hdb_bki_active_cc_max_limit"),
		OutstandSum:            getFloat("hdb_outstand_sum"),
		OverdueSum:             getFloat("ovrd_sum"),
		BlacklistFlag:          getInt("blacklist_flag"),
		TurnCurrentCreditAvgV2: getFloat("turn_cur_cr_avg_v2"),
	}
}

func (s *scoringService) extractFeatures(client *models.Client) (map[string]interface{}, error) {
	features := make(map[string]interface{})

	features["user_id"] = fmt.Sprintf("%d", client.ID)
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
