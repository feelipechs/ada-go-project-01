package repository

import (
	"context"
	"errors"
	"fmt"

	"orders-api/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{pool: pool}
}

func (r *ClientRepository) Create(ctx context.Context, client model.Client) (model.Client, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO clients (id, name, email, password_hash, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, name, email, created_at`,
		client.ID, client.Name, client.Email, client.PasswordHash, client.CreatedAt,
	)

	var created model.Client
	err := row.Scan(&created.ID, &created.Name, &created.Email, &created.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Client{}, model.ErrClientAlreadyExists
		}
		return model.Client{}, fmt.Errorf("create client: %w", err)
	}

	return created, nil
}

func (r *ClientRepository) FindAll(ctx context.Context) ([]model.Client, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, email, created_at FROM clients ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("find all clients: %w", err)
	}
	defer rows.Close()

	var clients []model.Client
	for rows.Next() {
		var c model.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan client: %w", err)
		}
		clients = append(clients, c)
	}

	return clients, rows.Err()
}

func (r *ClientRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Client, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, email, created_at FROM clients WHERE id = $1`,
		id,
	)

	var client model.Client
	err := row.Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Client{}, model.ErrClientNotFound
		}
		return model.Client{}, fmt.Errorf("find client by id: %w", err)
	}

	return client, nil
}
