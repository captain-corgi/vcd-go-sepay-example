package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/usecase"
)

// Mock repositories
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) List(ctx context.Context, limit, offset int) ([]*entity.Order, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*entity.Order), args.Error(1)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) List(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

func TestProcessWebhookUseCase_Execute(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		payload       *entity.WebhookPayload
		setupMocks    func(*MockOrderRepository, *MockTransactionRepository)
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
			setupMocks: func(or *MockOrderRepository, tr *MockTransactionRepository) {
				// Setup GetByID to return a valid order
				or.On("GetByID", mock.Anything, "ORDER123").Return(&entity.Order{
					ID:          "ORDER123",
					CustomerID:  "customer1",
					Amount:      100000,
					Status:      entity.OrderStatusPending,
					Description: "Test order",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)

				// Setup GetByOrderID to return a valid transaction
				tr.On("GetByOrderID", mock.Anything, "ORDER123").Return(&entity.Transaction{
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
				tr.On("Update", mock.Anything, mock.AnythingOfType("*entity.Transaction")).Return(nil)
				or.On("Update", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "order not found",
			payload: &entity.WebhookPayload{
				Description: "NONEXISTENT",
				Amount:      100000,
			},
			setupMocks: func(or *MockOrderRepository, tr *MockTransactionRepository) {
				or.On("GetByID", mock.Anything, "NONEXISTENT").Return(nil, errors.New("order not found"))
			},
			expectedError: true,
		},
		{
			name: "amount mismatch",
			payload: &entity.WebhookPayload{
				Description: "ORDER123",
				Amount:      200000, // Incorrect amount
			},
			setupMocks: func(or *MockOrderRepository, tr *MockTransactionRepository) {
				or.On("GetByID", mock.Anything, "ORDER123").Return(&entity.Order{
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
			setupMocks: func(or *MockOrderRepository, tr *MockTransactionRepository) {
				// Return valid order
				or.On("GetByID", mock.Anything, "ORDER123").Return(&entity.Order{
					ID:     "ORDER123",
					Amount: 100000,
				}, nil)

				// Return valid transaction
				tr.On("GetByOrderID", mock.Anything, "ORDER123").Return(&entity.Transaction{
					ID:      "txn_123",
					OrderID: "ORDER123",
					Amount:  100000,
					Status:  entity.TransactionStatusPending,
				}, nil)

				// But update fails
				tr.On("Update", mock.Anything, mock.AnythingOfType("*entity.Transaction")).Return(errors.New("database error"))
			},
			expectedError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			orderRepo := new(MockOrderRepository)
			transactionRepo := new(MockTransactionRepository)

			// Set up mock expectations
			tc.setupMocks(orderRepo, transactionRepo)

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

			// Verify that all expected mock calls were made
			orderRepo.AssertExpectations(t)
			transactionRepo.AssertExpectations(t)
		})
	}
}
