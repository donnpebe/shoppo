package promotioncondition

import "github.com/donnpebe/shoppo/pkg/domain"

type ProductQuantityDiscount struct {
	ProductID          string
	RequiredQuantity   int
	DiscountedQuantity int
}

func (cond ProductQuantityDiscount) CalculateDiscount(order *domain.Order) float64 {
	for _, line := range order.Lines {
		if line.ProductID == cond.ProductID {
			multiplier := line.Quantity / cond.RequiredQuantity
			return -float64(multiplier*cond.DiscountedQuantity) * float64(line.UnitPrice)
		}
	}

	return 0
}
