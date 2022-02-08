package promotioncondition

import "github.com/donnpebe/shoppo/pkg/domain"

type BuyXProductGetFreeProductCondition struct {
	XProductID    string
	FreeProductID string
}

func (cond BuyXProductGetFreeProductCondition) CalculateDiscount(order *domain.Order) float64 {
	promoProductQuantity := 0
	var freeProductLine *domain.OrderLine
	for _, line := range order.Lines {
		if line.ProductID == cond.XProductID {
			promoProductQuantity += line.Quantity
		}

		if line.ProductID == cond.FreeProductID {
			freeProductLine = line
		}
	}

	if freeProductLine == nil {
		return 0
	}

	if freeProductLine.Quantity >= promoProductQuantity {
		return -float64(promoProductQuantity) * freeProductLine.UnitPrice
	}

	return -float64(freeProductLine.Quantity) * freeProductLine.UnitPrice
}
