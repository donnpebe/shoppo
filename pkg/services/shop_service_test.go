package services

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/donnpebe/shoppo/pkg/domain"
	"github.com/donnpebe/shoppo/pkg/lib/promotioncondition/mock"
)

func TestShopService_CreateCart(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should return new order with new generated id",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			orderStore := make(map[string]*domain.Order)
			sut := NewShopService(nil, nil, orderStore)

			got := sut.CreateCart()
			assert.NotNil(t, got)
			assert.NotEmpty(t, got.ID)
			gotInStore, ok := orderStore[got.ID]
			assert.True(t, ok)
			assert.Equal(t, got.ID, gotInStore.ID)
		})
	}
}

func TestShopService_ListProducts(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should return all products in inventories",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inventories := map[string]*domain.Product{
				"p01": {
					ID:        "p01",
					SKU:       "120P90",
					Name:      "Google Home",
					UnitPrice: 49.99,
					Quantity:  10,
				},
				"p02": {
					ID:        "p02",
					SKU:       "43N23P",
					Name:      "MacBook Pro",
					UnitPrice: 5399.99,
					Quantity:  5,
				},
			}

			sut := NewShopService(inventories, nil, nil)
			got := sut.ListProducts()
			assert.Len(t, got, 2)

			for _, g := range got {
				assert.Equal(t, inventories[g.ID], g)
			}
		})
	}
}

func TestShopService_AddItemToCart(t *testing.T) {
	inventories := map[string]*domain.Product{
		"p01": {
			ID:        "p01",
			SKU:       "120P90",
			Name:      "Google Home",
			UnitPrice: 49.99,
			Quantity:  5,
		},
		"p02": {
			ID:        "p02",
			SKU:       "43N23P",
			Name:      "Macbook Pro",
			UnitPrice: 5399.99,
			Quantity:  4,
		},
	}

	type args struct {
		productID string
		quantity  int
	}
	type want struct {
		orderLinesLength int
		lines            []args
	}

	tests := []struct {
		name    string
		input   []args
		want    want
		wantErr error
	}{
		{
			name: "should return order with updated order lines",
			input: []args{
				{
					productID: "p01",
					quantity:  2,
				},
			},
			want: want{
				orderLinesLength: 1,
				lines: []args{
					{
						productID: "p01",
						quantity:  2,
					},
				},
			},
		},
		{
			name: "when the same product added twice, it will only update the quantity of existing line",
			input: []args{
				{
					productID: "p01",
					quantity:  2,
				},
				{
					productID: "p01",
					quantity:  1,
				},
			},
			want: want{
				orderLinesLength: 1,
				lines: []args{
					{
						productID: "p01",
						quantity:  3,
					},
				},
			},
		},
		{
			name: "should return different product per order line when different product added to cart",
			input: []args{
				{
					productID: "p01",
					quantity:  2,
				},
				{
					productID: "p02",
					quantity:  1,
				},
			},
			want: want{
				orderLinesLength: 2,
				lines: []args{
					{
						productID: "p01",
						quantity:  2,
					},
					{
						productID: "p02",
						quantity:  1,
					},
				},
			},
		},
		{
			name: "should return error when product id not in inventories",
			input: []args{
				{
					productID: "p04",
					quantity:  2,
				},
			},
			wantErr: domain.ErrProductNotFound,
		},
		{
			name: "should return error when adding a product with quantity more than what's in inventories",
			input: []args{
				{
					productID: "p01",
					quantity:  6,
				},
			},
			wantErr: domain.ErrNotEnoughStock,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			orderStore := make(map[string]*domain.Order)

			sut := NewShopService(inventories, nil, orderStore)
			order := sut.CreateCart()
			var (
				got *domain.Order
				err error
			)

			if test.wantErr != nil {
				for _, arg := range test.input {
					_, err = sut.AddItemToCart(order.ID, arg.productID, arg.quantity)
					assert.ErrorIs(t, err, test.wantErr)
				}
			} else {

				for _, arg := range test.input {
					got, err = sut.AddItemToCart(order.ID, arg.productID, arg.quantity)
					assert.NoError(t, err)
				}
				assert.Len(t, got.Lines, test.want.orderLinesLength)
				assert.NotEqual(t, got.Lines[0].ID, "")
				for idx, wantLine := range test.want.lines {

					assert.Equal(t, got.Lines[idx].Quantity, wantLine.quantity)
					assert.Equal(t, got.Lines[idx].ProductID, wantLine.productID)
				}
			}
		})
	}
}

