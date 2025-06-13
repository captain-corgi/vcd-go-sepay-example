package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/repository/mocks"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProcessWebhookUseCase_Execute(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		payload       *entity.WebhookPayload
		setupMocks    func(*gomock.Controller, *mocks.MockOrderRepository, *mocks.MockTransactionRepository)
		expectedError bool
	}{
		{
			name: "successful webhook processing with existing transaction",
			payload: &entity.WebhookPayload{
				ID:                12345,
				Gateway:           "Vietcombank",
				TransactionDate:   "2025-05-01 14:30:00",
				AccountNumber:     "0123456789",
				Amount:            100000,
				Description:       "ORDER123",
				CustomerInfo:      "John Doe",
				CreditAmount:      100000,
				DebitAmount:       0,
				Fee:               0,
				BankTransactionID: "FT12345678",
				WebhookURL:        "https://example.com/webhook",
			},
			setupMocks: func(ctrl *gomock.Controller, or *mocks.MockOrderRepository, tr *mocks.MockTransactionRepository) {
				// Setup GetByID to return a valid order
				or.EXPECT().GetByID(gomock.Any(), "ORDER123").Return(&entity.Order{
					ID:          "ORDER123",
					CustomerID:  "customer1",
					Amount:      100000,
					Status:      entity.OrderStatusPending,
					Description: "Test order",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)

				// Setup GetByOrderID to return a valid transaction
				tr.EXPECT().GetByOrderID(gomock.Any(), "ORDER123").Return(&entity.Transaction{
					ID:               "txn_123",
					OrderID:          "ORDER123",
					Amount:           100000,
					Status:           entity.TransactionStatusPending,
					PaymentMethod:    "bank_transfer",
					PaymentReference: "",
					BankName:         "",
					Description:      "ORDER123",
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				}, nil)

				// Setup Update to succeed
				tr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "order not found",
			payload: &entity.WebhookPayload{
				Description: "NONEXISTENT",
				Amount:      100000,
			},
			setupMocks: func(ctrl *gomock.Controller, or *mocks.MockOrderRepository, tr *mocks.MockTransactionRepository) {
				or.EXPECT().GetByID(gomock.Any(), "NONEXISTENT").Return(nil, errors.New("order not found"))
			},
			expectedError: true,
		},
		{
			name: "amount mismatch",
			payload: &entity.WebhookPayload{
				Description: "ORDER123",
				Amount:      200000, // Incorrect amount
			},
			setupMocks: func(ctrl *gomock.Controller, or *mocks.MockOrderRepository, tr *mocks.MockTransactionRepository) {
				or.EXPECT().GetByID(gomock.Any(), "ORDER123").Return(&entity.Order{
					ID:     "ORDER123",
					Amount: 100000, // Expected amount is different
				}, nil)
			},
			expectedError: true,
		},
		{
			name: "transaction update fails",
			payload: &entity.WebhookPayload{
				ID:                12345,
				Gateway:           "Vietcombank",
				Amount:            100000,
				Description:       "ORDER123",
				BankTransactionID: "FT12345678",
			},
			setupMocks: func(ctrl *gomock.Controller, or *mocks.MockOrderRepository, tr *mocks.MockTransactionRepository) {
				// Return valid order
				or.EXPECT().GetByID(gomock.Any(), "ORDER123").Return(&entity.Order{
					ID:     "ORDER123",
					Amount: 100000,
				}, nil)

				// Return valid transaction
				tr.EXPECT().GetByOrderID(gomock.Any(), "ORDER123").Return(&entity.Transaction{
					ID:      "txn_123",
					OrderID: "ORDER123",
					Amount:  100000,
					Status:  entity.TransactionStatusPending,
				}, nil)

				// But update fails
				tr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
			},
			expectedError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create controller and mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := mocks.NewMockOrderRepository(ctrl)
			transactionRepo := mocks.NewMockTransactionRepository(ctrl)

			// Set up mock expectations
			tc.setupMocks(ctrl, orderRepo, transactionRepo)

			// Create use case
			uc := usecase.NewProcessWebhookUseCase(orderRepo, transactionRepo)

			// Execute use case
			err := uc.Execute(context.Background(), tc.payload)

			// Check results
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// No need to explicitly verify expectations - gomock's controller will do this automatically
		})
	}
}

