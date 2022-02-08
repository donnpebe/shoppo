package domain

import "errors"

var (
	ErrCartNotFound                      = errors.New("cart not found")
	ErrProductNotFound                   = errors.New("product not found")
	ErrNotEnoughStock                    = errors.New("not enough stock")
	ErrItemNotFoundInCart                = errors.New("item not found in cart")
	ErrSomeProductInCartNotFound         = errors.New("some product in cart are not found")
	ErrSomeProductInCartNotEnoughInStock = errors.New("some product in cart are not enough in stock")
)
