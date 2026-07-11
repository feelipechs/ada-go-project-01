package repository

import (
	"context"
	"errors"
	"fmt"

	"orders-api/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, product model.Product) (model.Product, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO products (id, name, price, stock, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, name, price, stock, created_at`,
		product.ID, product.Name, product.Price, product.Stock, product.CreatedAt,
	)

	var created model.Product
	err := row.Scan(&created.ID, &created.Name, &created.Price, &created.Stock, &created.CreatedAt)
	if err != nil {
		return model.Product{}, fmt.Errorf("create product: %w", err)
	}

	return created, nil
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, price, stock, created_at FROM products ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, rows.Err()
}

func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, price, stock, created_at FROM products WHERE id = $1`,
		id,
	)

	var product model.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Product{}, model.ErrProductNotFound
		}
		return model.Product{}, fmt.Errorf("find product by id: %w", err)
	}

	return product, nil
}