func TestProcessWebhookUseCase_Execute_WithoutExistingTransaction(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)

	payload := &entity.WebhookPayload{
		ID:                12345,
		Gateway:           "Vietcombank",
		TransactionDate:   "2025-05-01 14:30:00",
		AccountNumber:     "0123456789",
		Amount:            100000,
		Description:       "ORDER123",
		CustomerInfo:      "John Doe",
		CreditAmount:      100000,
		DebitAmount:       0,
		Fee:               0,
		BankTransactionID: "FT12345678",
		WebhookURL:        "https://example.com/webhook",
	}

	// Setup GetByID to return a valid order
	orderRepo.EXPECT().GetByID(gomock.Any(), "ORDER123").Return(&entity.Order{
		ID:          "ORDER123",
		CustomerID:  "customer1",
		Amount:      100000,
		Status:      entity.OrderStatusPending,
		Description: "Test order",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	// Setup GetByOrderID to return no transaction (not found)
	transactionRepo.EXPECT().GetByOrderID(gomock.Any(), "ORDER123").Return(nil, errors.New("transaction not found"))

	// Note: Based on the implementation, even when creating a new transaction,
	// it will call Update (not Create) because the ID is always set
	// This seems like a bug in the implementation, but we test the current behavior
	transactionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, transaction *entity.Transaction) error {
			// Verify new transaction structure
			assert.Equal(t, "txn_12345", transaction.ID) // Should use Sepay's transaction ID
			assert.Equal(t, "ORDER123", transaction.OrderID)
			assert.Equal(t, int64(100000), transaction.Amount)
			assert.Equal(t, entity.TransactionStatusCompleted, transaction.Status)
			assert.Equal(t, "bank_transfer", transaction.PaymentMethod)
			assert.Equal(t, "FT12345678", transaction.PaymentReference)
			assert.Equal(t, "Vietcombank", transaction.BankName)
			assert.Equal(t, "ORDER123", transaction.Description)
			return nil
		},
	)

	// Setup order update
	orderRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, order *entity.Order) error {
			assert.Equal(t, entity.OrderStatusPaid, order.Status)
			return nil
		},
	)

	// Create use case
	uc := usecase.NewProcessWebhookUseCase(orderRepo, transactionRepo)

	// Execute
	err := uc.Execute(context.Background(), payload)

	// Assertions
	assert.NoError(t, err)
}

func TestProcessWebhookUseCase_Execute_OrderUpdateFails(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)

	payload := &entity.WebhookPayload{
		ID:                12345,
		Gateway:           "Vietcombank",
		Amount:            100000,
		Description:       "ORDER123",
		BankTransactionID: "FT12345678",
	}

	// Return valid order
	orderRepo.EXPECT().GetByID(gomock.Any(), "ORDER123").Return(&entity.Order{
		ID:     "ORDER123",
		Amount: 100000,
		Status: entity.OrderStatusPending,
	}, nil)

	// Return valid transaction
	transactionRepo.EXPECT().GetByOrderID(gomock.Any(), "ORDER123").Return(&entity.Transaction{
		ID:      "txn_123",
		OrderID: "ORDER123",
		Amount:  100000,
		Status:  entity.TransactionStatusPending,
	}, nil)

	// Transaction update succeeds
	transactionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	// But order update fails
	orderRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("order update failed"))

	// Create use case
	uc := usecase.NewProcessWebhookUseCase(orderRepo, transactionRepo)

	// Execute
	err := uc.Execute(context.Background(), payload)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update order status")
}

func TestProcessWebhookUseCase_Execute_EmptyOrderID(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)

	payload := &entity.WebhookPayload{
		ID:                12345,
		Gateway:           "Vietcombank",
		Amount:            100000,
		Description:       "", // Empty description results in empty order ID
		BankTransactionID: "FT12345678",
	}

	// Create use case
	uc := usecase.NewProcessWebhookUseCase(orderRepo, transactionRepo)

	// Execute
	err := uc.Execute(context.Background(), payload)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid webhook: order ID not found in description")
}

