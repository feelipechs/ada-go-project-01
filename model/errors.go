package model

import "errors"

var (
	ErrClientNotFound         = errors.New("client not found")
	ErrClientNameRequired     = errors.New("name is required")
	ErrClientNameTooShort     = errors.New("name must have at least 3 characters")
	ErrClientNameTooLong      = errors.New("name must have at most 255 characters")
	ErrClientEmailRequired    = errors.New("email is required")
	ErrClientEmailInvalid     = errors.New("email must be valid")
	ErrClientEmailTooLong     = errors.New("email must have at most 255 characters")
	ErrClientAlreadyExists    = errors.New("email already exists")
	ErrClientPasswordRequired = errors.New("password is required")
	ErrClientPasswordTooShort = errors.New("password must have at least 8 characters")
	ErrClientPasswordTooLong  = errors.New("password must have at most 72 bytes")

	ErrClientPasswordHashRequired = errors.New("password hash is required")

	ErrProductNotFound      = errors.New("product not found")
	ErrProductNameRequired  = errors.New("product name is required")
	ErrProductNameTooLong   = errors.New("product name must have at most 255 characters")
	ErrProductPriceInvalid  = errors.New("price must be greater than zero")
	ErrProductStockInvalid  = errors.New("stock cannot be negative")

	ErrOrderNotFound        = errors.New("order not found")
	ErrClientRequired       = errors.New("client is required")
	ErrOrderItemsRequired   = errors.New("order must have at least one item")
	ErrProductRequired      = errors.New("product is required")
	ErrInvalidQuantity      = errors.New("quantity must be greater than zero")
	ErrInvalidUnitPrice     = errors.New("unit price must be greater than zero")
	ErrInvalidOrderStatus   = errors.New("order cannot change from current status")
	ErrInsufficientStock    = errors.New("insufficient stock")
)
