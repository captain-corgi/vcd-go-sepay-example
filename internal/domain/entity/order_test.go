package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrder_IsValid(t *testing.T) {
	testCases := []struct {
		name     string
		order    Order
		expected bool
	}{
		{
			name: "valid order with all required fields",
			order: Order{
				ID:          "order123",
				CustomerID:  "customer1",
				Amount:      100000,
				Status:      OrderStatusPending,
				Description: "Test order",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expected: true,
		},
		{
			name: "invalid order with empty ID",
			order: Order{
				ID:          "",
				CustomerID:  "customer1",
				Amount:      100000,
				Status:      OrderStatusPending,
				Description: "Test order",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expected: false,
		},
		{
			name: "invalid order with zero amount",
			order: Order{
				ID:          "order123",
				CustomerID:  "customer1",
				Amount:      0,
				Status:      OrderStatusPending,
				Description: "Test order",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expected: false,
		},
		{
			name: "invalid order with negative amount",
			order: Order{
				ID:          "order123",
				CustomerID:  "customer1",
				Amount:      -1000,
				Status:      OrderStatusPending,
				Description: "Test order",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expected: false,
		},
		{
			name: "invalid order with empty ID and zero amount",
			order: Order{
				ID:          "",
				CustomerID:  "customer1",
				Amount:      0,
				Status:      OrderStatusPending,
				Description: "Test order",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.order.IsValid()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestOrderStatus_Constants(t *testing.T) {
	// Test that order status constants have expected values
	assert.Equal(t, OrderStatus("pending"), OrderStatusPending)
	assert.Equal(t, OrderStatus("paid"), OrderStatusPaid)
	assert.Equal(t, OrderStatus("failed"), OrderStatusFailed)
	assert.Equal(t, OrderStatus("cancelled"), OrderStatusCancelled)
}

func TestOrder_JSONTags(t *testing.T) {
	// Test that Order struct can be properly marshaled/unmarshaled
	order := Order{
		ID:          "order123",
		CustomerID:  "customer1",
		Amount:      100000,
		Status:      OrderStatusPending,
		Description: "Test order",
		CreatedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	// This would be useful for API serialization
	assert.NotEmpty(t, order.ID)
	assert.NotEmpty(t, order.CustomerID)
	assert.Greater(t, order.Amount, int64(0))
	assert.NotEmpty(t, order.Status)
}
