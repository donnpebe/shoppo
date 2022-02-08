package services

import (
	"sync"
	"time"

	"github.com/rs/xid"

	"github.com/donnpebe/shoppo/pkg/domain"
)

type ShopService struct {
	inventories map[string]*domain.Product
	invMutex    sync.RWMutex

	promotions []domain.Promotion

	orderStore map[string]*domain.Order
}

func NewShopService(inventories map[string]*domain.Product, promotions []domain.Promotion, orderStore map[string]*domain.Order) *ShopService {
	return &ShopService{
		inventories: inventories,
		promotions:  promotions,
		orderStore:  orderStore,
	}
}

func (service *ShopService) CreateCart() *domain.Order {
	order := &domain.Order{
		ID: xid.New().String(),
	}

	service.orderStore[order.ID] = order
	return order
}

func (service *ShopService) ListProducts() []*domain.Product {
	service.invMutex.RLock()

	var products []*domain.Product
	for _, product := range service.inventories {
		products = append(products, product)
	}

	service.invMutex.RUnlock()

	return products
}

func (service *ShopService) AddItemToCart(orderID string, productID string, quantity int) (*domain.Order, error) {
	order, ok := service.orderStore[orderID]
	if !ok {
		return nil, domain.ErrCartNotFound
	}

	service.invMutex.RLock()
	defer service.invMutex.RUnlock()
	product, ok := service.inventories[productID]
	if !ok {
		return nil, domain.ErrProductNotFound
	}

	foundLine, _ := findLineInOrder(order, productID)
	if foundLine == nil {
		if (product.Quantity - quantity) < 0 {
			return nil, domain.ErrNotEnoughStock
		}

		order.Lines = append(order.Lines, &domain.OrderLine{
			ID:        xid.New().String(),
			ProductID: productID,
			Quantity:  quantity,
			UnitPrice: product.UnitPrice,
		})

		return order, nil
	}

	if (product.Quantity - (foundLine.Quantity + quantity)) < 0 {
		return nil, domain.ErrNotEnoughStock
	}

	foundLine.Quantity += quantity

	return order, nil
}

func (service *ShopService) RemoveItemFromCart(orderID string, productID string) (*domain.Order, error) {
	order, ok := service.orderStore[orderID]
	if !ok {
		return nil, domain.ErrCartNotFound
	}

	_, foundIdx := findLineInOrder(order, productID)
	if foundIdx >= 0 {
		order.Lines = append(order.Lines[:foundIdx], order.Lines[foundIdx+1:]...)
		return order, nil
	}

	return nil, domain.ErrItemNotFoundInCart
}

func (service *ShopService) Checkout(orderID string) (totalAmount float64, err error) {
	order, ok := service.orderStore[orderID]
	if !ok {
		return 0, domain.ErrCartNotFound
	}

	service.invMutex.Lock()
	defer service.invMutex.Unlock()

	for _, line := range order.Lines {
		product, ok := service.inventories[line.ProductID]
		if !ok {
			return 0, domain.ErrSomeProductInCartNotFound
		}

		if (product.Quantity - line.Quantity) < 0 {
			return 0, domain.ErrSomeProductInCartNotEnoughInStock
		}

		product.Quantity -= line.Quantity

		totalAmount += line.UnitPrice * float64(line.Quantity)
	}

	if len(service.promotions) == 0 {
		return totalAmount, nil
	}

	now := time.Now()
	for _, promotion := range service.promotions {
		if !promotion.StartDate.IsZero() && promotion.StartDate.After(now) {
			continue
		}

		if !promotion.EndDate.IsZero() && promotion.EndDate.Before(now) {
			continue
		}

		if promotion.Condition == nil {
			continue
		}

		totalAmount += promotion.Condition.CalculateDiscount(order)
	}

	return totalAmount, nil
}

func findLineInOrder(order *domain.Order, productID string) (foundLine *domain.OrderLine, index int) {
	for idx, line := range order.Lines {
		if line.ProductID == productID {
			foundLine = line
			index = idx
			break
		}
	}

	if foundLine == nil {
		return nil, -1
	}

	return
}
