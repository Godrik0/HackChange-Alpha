package dto

type ScoringResponse struct {
	Id                        int64    `json:"id"`
	FirstName                 string   `json:"first_name"`
	LastName                  string   `json:"last_name"`
	MiddleName                string   `json:"middle_name,omitempty"`
	BirthDate                 string   `json:"birth_date"`
	Income                    int64    `json:"income,omitempty"`
	PredictIncome             float64  `json:"predict_income"`
	RecommendationCreditLimit float64  `json:"credit_limit"`
	MaxCreditLimit            float64  `json:"max_credit_limit,omitempty"`
	Recommendations           []string `json:"recommendations"`
	PositiveFactors           []string `json:"positive_factors"`
	NegativeFactors           []string `json:"negative_factors"`
}
