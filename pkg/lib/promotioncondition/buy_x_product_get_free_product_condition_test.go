package promotioncondition

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/donnpebe/shoppo/pkg/domain"
)

func TestBuyXProductGetFreeProductCondition_CalculateDiscount(t *testing.T) {
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
						ProductID: "p02",
						Quantity:  1,
						UnitPrice: 5399.99,
					},
					{
						ID:        "line2",
						ProductID: "p04",
						Quantity:  1,
						UnitPrice: 30.0,
					},
				},
			},
			want: -30.0,
		},
		{
			name: "should return zero discount if cannot find XProduct id in order",
			input: &domain.Order{
				Lines: []*domain.OrderLine{
					{
						ID:        "line1",
						ProductID: "p03",
						Quantity:  1,
						UnitPrice: 109.50,
					},
					{
						ID:        "line2",
						ProductID: "p04",
						Quantity:  1,
						UnitPrice: 30.0,
					},
				},
			},
			want: 0.0,
		},
		{
			name: "should return zero discount if free product line not included",
			input: &domain.Order{
				Lines: []*domain.OrderLine{
					{
						ID:        "line1",
						ProductID: "p02",
						Quantity:  1,
						UnitPrice: 5399.99,
					},
				},
			},
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sut := BuyXProductGetFreeProductCondition{
				XProductID:    "p02",
				FreeProductID: "p04",
			}

			got := sut.CalculateDiscount(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
