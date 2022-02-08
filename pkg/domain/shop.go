package domain

type ShopService interface {
	CreateCart() *Order
	ListProducts() []*Product
	AddItemToCart(orderID string, productID string, quantity int) (*Order, error)
	RemoveItemFromCart(orderID string, productID string) (*Order, error)
	Checkout(orderID string) (totalAmount float64, err error)
}
