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

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Create(ctx context.Context, order model.Order, items []model.OrderItem) (model.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.Order{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO orders (id, client_id, status, total, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		order.ID, order.ClientID, order.Status, order.Total, order.CreatedAt,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("insert order: %w", err)
	}

	for _, item := range items {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (id, order_id, product_id, quantity, unit_price)
			 VALUES ($1, $2, $3, $4, $5)`,
			item.ID, order.ID, item.ProductID, item.Quantity, item.UnitPrice,
		)
		if err != nil {
			return model.Order{}, fmt.Errorf("insert order item: %w", err)
		}

		tag, err := tx.Exec(ctx,
			`UPDATE products SET stock = stock - $1
			 WHERE id = $2 AND stock >= $1`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return model.Order{}, fmt.Errorf("update stock: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return model.Order{}, model.ErrInsufficientStock
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Order{}, fmt.Errorf("commit transaction: %w", err)
	}

	order.Items = items
	return order, nil
}

func (r *OrderRepository) FindAll(ctx context.Context, limit, offset int) ([]model.Order, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, client_id, status, total, created_at
		 FROM orders
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("find all orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.ClientID, &o.Status, &o.Total, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range orders {
		items, err := r.findItemsByOrderID(ctx, orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("find items for order %s: %w", orders[i].ID, err)
		}
		orders[i].Items = items
	}

	return orders, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Order, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, client_id, status, total, created_at
		 FROM orders WHERE id = $1`,
		id,
	)

	var order model.Order
	err := row.Scan(&order.ID, &order.ClientID, &order.Status, &order.Total, &order.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, fmt.Errorf("find order by id: %w", err)
	}

	items, err := r.findItemsByOrderID(ctx, id)
	if err != nil {
		return model.Order{}, err
	}
	order.Items = items

	return order, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) (model.Order, error) {
	tag, err := r.pool.Exec(ctx,
		`UPDATE orders SET status = $1
		 WHERE id = $2 AND status = 'PENDING'`,
		status, id,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("update order status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		order, findErr := r.FindByID(ctx, id)
		if findErr != nil {
			return model.Order{}, findErr
		}
		if !order.IsPending() {
			return model.Order{}, model.ErrInvalidOrderStatus
		}
		return model.Order{}, model.ErrOrderNotFound
	}

	return r.FindByID(ctx, id)
}

func (r *OrderRepository) Cancel(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, err := r.FindByID(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	if !order.IsPending() {
		return model.Order{}, model.ErrInvalidOrderStatus
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.Order{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, item := range order.Items {
		_, err := tx.Exec(ctx,
			`UPDATE products SET stock = stock + $1 WHERE id = $2`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return model.Order{}, fmt.Errorf("restore stock: %w", err)
		}
	}

	tag, err := tx.Exec(ctx,
		`UPDATE orders SET status = 'CANCELED' WHERE id = $1 AND status = 'PENDING'`,
		id,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("update order status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.Order{}, model.ErrInvalidOrderStatus
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Order{}, fmt.Errorf("commit transaction: %w", err)
	}

	order.Status = model.StatusCanceled
	return order, nil
}

func (r *OrderRepository) findItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.OrderItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, order_id, product_id, quantity, unit_price
		 FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("find items: %w", err)
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
