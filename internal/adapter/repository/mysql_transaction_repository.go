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

// MySQLTransactionRepository implements the TransactionRepository interface using MySQL
type MySQLTransactionRepository struct {
	db *persistence.MySQLDB
}

// NewMySQLTransactionRepository creates a new MySQL transaction repository
func NewMySQLTransactionRepository(db *persistence.MySQLDB) repository.TransactionRepository {
	return &MySQLTransactionRepository{
		db: db,
	}
}

// Create inserts a new transaction into the database
func (r *MySQLTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, order_id, amount, status, payment_method, 
			payment_reference, bank_name, description, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = now
	}
	if transaction.UpdatedAt.IsZero() {
		transaction.UpdatedAt = now
	}

	_, err := r.db.DB.ExecContext(
		ctx,
		query,
		transaction.ID,
		transaction.OrderID,
		transaction.Amount,
		transaction.Status,
		transaction.PaymentMethod,
		transaction.PaymentReference,
		transaction.BankName,
		transaction.Description,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a transaction by its ID
func (r *MySQLTransactionRepository) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	query := `
		SELECT id, order_id, amount, status, payment_method, 
		       payment_reference, bank_name, description, created_at, updated_at
		FROM transactions
		WHERE id = ?
	`

	var transaction entity.Transaction
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.OrderID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.PaymentMethod,
		&transaction.PaymentReference,
		&transaction.BankName,
		&transaction.Description,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// GetByOrderID retrieves a transaction by its order ID
func (r *MySQLTransactionRepository) GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error) {
	query := `
		SELECT id, order_id, amount, status, payment_method, 
		       payment_reference, bank_name, description, created_at, updated_at
		FROM transactions
		WHERE order_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var transaction entity.Transaction
	err := r.db.DB.QueryRowContext(ctx, query, orderID).Scan(
		&transaction.ID,
		&transaction.OrderID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.PaymentMethod,
		&transaction.PaymentReference,
		&transaction.BankName,
		&transaction.Description,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found for order: %s", orderID)
		}
		return nil, fmt.Errorf("failed to get transaction by order ID: %w", err)
	}

	return &transaction, nil
}

// Update updates an existing transaction
func (r *MySQLTransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	query := `
		UPDATE transactions
		SET order_id = ?, amount = ?, status = ?, payment_method = ?,
		    payment_reference = ?, bank_name = ?, description = ?, updated_at = ?
		WHERE id = ?
	`

	transaction.UpdatedAt = time.Now()

	_, err := r.db.DB.ExecContext(
		ctx,
		query,
		transaction.OrderID,
		transaction.Amount,
		transaction.Status,
		transaction.PaymentMethod,
		transaction.PaymentReference,
		transaction.BankName,
		transaction.Description,
		transaction.UpdatedAt,
		transaction.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

// List retrieves a list of transactions with pagination
func (r *MySQLTransactionRepository) List(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	query := `
		SELECT id, order_id, amount, status, payment_method, 
		       payment_reference, bank_name, description, created_at, updated_at
		FROM transactions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]*entity.Transaction, 0)
	for rows.Next() {
		var transaction entity.Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.OrderID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.PaymentMethod,
			&transaction.PaymentReference,
			&transaction.BankName,
			&transaction.Description,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}
