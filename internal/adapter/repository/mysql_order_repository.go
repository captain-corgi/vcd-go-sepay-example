package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/repository"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/infrastructure/persistence"
)

// MySQLOrderRepository implements the OrderRepository interface using MySQL
type MySQLOrderRepository struct {
	db *persistence.MySQLDB
}

// NewMySQLOrderRepository creates a new MySQL order repository
func NewMySQLOrderRepository(db *persistence.MySQLDB) repository.OrderRepository {
	return &MySQLOrderRepository{
		db: db,
	}
}

// Create inserts a new order into the database
func (r *MySQLOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	query := `
		INSERT INTO orders (id, customer_id, amount, status, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	if order.UpdatedAt.IsZero() {
		order.UpdatedAt = now
	}

	_, err := r.db.DB.ExecContext(
		ctx,
		query,
		order.ID,
		order.CustomerID,
		order.Amount,
		order.Status,
		order.Description,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// GetByID retrieves an order by its ID
func (r *MySQLOrderRepository) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	query := `
		SELECT id, customer_id, amount, status, description, created_at, updated_at
		FROM orders
		WHERE id = ?
	`

	var order entity.Order
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.CustomerID,
		&order.Amount,
		&order.Status,
		&order.Description,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// Update updates an existing order
func (r *MySQLOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	query := `
		UPDATE orders
		SET customer_id = ?, amount = ?, status = ?, description = ?, updated_at = ?
		WHERE id = ?
	`

	order.UpdatedAt = time.Now()

	_, err := r.db.DB.ExecContext(
		ctx,
		query,
		order.CustomerID,
		order.Amount,
		order.Status,
		order.Description,
		order.UpdatedAt,
		order.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// List retrieves a list of orders with pagination
func (r *MySQLOrderRepository) List(ctx context.Context, limit, offset int) ([]*entity.Order, error) {
	query := `
		SELECT id, customer_id, amount, status, description, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.Amount,
			&order.Status,
			&order.Description,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order rows: %w", err)
	}

	return orders, nil
}
