package promotioncondition

import (
	"testing"

	"github.com/donnpebe/shoppo/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestProductQuantityDiscountCondition_CalculateDiscount(t *testing.T) {
	tests := []struct {
		name  string
		input *domain.Order
		want  float64
	}{
		{
			name: "should return correct discount amount",
			input: &domain.Order{
				Lines: []*domain.OrderLine{
					{
						ID:        "line1",
						ProductID: "p01",
						Quantity:  3,
						UnitPrice: 49.99,
					},
				},
			},
			want: -49.99,
		},
		{
			name: "should return correct discount amount even with multiple of required quantity",
			input: &domain.Order{
				Lines: []*domain.OrderLine{
					{
						ID:        "line1",
						ProductID: "p01",
						Quantity:  3 * 2,
						UnitPrice: 49.99,
					},
				},
			},
			want: -49.99 * 2,
		},
		{
			name: "should return zero if quantity is less than required promo quantity",
			input: &domain.Order{
				Lines: []*domain.OrderLine{
					{
						ID:        "line1",
						ProductID: "p01",
						Quantity:  2,
						UnitPrice: 49.99,
					},
				},
			},
			want: 0,
		},
		{
			name: "should return zero if product id not found in order",
			input: &domain.Order{
				Lines: []*domain.OrderLine{
					{
						ID:        "line1",
						ProductID: "p03",
						Quantity:  3,
						UnitPrice: 109.50,
					},
				},
			},
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sut := &ProductQuantityDiscount{
				ProductID:          "p01",
				RequiredQuantity:   3,
				DiscountedQuantity: 1,
			}

			got := sut.CalculateDiscount(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
