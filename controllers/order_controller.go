package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"orders-api/dto"
	"orders-api/model"

	"github.com/gin-gonic/gin"
)

type OrderService interface {
	Create(ctx context.Context, req dto.CreateOrderRequest) (model.Order, error)
	FindAll(ctx context.Context, limit, offset int) ([]model.Order, error)
	FindByID(ctx context.Context, id string) (model.Order, error)
	Pay(ctx context.Context, id string) (model.Order, error)
	Cancel(ctx context.Context, id string) (model.Order, error)
}

type OrderController struct {
	service OrderService
}

func NewOrderController(service OrderService) *OrderController {
	return &OrderController{service: service}
}

func (c *OrderController) Create(ctx *gin.Context) {
	var req dto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		writeOrderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewOrderResponse(order))
}

func (c *OrderController) FindAll(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	orders, err := c.service.FindAll(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	responses := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = dto.NewOrderResponse(order)
	}

	ctx.JSON(http.StatusOK, responses)
}

func (c *OrderController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		writeOrderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewOrderResponse(order))
}

func (c *OrderController) Pay(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := c.service.Pay(ctx.Request.Context(), id)
	if err != nil {
		writeOrderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewOrderResponse(order))
}

func (c *OrderController) Cancel(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := c.service.Cancel(ctx.Request.Context(), id)
	if err != nil {
		writeOrderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewOrderResponse(order))
}

func writeOrderError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrOrderNotFound),
		errors.Is(err, model.ErrClientNotFound),
		errors.Is(err, model.ErrProductNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrInsufficientStock),
		errors.Is(err, model.ErrInvalidOrderStatus):
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrClientRequired),
		errors.Is(err, model.ErrOrderItemsRequired),
		errors.Is(err, model.ErrProductRequired),
		errors.Is(err, model.ErrInvalidQuantity),
		errors.Is(err, model.ErrInvalidUnitPrice):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
