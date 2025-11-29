package models

type ScoringResult struct {
	Score           float64            `json:"score"`
	Recommendations []string           `json:"recommendations"`
	Factors         map[string]float64 `json:"factors"`
	PositiveFactors []string           `json:"positive_factors"`
	NegativeFactors []string           `json:"negative_factors"`
	CreditLimit     float64            `json:"credit_limit"`
	MaxCreditLimit  float64            `json:"max_credit_limit"`
}

func (s *ScoringResult) IsValid() bool {
	return s.Score >= 0
}
