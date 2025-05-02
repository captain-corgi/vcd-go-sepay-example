package repository

import (
	"context"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
)

//go:generate mockgen --source=order_repository.go --destination=mocks/order_repository_mock.go --package=mocks

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id string) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	List(ctx context.Context, limit, offset int) ([]*entity.Order, error)
}
