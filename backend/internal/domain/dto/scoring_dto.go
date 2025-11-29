package dto

type ScoringResponse struct {
	Score           float64            `json:"score"`
	Recommendations []string           `json:"recommendations"`
	Factors         map[string]float64 `json:"factors"`
}
