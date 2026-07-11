package model

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}

func NewOrderItem(productID uuid.UUID, quantity int, unitPrice float64) OrderItem {
	return OrderItem{
		ID:        uuid.New(),
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
}

func (i OrderItem) Validate() error {
	if i.ProductID == uuid.Nil {
		return ErrProductRequired
	}
	if i.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	if i.UnitPrice <= 0 {
		return ErrInvalidUnitPrice
	}
	return nil
}

func (i OrderItem) Total() float64 {
	return float64(i.Quantity) * i.UnitPrice
}
