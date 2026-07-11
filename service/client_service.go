package service

import (
	"context"
	"errors"

	"orders-api/auth"
	"orders-api/dto"
	"orders-api/model"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Create(ctx context.Context, client model.Client) (model.Client, error)
	FindAll(ctx context.Context) ([]model.Client, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.Client, error)
}

type ClientService struct {
	repo ClientRepository
}

func NewClientService(repo ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) Create(ctx context.Context, req dto.CreateClientRequest) (model.Client, error) {
	if err := model.ValidateClientPassword(req.Password); err != nil {
		return model.Client{}, err
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return model.Client{}, err
	}

	client, err := model.NewClient(req.Name, req.Email, passwordHash)
	if err != nil {
		return model.Client{}, err
	}

	created, err := s.repo.Create(ctx, client)
	if err != nil {
		return model.Client{}, err
	}

	return created, nil
}

func (s *ClientService) FindAll(ctx context.Context) ([]model.Client, error) {
	clients, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (s *ClientService) FindByID(ctx context.Context, id string) (model.Client, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return model.Client{}, model.ErrClientNotFound
	}

	client, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if errors.Is(err, model.ErrClientNotFound) {
			return model.Client{}, model.ErrClientNotFound
		}
		return model.Client{}, err
	}

	return client, nil
}
