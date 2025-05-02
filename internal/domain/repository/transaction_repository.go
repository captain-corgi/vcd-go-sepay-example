package repository

import (
	"context"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
)

//go:generate mockgen --source=transaction_repository.go --destination=mocks/transaction_repository_mock.go --package=mocks

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	List(ctx context.Context, limit, offset int) ([]*entity.Transaction, error)
}
