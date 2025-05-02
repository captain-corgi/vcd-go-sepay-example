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
