package model

import (
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	productNameMaxLength = 255
)

type Product struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price float64, stock int) (Product, error) {
	p := Product{
		ID:        uuid.New(),
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: time.Now(),
	}

	if err := validateProductName(p.Name); err != nil {
		return Product{}, err
	}
	if err := validateProductPrice(p.Price); err != nil {
		return Product{}, err
	}
	if err := validateProductStock(p.Stock); err != nil {
		return Product{}, err
	}

	return p, nil
}

func validateProductName(name string) error {
	if name == "" {
		return ErrProductNameRequired
	}
	if utf8.RuneCountInString(name) > productNameMaxLength {
		return ErrProductNameTooLong
	}
	return nil
}

func validateProductPrice(price float64) error {
	if price <= 0 {
		return ErrProductPriceInvalid
	}
	return nil
}

func validateProductStock(stock int) error {
	if stock < 0 {
		return ErrProductStockInvalid
	}
	return nil
}
