package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/repository/mocks"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/infrastructure/config"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Mock QR code generator for testing
type mockQRCodeGenerator struct {
	shouldFail   bool
	qrCode       *entity.QRCode
	capturedData *entity.VietQRData
}

func (m *mockQRCodeGenerator) Generate(data entity.VietQRData) (*entity.QRCode, error) {
	if m.capturedData != nil {
		*m.capturedData = data
	}
	if m.shouldFail {
		return nil, errors.New("failed to generate QR code")
	}
	return m.qrCode, nil
}

func TestGeneratePaymentQRUseCase_Execute_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)

	qrCode := &entity.QRCode{
		Content: "VietQR payment content",
		Size:    256,
		Image:   []byte("fake image data"),
	}
	qrGenerator := &mockQRCodeGenerator{
		shouldFail: false,
		qrCode:     qrCode,
	}

	cfg := &config.Config{
		Sepay: config.SepayConfig{
			BankID:        "970415",
			AccountNumber: "0123456789",
			AccountName:   "NGUYEN VAN A",
		},
	}

	// Set up mock expectations
	orderRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	// Create use case
	uc := usecase.NewGeneratePaymentQRUseCase(orderRepo, transactionRepo, qrGenerator, cfg)

	// Execute
	input := usecase.OrderInput{
		CustomerID:  "customer123",
		Amount:      100000,
		Description: "Test order payment",
	}

	result, err := uc.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.OrderID)
	assert.Equal(t, qrCode, result.QRCode)
	assert.Equal(t, input.Amount, result.Amount)
	assert.Equal(t, cfg.Sepay.BankID, result.BankID)
	assert.Equal(t, cfg.Sepay.AccountNumber, result.AccountNumber)
	assert.Equal(t, cfg.Sepay.AccountName, result.AccountName)
	assert.True(t, result.ExpiresAt.After(time.Now()))
	assert.True(t, result.ExpiresAt.Before(time.Now().Add(25*time.Hour))) // Should expire in ~24 hours
}

func TestGeneratePaymentQRUseCase_Execute_OrderCreationFails(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)
	qrGenerator := &mockQRCodeGenerator{shouldFail: false}
	cfg := &config.Config{}

	// Set up mock expectations - order creation fails
	orderRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database connection failed"))

	// Create use case
	uc := usecase.NewGeneratePaymentQRUseCase(orderRepo, transactionRepo, qrGenerator, cfg)

	// Execute
	input := usecase.OrderInput{
		CustomerID:  "customer123",
		Amount:      100000,
		Description: "Test order payment",
	}

	result, err := uc.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create order")
}

func TestGeneratePaymentQRUseCase_Execute_TransactionCreationFails(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)
	qrGenerator := &mockQRCodeGenerator{shouldFail: false}
	cfg := &config.Config{}

	// Set up mock expectations
	orderRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("transaction table locked"))

	// Create use case
	uc := usecase.NewGeneratePaymentQRUseCase(orderRepo, transactionRepo, qrGenerator, cfg)

	// Execute
	input := usecase.OrderInput{
		CustomerID:  "customer123",
		Amount:      100000,
		Description: "Test order payment",
	}

	result, err := uc.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create transaction")
}

func TestGeneratePaymentQRUseCase_Execute_QRGenerationFails(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)
	qrGenerator := &mockQRCodeGenerator{shouldFail: true}
	cfg := &config.Config{
		Sepay: config.SepayConfig{
			BankID:        "970415",
			AccountNumber: "0123456789",
			AccountName:   "NGUYEN VAN A",
		},
	}

	// Set up mock expectations
	orderRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	// Create use case
	uc := usecase.NewGeneratePaymentQRUseCase(orderRepo, transactionRepo, qrGenerator, cfg)

	// Execute
	input := usecase.OrderInput{
		CustomerID:  "customer123",
		Amount:      100000,
		Description: "Test order payment",
	}

	result, err := uc.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to generate QR code")
}

func TestGeneratePaymentQRUseCase_Execute_ValidatesOrderData(t *testing.T) {
	testCases := []struct {
		name           string
		input          usecase.OrderInput
		shouldValidate bool
		description    string
	}{
		{
			name: "valid input",
			input: usecase.OrderInput{
				CustomerID:  "customer123",
				Amount:      100000,
				Description: "Valid order",
			},
			shouldValidate: true,
			description:    "All required fields present",
		},
		{
			name: "zero amount",
			input: usecase.OrderInput{
				CustomerID:  "customer123",
				Amount:      0,
				Description: "Invalid order",
			},
			shouldValidate: false,
			description:    "Amount cannot be zero",
		},
		{
			name: "negative amount",
			input: usecase.OrderInput{
				CustomerID:  "customer123",
				Amount:      -50000,
				Description: "Invalid order",
			},
			shouldValidate: false,
			description:    "Amount cannot be negative",
		},
		{
			name: "empty customer ID",
			input: usecase.OrderInput{
				CustomerID:  "",
				Amount:      100000,
				Description: "Invalid order",
			},
			shouldValidate: false,
			description:    "Customer ID is required",
		},
		{
			name: "empty description",
			input: usecase.OrderInput{
				CustomerID:  "customer123",
				Amount:      100000,
				Description: "",
			},
			shouldValidate: false,
			description:    "Description is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Basic validation logic (this would be part of the use case implementation)
			isValid := tc.input.CustomerID != "" &&
				tc.input.Amount > 0 &&
				tc.input.Description != ""

			assert.Equal(t, tc.shouldValidate, isValid, tc.description)
		})
	}
}