func TestShopService_RemoveItemFromCart(t *testing.T) {
	type args struct {
		orderID   string
		productID string
	}

	tests := []struct {
		name    string
		input   args
		wantErr error
	}{
		{
			name: "should return order with updated order lines",
			input: args{
				orderID:   "order1",
				productID: "p01",
			},
		},
		{
			name: "should return error when trying to remove product that is not in order",
			input: args{
				orderID:   "order1",
				productID: "p04",
			},
			wantErr: domain.ErrItemNotFoundInCart,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			orderStore := map[string]*domain.Order{
				"order1": {
					Lines: []*domain.OrderLine{
						{
							ID:        "line1",
							ProductID: "p01",
							Quantity:  2,
							UnitPrice: 10,
						},
						{
							ID:        "line1",
							ProductID: "p02",
							Quantity:  5,
							UnitPrice: 40,
						},
					},
				},
			}

			sut := NewShopService(nil, nil, orderStore)

			order, err := sut.RemoveItemFromCart(test.input.orderID, test.input.productID)

			if test.wantErr != nil {
				assert.ErrorIs(t, err, test.wantErr)
			} else {
				assert.Len(t, order.Lines, 1)
				assert.NoError(t, err)
			}
		})
	}
}

func TestShopService_Checkout(t *testing.T) {
	type item struct {
		productID string
		quantity  int
	}

	type args struct {
		useInvalidOrderID bool
		items             []item
		promotionsCount   int
	}

	type mockBehavior func(ms ...*mock.MockPromotionCondition)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    float64
		wantErr error
	}{
		{
			name: "should return full total amount if promotions not provided",
			input: args{promotionsCount: 0, items: []item{
				{
					productID: "p01",
					quantity:  1,
				},
				{
					productID: "p02",
					quantity:  1,
				},
			}},
			want: 59.99,
		},
		{
			name: "should return discounted total amount if promotion provided",
			input: args{promotionsCount: 1, items: []item{
				{
					productID: "p01",
					quantity:  1,
				},
				{
					productID: "p02",
					quantity:  1,
				},
			}},
			mock: func(ms ...*mock.MockPromotionCondition) {
				ms[0].EXPECT().CalculateDiscount(gomock.Any()).Return(-10.0)
			},
			want: 49.99,
		},
		{
			name: "should return discounted total amount if multiple promotion provided",
			input: args{promotionsCount: 2, items: []item{
				{
					productID: "p01",
					quantity:  1,
				},
				{
					productID: "p02",
					quantity:  1,
				},
			}},
			mock: func(ms ...*mock.MockPromotionCondition) {
				for _, m := range ms {
					m.EXPECT().CalculateDiscount(gomock.Any()).Return(-10.0)
				}
			},
			want: 39.99,
		},
		{
			name:    "should return error if provided with invalid order id",
			input:   args{useInvalidOrderID: true},
			wantErr: domain.ErrCartNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inventories := map[string]*domain.Product{
				"p01": {
					ID:        "p01",
					SKU:       "120P90",
					Name:      "Google Home",
					UnitPrice: 49.99,
					Quantity:  5,
				},
				"p02": {
					ID:        "p02",
					SKU:       "43N23P",
					Name:      "Pen",
					UnitPrice: 10,
					Quantity:  4,
				},
			}

			var promotions []domain.Promotion

			if test.input.promotionsCount > 0 {
				c := gomock.NewController(t)
				defer c.Finish()

				conds := make([]*mock.MockPromotionCondition, 0, test.input.promotionsCount)
				for i := 0; i < test.input.promotionsCount; i++ {
					cond := mock.NewMockPromotionCondition(c)
					conds = append(conds, cond)
					promotions = append(promotions, domain.Promotion{Condition: cond})
				}

				test.mock(conds...)
			}

			orderStore := make(map[string]*domain.Order)
			sut := NewShopService(inventories, promotions, orderStore)
			order := sut.CreateCart()

			for _, item := range test.input.items {
				sut.AddItemToCart(order.ID, item.productID, item.quantity)
			}

			orderID := order.ID
			if test.input.useInvalidOrderID {
				orderID = "invalid"
			}
			totalAmount, err := sut.Checkout(orderID)

			if test.wantErr != nil {
				assert.ErrorIs(t, err, test.wantErr)
			} else {

				assert.NoError(t, err)
				assert.Equal(t, test.want, totalAmount)
			}
		})
	}
}
