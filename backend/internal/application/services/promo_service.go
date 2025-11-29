package services

import (
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
)

type PromoService struct {
	promoData models.PromoCategory
}

func NewPromoService() *PromoService {
	return &PromoService{
		promoData: models.PromoCategory{
			PromoActions: []models.PromoAction{
				{MinIncome: 0, MaxIncome: 30000, Promo: "Дарим 500р за отзыв!"},
				{MinIncome: 30000, MaxIncome: 60000, Promo: "Бигфест с кэшбэком 50%!"},
				{MinIncome: 60000, MaxIncome: 120000, Promo: "Пятничный суперкэшбек!"},
				{MinIncome: 120000, MaxIncome: 250000, Promo: "Дополнительная категория кэшбека!"},
				{MinIncome: 250000, MaxIncome: 500000, Promo: "Повышенная ставка по накопительному счёту альфа-банка"},
				{MinIncome: 500000, MaxIncome: 1000000, Promo: "Счёт для бизнеса за 0 рублей!"},
				{MinIncome: 1000000, MaxIncome: 1000000000, Promo: "Дарим платёжное кольцо"},
			},
			Cards: []models.PromoAction{
				{MinIncome: 0, MaxIncome: 30000, Promo: "Альфа-Стикер"},
				{MinIncome: 30000, MaxIncome: 60000, Promo: "Сверхтонкий стикер с котами"},
				{MinIncome: 60000, MaxIncome: 120000, Promo: "Карта Альфа и Золотое Яблоко"},
				{MinIncome: 120000, MaxIncome: 250000, Promo: "Карта Альфа и Баста"},
				{MinIncome: 250000, MaxIncome: 1000000, Promo: "Дебетовая карта Alfa Only Аэрофлот"},
				{MinIncome: 1000000, MaxIncome: 1000000000, Promo: "Дебетовая карта Alfa Only"},
			},
			Investment: []models.PromoAction{
				{MinIncome: 0, MaxIncome: 30000, Promo: "Платим 5 000 ₽ каждому"},
				{MinIncome: 30000, MaxIncome: 60000, Promo: "Старт в инвестициях всего со 100 рублей!"},
				{MinIncome: 60000, MaxIncome: 120000, Promo: "Платим 10 000 ₽ каждому"},
				{MinIncome: 120000, MaxIncome: 250000, Promo: "Акции в подарок"},
				{MinIncome: 250000, MaxIncome: 1000000, Promo: "Инвест-копилка"},
				{MinIncome: 1000000, MaxIncome: 1000000000, Promo: "Тарифный план: Персональный брокер"},
			},
		},
	}
}

func (s *PromoService) GetPromoByIncome(income int64) []string {
	result := make([]string, 0, 3)

	for _, promo := range s.promoData.PromoActions {
		if income >= promo.MinIncome && income < promo.MaxIncome {
			result = append(result, promo.Promo)
			break
		}
	}

	for _, card := range s.promoData.Cards {
		if income >= card.MinIncome && income < card.MaxIncome {
			result = append(result, card.Promo)
			break
		}
	}

	for _, investment := range s.promoData.Investment {
		if income >= investment.MinIncome && income < investment.MaxIncome {
			result = append(result, investment.Promo)
			break
		}
	}

	return result
}

func (s *PromoService) GetAllPromos() models.PromoCategory {
	return s.promoData
}
