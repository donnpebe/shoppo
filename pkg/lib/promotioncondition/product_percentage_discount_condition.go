package promotioncondition

import "github.com/donnpebe/shoppo/pkg/domain"

type ProductPercentageDiscount struct {
	ProductID         string
	MinQuantity       int
	DiscountInPercent float64
}

func (cond ProductPercentageDiscount) CalculateDiscount(order *domain.Order) float64 {
	for _, line := range order.Lines {
		if line.ProductID == cond.ProductID && line.Quantity >= cond.MinQuantity {
			return -line.UnitPrice * float64(line.Quantity) * (cond.DiscountInPercent / 100)
		}
	}

	return 0
}
