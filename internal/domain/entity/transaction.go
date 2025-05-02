package entity

import (
	"time"
)

// TransactionStatus represents the current status of a payment transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// Transaction represents a payment transaction in the system
type Transaction struct {
	ID               string            `json:"id"`
	OrderID          string            `json:"order_id"`
	Amount           int64             `json:"amount"`
	Status           TransactionStatus `json:"status"`
	PaymentMethod    string            `json:"payment_method"`
	PaymentReference string            `json:"payment_reference"`
	BankName         string            `json:"bank_name"`
	Description      string            `json:"description"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}
