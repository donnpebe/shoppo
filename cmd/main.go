package main

import (
	"fmt"
	"log"

	"github.com/donnpebe/shoppo/pkg/domain"
	"github.com/donnpebe/shoppo/pkg/lib/promotioncondition"
	"github.com/donnpebe/shoppo/pkg/services"
)

var inventories = map[string]*domain.Product{
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

func main() {
	promotions := setupPromotion()
	shopService := services.NewShopService(inventories, promotions, make(map[string]*domain.Order))

	fmt.Println("Welcome to shoppo")
	fmt.Println("=================")

	// Scenario 1: Each sale of a Macbook Pro comes with a free Raspberry Pi B
	order := shopService.CreateCart()
	order, err := shopService.AddItemToCart(order.ID, "macbookpro", 1)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}

	order, err = shopService.AddItemToCart(order.ID, "raspberrypi", 1)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}

	totalAmount, err := shopService.Checkout(order.ID)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}

	fmt.Printf("For scenario 1 you need to pay: %.2f\n", totalAmount)

	// Scenario 2: Buy 3 Google Homes for the price of 2
	order2 := shopService.CreateCart()
	order2, err = shopService.AddItemToCart(order2.ID, "googlehome", 1)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}
	order2, err = shopService.AddItemToCart(order2.ID, "googlehome", 1)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}
	order2, err = shopService.AddItemToCart(order2.ID, "googlehome", 1)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}

	totalAmount, err = shopService.Checkout(order2.ID)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}

	fmt.Printf("For scenario 2 you need to pay: %.2f\n", totalAmount)

	// Scenario 3: Buy more than 3 Alexa Speakers will have 10% discount on all alexa speakers
	order3 := shopService.CreateCart()
	order3, err = shopService.AddItemToCart(order3.ID, "alexaspeaker", 1)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}
	order3, err = shopService.AddItemToCart(order3.ID, "alexaspeaker", 1)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}
	order3, err = shopService.AddItemToCart(order3.ID, "alexaspeaker", 1)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}

	totalAmount, err = shopService.Checkout(order3.ID)
	if err != nil {
		log.Fatalf("cannot add item to cart: %v", err)
	}

	fmt.Printf("For scenario 3 you need to pay: %.2f\n", totalAmount)
}

func setupPromotion() []domain.Promotion {
	return []domain.Promotion{
		{
			Condition: promotioncondition.BuyXProductGetFreeProductCondition{
				XProductID:    "macbookpro",
				FreeProductID: "raspberrypi",
			},
		},
		{
			Condition: promotioncondition.ProductQuantityDiscount{
				ProductID:          "googlehome",
				RequiredQuantity:   3,
				DiscountedQuantity: 1,
			},
		},
		{
			Condition: promotioncondition.ProductPercentageDiscount{
				ProductID:         "alexaspeaker",
				MinQuantity:       3,
				DiscountInPercent: 10,
			},
		},
	}
}
