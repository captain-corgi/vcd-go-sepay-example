package entity

import (
	"time"
)

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a customer order that requires payment
type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	Amount      int64       `json:"amount"` // Amount in Vietnamese Dong (VND)
	Status      OrderStatus `json:"status"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// IsValid returns whether the order is in a valid state
func (o *Order) IsValid() bool {
	return o.ID != "" && o.Amount > 0
}
