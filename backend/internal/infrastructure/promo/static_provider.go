package promo

import (
	"context"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type staticPromoProvider struct {
	promoActions []models.PromoAction
	cards        []models.PromoAction
	investment   []models.PromoAction
}

func NewStaticPromoProvider() interfaces.PromoProvider {
	return &staticPromoProvider{
		promoActions: []models.PromoAction{
			{MinIncome: 0, MaxIncome: 30000, Promo: "Дарим 500р за отзыв!"},
			{MinIncome: 30000, MaxIncome: 60000, Promo: "Бигфест с кэшбэком 50%!"},
			{MinIncome: 60000, MaxIncome: 120000, Promo: "Пятничный суперкэшбек!"},
			{MinIncome: 120000, MaxIncome: 250000, Promo: "Дополнительная категория кэшбека!"},
			{MinIncome: 250000, MaxIncome: 500000, Promo: "Повышенная ставка по накопительному счёту альфа-банка"},
			{MinIncome: 500000, MaxIncome: 1000000, Promo: "Счёт для бизнеса за 0 рублей!"},
			{MinIncome: 1000000, MaxIncome: 1000000000, Promo: "Дарим платёжное кольцо"},
		},
		cards: []models.PromoAction{
			{MinIncome: 0, MaxIncome: 30000, Promo: "Альфа-Стикер"},
			{MinIncome: 30000, MaxIncome: 60000, Promo: "Сверхтонкий стикер с котами"},
			{MinIncome: 60000, MaxIncome: 120000, Promo: "Карта Альфа и Золотое Яблоко"},
			{MinIncome: 120000, MaxIncome: 250000, Promo: "Карта Альфа и Баста"},
			{MinIncome: 250000, MaxIncome: 1000000, Promo: "Дебетовая карта Alfa Only Аэрофлот"},
			{MinIncome: 1000000, MaxIncome: 1000000000, Promo: "Дебетовая карта Alfa Only"},
		},
		investment: []models.PromoAction{
			{MinIncome: 0, MaxIncome: 30000, Promo: "Платим 5 000 ₽ каждому"},
			{MinIncome: 30000, MaxIncome: 60000, Promo: "Старт в инвестициях всего со 100 рублей!"},
			{MinIncome: 60000, MaxIncome: 120000, Promo: "Платим 10 000 ₽ каждому"},
			{MinIncome: 120000, MaxIncome: 250000, Promo: "Акции в подарок"},
			{MinIncome: 250000, MaxIncome: 1000000, Promo: "Инвест-копилка"},
			{MinIncome: 1000000, MaxIncome: 1000000000, Promo: "Тарифный план: Персональный брокер"},
		},
	}
}

func (p *staticPromoProvider) GetPromos(ctx context.Context, income int64, predictIncome float64) ([]string, error) {
	result := make([]string, 0, 3)

	for _, promo := range p.promoActions {
		if income >= promo.MinIncome && income < promo.MaxIncome {
			result = append(result, promo.Promo)
			break
		}
	}

	for _, card := range p.cards {
		if income >= card.MinIncome && income < card.MaxIncome {
			result = append(result, card.Promo)
			break
		}
	}

	for _, investment := range p.investment {
		if income >= investment.MinIncome && income < investment.MaxIncome {
			result = append(result, investment.Promo)
			break
		}
	}

	return result, nil
}