func TestProcessWebhookUseCase_Execute_WebhookDataValidation(t *testing.T) {
	testCases := []struct {
		name          string
		payload       *entity.WebhookPayload
		expectedError string
		description   string
	}{
		{
			name: "negative amount",
			payload: &entity.WebhookPayload{
				ID:                12345,
				Gateway:           "Vietcombank",
				Amount:            -100000,
				Description:       "ORDER123",
				BankTransactionID: "FT12345678",
			},
			expectedError: "payment amount mismatch",
			description:   "Negative amounts should cause validation error",
		},
		{
			name: "zero amount",
			payload: &entity.WebhookPayload{
				ID:                12345,
				Gateway:           "Vietcombank",
				Amount:            0,
				Description:       "ORDER123",
				BankTransactionID: "FT12345678",
			},
			expectedError: "payment amount mismatch",
			description:   "Zero amounts should cause validation error",
		},
		{
			name: "missing bank transaction ID",
			payload: &entity.WebhookPayload{
				ID:                12345,
				Gateway:           "Vietcombank",
				Amount:            100000,
				Description:       "ORDER123",
				BankTransactionID: "",
			},
			expectedError: "", // This will succeed with empty bank transaction ID
			description:   "Missing bank transaction ID",
		},
		{
			name: "missing gateway",
			payload: &entity.WebhookPayload{
				ID:                12345,
				Gateway:           "",
				Amount:            100000,
				Description:       "ORDER123",
				BankTransactionID: "FT12345678",
			},
			expectedError: "", // This will succeed with empty gateway
			description:   "Missing gateway information",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := mocks.NewMockOrderRepository(ctrl)
			transactionRepo := mocks.NewMockTransactionRepository(ctrl)

			if tc.payload.Description != "" {
				// Only set up mocks if we expect to get past the order ID validation
				orderRepo.EXPECT().GetByID(gomock.Any(), tc.payload.Description).Return(&entity.Order{
					ID:     tc.payload.Description,
					Amount: 100000,
					Status: entity.OrderStatusPending,
				}, nil).AnyTimes()

				// If we expect no error, set up the full flow
				if tc.expectedError == "" {
					// Setup transaction repository calls for successful flow
					transactionRepo.EXPECT().GetByOrderID(gomock.Any(), tc.payload.Description).Return(&entity.Transaction{
						ID:      "txn_existing",
						OrderID: tc.payload.Description,
						Amount:  100000,
						Status:  entity.TransactionStatusPending,
					}, nil).AnyTimes()

					transactionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
					orderRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				}
			}

			// Create use case
			uc := usecase.NewProcessWebhookUseCase(orderRepo, transactionRepo)

			// Execute
			err := uc.Execute(context.Background(), tc.payload)

			// Assertions
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessWebhookUseCase_Execute_ConcurrentWebhooks(t *testing.T) {
	// Test concurrent processing of webhooks for the same order
	// This tests that the use case handles race conditions appropriately

	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)

	payload := &entity.WebhookPayload{
		ID:                12345,
		Gateway:           "Vietcombank",
		Amount:            100000,
		Description:       "ORDER123",
		BankTransactionID: "FT12345678",
	}

	// Setup mocks to be called multiple times (simulating concurrent calls)
	orderRepo.EXPECT().GetByID(gomock.Any(), "ORDER123").Return(&entity.Order{
		ID:     "ORDER123",
		Amount: 100000,
		Status: entity.OrderStatusPending,
	}, nil).Times(2)

	transactionRepo.EXPECT().GetByOrderID(gomock.Any(), "ORDER123").Return(&entity.Transaction{
		ID:      "txn_123",
		OrderID: "ORDER123",
		Amount:  100000,
		Status:  entity.TransactionStatusPending,
	}, nil).Times(2)

	transactionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	orderRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(2)

	// Create use case
	uc := usecase.NewProcessWebhookUseCase(orderRepo, transactionRepo)

	// Execute concurrently
	done := make(chan error, 2)

	go func() {
		done <- uc.Execute(context.Background(), payload)
	}()

	go func() {
		done <- uc.Execute(context.Background(), payload)
	}()

	// Wait for both to complete
	err1 := <-done
	err2 := <-done

	// Both should succeed (in a real scenario, you might want to handle duplicates)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestProcessWebhookUseCase_Execute_DifferentBankGateways(t *testing.T) {
	testCases := []struct {
		name    string
		gateway string
		valid   bool
	}{
		{
			name:    "Vietcombank",
			gateway: "Vietcombank",
			valid:   true,
		},
		{
			name:    "BIDV",
			gateway: "BIDV",
			valid:   true,
		},
		{
			name:    "Techcombank",
			gateway: "Techcombank",
			valid:   true,
		},
		{
			name:    "Unknown bank",
			gateway: "Unknown Bank",
			valid:   true, // Assuming the system accepts any bank
		},
		{
			name:    "Empty gateway",
			gateway: "",
			valid:   true, // This might be business rule dependent
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := mocks.NewMockOrderRepository(ctrl)
			transactionRepo := mocks.NewMockTransactionRepository(ctrl)

			payload := &entity.WebhookPayload{
				ID:                12345,
				Gateway:           tc.gateway,
				Amount:            100000,
				Description:       "ORDER123",
				BankTransactionID: "FT12345678",
			}

			if tc.valid {
				// Setup successful flow
				orderRepo.EXPECT().GetByID(gomock.Any(), "ORDER123").Return(&entity.Order{
					ID:     "ORDER123",
					Amount: 100000,
					Status: entity.OrderStatusPending,
				}, nil)

				transactionRepo.EXPECT().GetByOrderID(gomock.Any(), "ORDER123").Return(&entity.Transaction{
					ID:      "txn_123",
					OrderID: "ORDER123",
					Amount:  100000,
					Status:  entity.TransactionStatusPending,
				}, nil)

				transactionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, transaction *entity.Transaction) error {
						// Verify bank name is set correctly
						assert.Equal(t, tc.gateway, transaction.BankName)
						return nil
					},
				)

				orderRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			}

			// Create use case
			uc := usecase.NewProcessWebhookUseCase(orderRepo, transactionRepo)

			// Execute
			err := uc.Execute(context.Background(), payload)

			// Assertions
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
