package domain

type OrderLine struct {
	ID        string
	ProductID string
	Quantity  int
	UnitPrice float64
}
