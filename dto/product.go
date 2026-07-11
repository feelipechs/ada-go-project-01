package dto

import (
	"time"

	"orders-api/model"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name  string  `json:"name"  binding:"required,max=255"`
	Price float64 `json:"price" binding:"required,gt=0"`
	Stock int     `json:"stock" binding:"gte=0"`
}

type ProductResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProductResponse(product model.Product) ProductResponse {
	return ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
	}
}
