package dto

type MLScoringResponse struct {
	Prediction  float64            `json:"prediction"`
	Explanation map[string]map[string]float64 `json:"explanation"`
	ID          string             `json:"id"`
}

type CreditLimitInput struct {
	PredictedIncome        float64
	ActiveCCMaxLimit       float64
	OutstandSum            float64
	OverdueSum             float64
	BlacklistFlag          int
	TurnCurrentCreditAvgV2 float64
}

type CreditLimitResult struct {
	LimitLegal                float64 `json:"limit_legal"`
	RecommendationCreditLimit float64 `json:"recommendation_credit_limit"`
}
