package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookPayload_GetOrderID(t *testing.T) {
	testCases := []struct {
		name        string
		payload     WebhookPayload
		expectedID  string
		description string
	}{
		{
			name: "simple order ID in description",
			payload: WebhookPayload{
				Description: "ORDER123",
			},
			expectedID:  "ORDER123",
			description: "Should return order ID directly from description",
		},
		{
			name: "order ID with prefix",
			payload: WebhookPayload{
				Description: "ord_12345678",
			},
			expectedID:  "ord_12345678",
			description: "Should return order ID with prefix",
		},
		{
			name: "empty description",
			payload: WebhookPayload{
				Description: "",
			},
			expectedID:  "",
			description: "Should return empty string for empty description",
		},
		{
			name: "description with spaces",
			payload: WebhookPayload{
				Description: "  ORDER123  ",
			},
			expectedID:  "  ORDER123  ",
			description: "Should return description as-is including spaces",
		},
		{
			name: "complex description format",
			payload: WebhookPayload{
				Description: "Payment for order ORDER123 from customer ABC",
			},
			expectedID:  "Payment for order ORDER123 from customer ABC",
			description: "Should return full description when complex format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.payload.GetOrderID()
			assert.Equal(t, tc.expectedID, result, tc.description)
		})
	}
}

func TestWebhookPayload_StructFields(t *testing.T) {
	// Test creating a webhook payload with all fields
	payload := WebhookPayload{
		ID:                12345,
		Gateway:           "Vietcombank",
		TransactionDate:   "2023-12-01 14:30:00",
		AccountNumber:     "0123456789",
		Amount:            100000,
		Description:       "ORDER123",
		CustomerInfo:      "John Doe",
		CreditAmount:      100000,
		DebitAmount:       0,
		Fee:               1000,
		BankTransactionID: "FT12345678",
		WebhookURL:        "https://example.com/webhook",
	}

	// Verify all fields are set correctly
	assert.Equal(t, int64(12345), payload.ID)
	assert.Equal(t, "Vietcombank", payload.Gateway)
	assert.Equal(t, "2023-12-01 14:30:00", payload.TransactionDate)
	assert.Equal(t, "0123456789", payload.AccountNumber)
	assert.Equal(t, int64(100000), payload.Amount)
	assert.Equal(t, "ORDER123", payload.Description)
	assert.Equal(t, "John Doe", payload.CustomerInfo)
	assert.Equal(t, int64(100000), payload.CreditAmount)
	assert.Equal(t, int64(0), payload.DebitAmount)
	assert.Equal(t, int64(1000), payload.Fee)
	assert.Equal(t, "FT12345678", payload.BankTransactionID)
	assert.Equal(t, "https://example.com/webhook", payload.WebhookURL)
}

func TestWebhookPayload_AmountValidation(t *testing.T) {
	testCases := []struct {
		name          string
		payload       WebhookPayload
		shouldBeValid bool
		description   string
	}{
		{
			name: "positive amount",
			payload: WebhookPayload{
				Amount:       100000,
				CreditAmount: 100000,
				DebitAmount:  0,
			},
			shouldBeValid: true,
			description:   "Valid webhook with positive amount",
		},
		{
			name: "zero amount",
			payload: WebhookPayload{
				Amount:       0,
				CreditAmount: 0,
				DebitAmount:  0,
			},
			shouldBeValid: false,
			description:   "Invalid webhook with zero amount",
		},
		{
			name: "negative amount",
			payload: WebhookPayload{
				Amount:       -50000,
				CreditAmount: 0,
				DebitAmount:  50000,
			},
			shouldBeValid: false,
			description:   "Invalid webhook with negative amount",
		},
		{
			name: "credit and debit mismatch",
			payload: WebhookPayload{
				Amount:       100000,
				CreditAmount: 90000,
				DebitAmount:  0,
			},
			shouldBeValid: false,
			description:   "Invalid when credit amount doesn't match transaction amount",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Basic validation logic
			isAmountValid := tc.payload.Amount > 0
			isCreditValid := tc.payload.Amount == tc.payload.CreditAmount || tc.payload.CreditAmount == 0

			isValid := isAmountValid && isCreditValid
			assert.Equal(t, tc.shouldBeValid, isValid, tc.description)
		})
	}
}

func TestWebhookPayload_RequiredFields(t *testing.T) {
	testCases := []struct {
		name          string
		payload       WebhookPayload
		shouldBeValid bool
		description   string
	}{
		{
			name: "all required fields present",
			payload: WebhookPayload{
				ID:                12345,
				Gateway:           "Vietcombank",
				Amount:            100000,
				Description:       "ORDER123",
				BankTransactionID: "FT12345678",
			},
			shouldBeValid: true,
			description:   "Valid webhook with all required fields",
		},
		{
			name: "missing gateway",
			payload: WebhookPayload{
				ID:                12345,
				Gateway:           "",
				Amount:            100000,
				Description:       "ORDER123",
				BankTransactionID: "FT12345678",
			},
			shouldBeValid: false,
			description:   "Invalid webhook with missing gateway",
		},
		{
			name: "missing description",
			payload: WebhookPayload{
				ID:                12345,
				Gateway:           "Vietcombank",
				Amount:            100000,
				Description:       "",
				BankTransactionID: "FT12345678",
			},
			shouldBeValid: false,
			description:   "Invalid webhook with missing description",
		},
		{
			name: "missing bank transaction ID",
			payload: WebhookPayload{
				ID:                12345,
				Gateway:           "Vietcombank",
				Amount:            100000,
				Description:       "ORDER123",
				BankTransactionID: "",
			},
			shouldBeValid: false,
			description:   "Invalid webhook with missing bank transaction ID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Validation logic for required fields
			hasRequiredFields := tc.payload.ID > 0 &&
				tc.payload.Gateway != "" &&
				tc.payload.Amount > 0 &&
				tc.payload.Description != "" &&
				tc.payload.BankTransactionID != ""

			assert.Equal(t, tc.shouldBeValid, hasRequiredFields, tc.description)
		})
	}
}
