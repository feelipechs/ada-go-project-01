package service

import (
	"context"

	"orders-api/dto"
	"orders-api/model"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product model.Product) (model.Product, error)
	FindAll(ctx context.Context) ([]model.Product, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, req dto.CreateProductRequest) (model.Product, error) {
	product, err := model.NewProduct(req.Name, req.Price, req.Stock)
	if err != nil {
		return model.Product{}, err
	}

	created, err := s.repo.Create(ctx, product)
	if err != nil {
		return model.Product{}, err
	}

	return created, nil
}

func (s *ProductService) FindAll(ctx context.Context) ([]model.Product, error) {
	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) FindByID(ctx context.Context, id string) (model.Product, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return model.Product{}, model.ErrProductNotFound
	}

	product, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}
