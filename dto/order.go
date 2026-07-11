package dto

import (
	"time"

	"orders-api/model"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	ClientID uuid.UUID            `json:"client_id" binding:"required"`
	Items    []CreateOrderItemDTO `json:"items"     binding:"required,min=1,dive"`
}

type CreateOrderItemDTO struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity"   binding:"required,gt=0"`
}

type OrderResponse struct {
	ID        uuid.UUID          `json:"id"`
	ClientID  uuid.UUID          `json:"client_id"`
	Status    string             `json:"status"`
	Total     float64            `json:"total"`
	CreatedAt time.Time          `json:"created_at"`
	Items     []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}

func NewOrderResponse(order model.Order) OrderResponse {
	items := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return OrderResponse{
		ID:        order.ID,
		ClientID:  order.ClientID,
		Status:    string(order.Status),
		Total:     order.Total,
		CreatedAt: order.CreatedAt,
		Items:     items,
	}
}