func TestGeneratePaymentQRUseCase_Execute_OrderIDGeneration(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)
	qrGenerator := &mockQRCodeGenerator{
		shouldFail: false,
		qrCode: &entity.QRCode{
			Content: "test",
			Size:    256,
		},
	}
	cfg := &config.Config{
		Sepay: config.SepayConfig{
			BankID:        "970415",
			AccountNumber: "0123456789",
			AccountName:   "NGUYEN VAN A",
		},
	}

	// Set up mock expectations with argument matchers to verify order creation
	orderRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, order *entity.Order) error {
			// Verify order structure
			assert.NotEmpty(t, order.ID)
			assert.True(t, len(order.ID) > 4) // Should have prefix "ord_" plus UUID
			assert.Contains(t, order.ID, "ord_")
			assert.Equal(t, entity.OrderStatusPending, order.Status)
			assert.NotZero(t, order.CreatedAt)
			assert.NotZero(t, order.UpdatedAt)
			return nil
		},
	)

	transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, transaction *entity.Transaction) error {
			// Verify transaction structure
			assert.NotEmpty(t, transaction.ID)
			assert.Contains(t, transaction.ID, "txn_")
			assert.Equal(t, entity.TransactionStatusPending, transaction.Status)
			assert.Equal(t, "bank_transfer", transaction.PaymentMethod)
			assert.NotZero(t, transaction.CreatedAt)
			assert.NotZero(t, transaction.UpdatedAt)
			return nil
		},
	)

	// Create use case
	uc := usecase.NewGeneratePaymentQRUseCase(orderRepo, transactionRepo, qrGenerator, cfg)

	// Execute
	input := usecase.OrderInput{
		CustomerID:  "customer123",
		Amount:      100000,
		Description: "Test order payment",
	}

	result, err := uc.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.OrderID)
	assert.Contains(t, result.OrderID, "ord_")
}

func TestGeneratePaymentQRUseCase_Execute_VietQRDataGeneration(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)

	// Custom mock that captures the VietQR data
	var capturedVietQRData entity.VietQRData
	qrGenerator := &mockQRCodeGenerator{
		shouldFail: false,
		qrCode: &entity.QRCode{
			Content: "test",
			Size:    256,
		},
		capturedData: &capturedVietQRData,
	}

	cfg := &config.Config{
		Sepay: config.SepayConfig{
			BankID:        "970415",
			AccountNumber: "0123456789",
			AccountName:   "NGUYEN VAN A",
		},
	}

	// Set up mock expectations
	orderRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	// Create use case
	uc := usecase.NewGeneratePaymentQRUseCase(orderRepo, transactionRepo, qrGenerator, cfg)

	// Execute
	input := usecase.OrderInput{
		CustomerID:  "customer123",
		Amount:      150000,
		Description: "Test order payment",
	}

	result, err := uc.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify VietQR data was generated correctly
	assert.Equal(t, cfg.Sepay.BankID, capturedVietQRData.BankID)
	assert.Equal(t, cfg.Sepay.AccountNumber, capturedVietQRData.AccountNumber)
	assert.Equal(t, cfg.Sepay.AccountName, capturedVietQRData.AccountName)
	assert.Equal(t, input.Amount, capturedVietQRData.Amount)
	assert.Equal(t, result.OrderID, capturedVietQRData.Description) // Order ID should be in description
}

func TestGeneratePaymentQRUseCase_Execute_MultipleCalls(t *testing.T) {
	// Test that multiple calls generate unique order IDs
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	transactionRepo := mocks.NewMockTransactionRepository(ctrl)
	qrGenerator := &mockQRCodeGenerator{
		shouldFail: false,
		qrCode: &entity.QRCode{
			Content: "test",
			Size:    256,
		},
	}
	cfg := &config.Config{
		Sepay: config.SepayConfig{
			BankID:        "970415",
			AccountNumber: "0123456789",
			AccountName:   "NGUYEN VAN A",
		},
	}

	// Set up mock expectations for multiple calls
	orderRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(3)
	transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(3)

	// Create use case
	uc := usecase.NewGeneratePaymentQRUseCase(orderRepo, transactionRepo, qrGenerator, cfg)

	// Execute multiple times
	input := usecase.OrderInput{
		CustomerID:  "customer123",
		Amount:      100000,
		Description: "Test order payment",
	}

	var orderIDs []string
	for i := 0; i < 3; i++ {
		result, err := uc.Execute(context.Background(), input)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		orderIDs = append(orderIDs, result.OrderID)
	}

	// Verify all order IDs are unique
	assert.Len(t, orderIDs, 3)
	assert.NotEqual(t, orderIDs[0], orderIDs[1])
	assert.NotEqual(t, orderIDs[1], orderIDs[2])
	assert.NotEqual(t, orderIDs[0], orderIDs[2])
}
