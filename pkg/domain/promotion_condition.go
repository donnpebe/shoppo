package domain

type PromotionCondition interface {
	CalculateDiscount(order *Order) float64
}
