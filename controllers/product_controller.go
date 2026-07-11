package controllers

import (
	"context"
	"errors"
	"net/http"

	"orders-api/dto"
	"orders-api/model"

	"github.com/gin-gonic/gin"
)

type ProductService interface {
	Create(ctx context.Context, req dto.CreateProductRequest) (model.Product, error)
	FindAll(ctx context.Context) ([]model.Product, error)
	FindByID(ctx context.Context, id string) (model.Product, error)
}

type ProductController struct {
	service ProductService
}

func NewProductController(service ProductService) *ProductController {
	return &ProductController{service: service}
}

func (c *ProductController) Create(ctx *gin.Context) {
	var req dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		writeProductError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewProductResponse(product))
}

func (c *ProductController) FindAll(ctx *gin.Context) {
	products, err := c.service.FindAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = dto.NewProductResponse(product)
	}

	ctx.JSON(http.StatusOK, responses)
}

func (c *ProductController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		writeProductError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewProductResponse(product))
}

func writeProductError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrProductNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrProductNameRequired),
		errors.Is(err, model.ErrProductNameTooLong),
		errors.Is(err, model.ErrProductPriceInvalid),
		errors.Is(err, model.ErrProductStockInvalid):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
