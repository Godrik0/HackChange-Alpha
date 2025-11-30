package interfaces

import "context"

type PromoProvider interface {
	GetPromos(ctx context.Context, income int64, predictIncome float64) ([]string, error)
}
