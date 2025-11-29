package models

type PromoAction struct {
	MinIncome int64  `json:"min_income"`
	MaxIncome int64  `json:"max_income"`
	Promo     string `json:"promo"`
}

type PromoCategory struct {
	PromoActions []PromoAction `json:"promo_actions"`
	Cards        []PromoAction `json:"cards"`
	Investment   []PromoAction `json:"investment"`
}
