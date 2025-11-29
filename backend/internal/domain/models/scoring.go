package models

type ScoringResult struct {
	Score           float64            `json:"score"`
	Recommendations []string           `json:"recommendations"`
	Factors         map[string]float64 `json:"factors"`
}

func (s *ScoringResult) IsValid() bool {
	return s.Score >= 0 && s.Score <= 1
}
