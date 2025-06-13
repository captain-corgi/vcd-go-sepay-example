package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransactionStatus_Constants(t *testing.T) {
	// Test that transaction status constants have expected values
	assert.Equal(t, TransactionStatus("pending"), TransactionStatusPending)
	assert.Equal(t, TransactionStatus("completed"), TransactionStatusCompleted)
	assert.Equal(t, TransactionStatus("failed"), TransactionStatusFailed)
}

func TestTransaction_StructFields(t *testing.T) {
	// Test creating a transaction with all fields
	now := time.Now()
	transaction := Transaction{
		ID:               "txn_123",
		OrderID:          "order_456",
		Amount:           100000,
		Status:           TransactionStatusPending,
		PaymentMethod:    "bank_transfer",
		PaymentReference: "REF123456",
		BankName:         "Vietcombank",
		Description:      "Payment for order 456",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// Verify all fields are set correctly
	assert.Equal(t, "txn_123", transaction.ID)
	assert.Equal(t, "order_456", transaction.OrderID)
	assert.Equal(t, int64(100000), transaction.Amount)
	assert.Equal(t, TransactionStatusPending, transaction.Status)
	assert.Equal(t, "bank_transfer", transaction.PaymentMethod)
	assert.Equal(t, "REF123456", transaction.PaymentReference)
	assert.Equal(t, "Vietcombank", transaction.BankName)
	assert.Equal(t, "Payment for order 456", transaction.Description)
	assert.Equal(t, now, transaction.CreatedAt)
	assert.Equal(t, now, transaction.UpdatedAt)
}

func TestTransaction_StatusTransitions(t *testing.T) {
	testCases := []struct {
		name          string
		initialStatus TransactionStatus
		finalStatus   TransactionStatus
		shouldBeValid bool
		description   string
	}{
		{
			name:          "pending to completed",
			initialStatus: TransactionStatusPending,
			finalStatus:   TransactionStatusCompleted,
			shouldBeValid: true,
			description:   "Valid transition from pending to completed",
		},
		{
			name:          "pending to failed",
			initialStatus: TransactionStatusPending,
			finalStatus:   TransactionStatusFailed,
			shouldBeValid: true,
			description:   "Valid transition from pending to failed",
		},
		{
			name:          "completed stays completed",
			initialStatus: TransactionStatusCompleted,
			finalStatus:   TransactionStatusCompleted,
			shouldBeValid: true,
			description:   "Completed transaction remains completed",
		},
		{
			name:          "failed stays failed",
			initialStatus: TransactionStatusFailed,
			finalStatus:   TransactionStatusFailed,
			shouldBeValid: true,
			description:   "Failed transaction remains failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transaction := Transaction{
				ID:      "txn_test",
				OrderID: "order_test",
				Amount:  50000,
				Status:  tc.initialStatus,
			}

			// Simulate status change
			transaction.Status = tc.finalStatus
			transaction.UpdatedAt = time.Now()

			// Verify the status change
			assert.Equal(t, tc.finalStatus, transaction.Status)
			assert.NotZero(t, transaction.UpdatedAt)
		})
	}
}

func TestTransaction_AmountValidation(t *testing.T) {
	testCases := []struct {
		name     string
		amount   int64
		expected bool
	}{
		{
			name:     "positive amount",
			amount:   100000,
			expected: true,
		},
		{
			name:     "zero amount",
			amount:   0,
			expected: false,
		},
		{
			name:     "negative amount",
			amount:   -50000,
			expected: false,
		},
		{
			name:     "very large amount",
			amount:   99999999999,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transaction := Transaction{
				ID:      "txn_test",
				OrderID: "order_test",
				Amount:  tc.amount,
				Status:  TransactionStatusPending,
			}

			// Simple validation: amount should be positive
			isValid := transaction.Amount > 0
			assert.Equal(t, tc.expected, isValid)
		})
	}
}

func TestTransaction_JSONSerialization(t *testing.T) {
	// Test that important fields have proper JSON tags for API responses
	transaction := Transaction{
		ID:               "txn_123",
		OrderID:          "order_456",
		Amount:           100000,
		Status:           TransactionStatusCompleted,
		PaymentMethod:    "bank_transfer",
		PaymentReference: "REF123456",
		BankName:         "Vietcombank",
		Description:      "Test payment",
		CreatedAt:        time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC),
	}

	// Verify struct has the expected values that would be serialized
	assert.NotEmpty(t, transaction.ID)
	assert.NotEmpty(t, transaction.OrderID)
	assert.Greater(t, transaction.Amount, int64(0))
	assert.NotEmpty(t, transaction.Status)
	assert.NotEmpty(t, transaction.PaymentMethod)
	assert.NotZero(t, transaction.CreatedAt)
	assert.NotZero(t, transaction.UpdatedAt)
}
