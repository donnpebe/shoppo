package promotioncondition

import (
	"testing"

	"github.com/donnpebe/shoppo/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestProductPercentageDiscountCondition_CalculateDiscount(t *testing.T) {
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
						ProductID: "p03",
						Quantity:  3,
						UnitPrice: 109.5,
					},
				},
			},
			want: -32.85,
		},
		{
			name: "should return zero if quantity in order less than required min quantity",
			input: &domain.Order{
				Lines: []*domain.OrderLine{
					{
						ID:        "line1",
						ProductID: "p03",
						Quantity:  2,
						UnitPrice: 109.5,
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
						ProductID: "p02",
						Quantity:  3,
						UnitPrice: 5399.99,
					},
				},
			},
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sut := &ProductPercentageDiscount{
				ProductID:         "p03",
				MinQuantity:       3,
				DiscountInPercent: 10,
			}

			got := sut.CalculateDiscount(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
