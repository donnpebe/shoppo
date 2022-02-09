package main

import (
	"testing"

	"github.com/donnpebe/shoppo/pkg/domain"
	"github.com/donnpebe/shoppo/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite

	inventories map[string]*domain.Product
	orderStore  map[string]*domain.Order
}

func TestMain(t *testing.T) {
	ms := &MainTestSuite{}
	suite.Run(t, ms)
}

func (ms *MainTestSuite) SetupTest() {
	ms.inventories = map[string]*domain.Product{
		"googlehome": {
			ID:        "googlehome",
			SKU:       "120P90",
			Name:      "Google Home",
			UnitPrice: 49.99,
			Quantity:  10,
		},
		"macbookpro": {
			ID:        "macbookpro",
			SKU:       "43N23P",
			Name:      "MacBook Pro",
			UnitPrice: 5399.99,
			Quantity:  5,
		},
		"alexaspeaker": {
			ID:        "alexaspeaker",
			SKU:       "A304SD",
			Name:      "Alexa Speaker",
			UnitPrice: 109.50,
			Quantity:  10,
		},
		"raspberrypi": {
			ID:        "raspberrypi",
			SKU:       "234234",
			Name:      "Raspberry Pi B",
			UnitPrice: 30,
			Quantity:  2,
		},
	}

	ms.orderStore = make(map[string]*domain.Order)
}

func (ms *MainTestSuite) TestBuyXProductGetFreeProductCondition() {
	sut := services.NewShopService(ms.inventories, setupPromotion(), ms.orderStore)

	order := sut.CreateCart()
	order, err := sut.AddItemToCart(order.ID, "macbookpro", 1)
	assert.NoError(ms.T(), err)
	order, err = sut.AddItemToCart(order.ID, "raspberrypi", 1)
	assert.NoError(ms.T(), err)

	totalAmount, err := sut.Checkout(order.ID)
	assert.NoError(ms.T(), err)
	assert.Equal(ms.T(), 5399.99, totalAmount)
}

func (ms *MainTestSuite) TestProductPercentageDiscountCondition() {
	sut := services.NewShopService(ms.inventories, setupPromotion(), ms.orderStore)

	order := sut.CreateCart()
	order, err := sut.AddItemToCart(order.ID, "alexaspeaker", 1)
	assert.NoError(ms.T(), err)
	order, err = sut.AddItemToCart(order.ID, "alexaspeaker", 1)
	assert.NoError(ms.T(), err)
	order, err = sut.AddItemToCart(order.ID, "alexaspeaker", 1)
	assert.NoError(ms.T(), err)

	totalAmount, err := sut.Checkout(order.ID)
	assert.NoError(ms.T(), err)
	assert.Equal(ms.T(), 295.65, totalAmount)
}

func (ms *MainTestSuite) TestProductQuantityDiscountCondition() {
	sut := services.NewShopService(ms.inventories, setupPromotion(), ms.orderStore)

	order := sut.CreateCart()
	order, err := sut.AddItemToCart(order.ID, "googlehome", 1)
	assert.NoError(ms.T(), err)
	order, err = sut.AddItemToCart(order.ID, "googlehome", 1)
	assert.NoError(ms.T(), err)
	order, err = sut.AddItemToCart(order.ID, "googlehome", 1)
	assert.NoError(ms.T(), err)

	totalAmount, err := sut.Checkout(order.ID)
	assert.NoError(ms.T(), err)
	assert.Equal(ms.T(), 49.99*2, totalAmount)
}
