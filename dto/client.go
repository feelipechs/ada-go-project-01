package dto

import (
	"time"

	"orders-api/model"

	"github.com/google/uuid"
)

type CreateClientRequest struct {
	Name     string `json:"name"     binding:"required,min=3,max=255"`
	Email    string `json:"email"    binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type ClientResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewClientResponse(client model.Client) ClientResponse {
	return ClientResponse{
		ID:        client.ID,
		Name:      client.Name,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
	}
}
