package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending  OrderStatus = "PENDING"
	StatusPaid     OrderStatus = "PAID"
	StatusCanceled OrderStatus = "CANCELED"
)

type Order struct {
	ID        uuid.UUID   `json:"id"`
	ClientID  uuid.UUID   `json:"client_id"`
	Status    OrderStatus `json:"status"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	Items     []OrderItem `json:"items"`
}

func NewOrder(clientID uuid.UUID, items []OrderItem) (Order, error) {
	if clientID == uuid.Nil {
		return Order{}, ErrClientRequired
	}
	if len(items) == 0 {
		return Order{}, ErrOrderItemsRequired
	}

	total := 0.0
	for _, item := range items {
		if err := item.Validate(); err != nil {
			return Order{}, err
		}
		total += item.Total()
	}

	return Order{
		ID:        uuid.New(),
		ClientID:  clientID,
		Status:    StatusPending,
		Total:     total,
		CreatedAt: time.Now(),
		Items:     items,
	}, nil
}

func (o *Order) Pay() error {
	if o.Status != StatusPending {
		return ErrInvalidOrderStatus
	}
	o.Status = StatusPaid
	return nil
}

func (o *Order) Cancel() error {
	if o.Status != StatusPending {
		return ErrInvalidOrderStatus
	}
	o.Status = StatusCanceled
	return nil
}

func (o Order) IsPending() bool {
	return o.Status == StatusPending
}

func (o Order) IsPaid() bool {
	return o.Status == StatusPaid
}

func (o Order) IsCanceled() bool {
	return o.Status == StatusCanceled
}
