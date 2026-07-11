package service

import (
	"context"
	"errors"

	"orders-api/dto"
	"orders-api/model"

	"github.com/google/uuid"
)

type OrderClientRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (model.Client, error)
}

type OrderProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (model.Product, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order model.Order, items []model.OrderItem) (model.Order, error)
	FindAll(ctx context.Context, limit, offset int) ([]model.Order, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) (model.Order, error)
	Cancel(ctx context.Context, id uuid.UUID) (model.Order, error)
}

type OrderService struct {
	repo          OrderRepository
	clientRepo    OrderClientRepository
	productRepo   OrderProductRepository
}

func NewOrderService(
	repo OrderRepository,
	clientRepo OrderClientRepository,
	productRepo OrderProductRepository,
) *OrderService {
	return &OrderService{
		repo:        repo,
		clientRepo:  clientRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, req dto.CreateOrderRequest) (model.Order, error) {
	clientUID := req.ClientID

	if _, err := s.clientRepo.FindByID(ctx, clientUID); err != nil {
		if errors.Is(err, model.ErrClientNotFound) {
			return model.Order{}, err
		}
		return model.Order{}, err
	}

	var items []model.OrderItem
	for _, itemReq := range req.Items {
		product, err := s.productRepo.FindByID(ctx, itemReq.ProductID)
		if err != nil {
			if errors.Is(err, model.ErrProductNotFound) {
				return model.Order{}, model.ErrProductNotFound
			}
			return model.Order{}, err
		}

		if product.Stock < itemReq.Quantity {
			return model.Order{}, model.ErrInsufficientStock
		}

		items = append(items, model.NewOrderItem(itemReq.ProductID, itemReq.Quantity, product.Price))
	}

	order, err := model.NewOrder(clientUID, items)
	if err != nil {
		return model.Order{}, err
	}

	created, err := s.repo.Create(ctx, order, items)
	if err != nil {
		return model.Order{}, err
	}

	return created, nil
}

func (s *OrderService) FindAll(ctx context.Context, limit, offset int) ([]model.Order, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	orders, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) FindByID(ctx context.Context, id string) (model.Order, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return model.Order{}, model.ErrOrderNotFound
	}

	order, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *OrderService) Pay(ctx context.Context, id string) (model.Order, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return model.Order{}, model.ErrOrderNotFound
	}

	order, err := s.repo.UpdateStatus(ctx, uid, model.StatusPaid)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *OrderService) Cancel(ctx context.Context, id string) (model.Order, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return model.Order{}, model.ErrOrderNotFound
	}

	order, err := s.repo.Cancel(ctx, uid)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
