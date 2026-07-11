package controllers

import (
	"context"
	"errors"
	"net/http"

	"orders-api/dto"
	"orders-api/model"

	"github.com/gin-gonic/gin"
)

type ClientService interface {
	Create(ctx context.Context, req dto.CreateClientRequest) (model.Client, error)
	FindAll(ctx context.Context) ([]model.Client, error)
	FindByID(ctx context.Context, id string) (model.Client, error)
}

type ClientController struct {
	service ClientService
}

func NewClientController(service ClientService) *ClientController {
	return &ClientController{service: service}
}

func (c *ClientController) Create(ctx *gin.Context) {
	var req dto.CreateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		writeClientError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewClientResponse(client))
}

func (c *ClientController) FindAll(ctx *gin.Context) {
	clients, err := c.service.FindAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	responses := make([]dto.ClientResponse, len(clients))
	for i, client := range clients {
		responses[i] = dto.NewClientResponse(client)
	}

	ctx.JSON(http.StatusOK, responses)
}

func (c *ClientController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	client, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		writeClientError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewClientResponse(client))
}

func writeClientError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrClientNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrClientAlreadyExists):
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrClientNameRequired),
		errors.Is(err, model.ErrClientNameTooShort),
		errors.Is(err, model.ErrClientNameTooLong),
		errors.Is(err, model.ErrClientEmailRequired),
		errors.Is(err, model.ErrClientEmailInvalid),
		errors.Is(err, model.ErrClientEmailTooLong),
		errors.Is(err, model.ErrClientPasswordRequired),
		errors.Is(err, model.ErrClientPasswordTooShort),
		errors.Is(err, model.ErrClientPasswordTooLong),
		errors.Is(err, model.ErrClientPasswordHashRequired):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
