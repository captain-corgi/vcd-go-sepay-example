# Sepay Payment Gateway Integration Instructions

This document provides comprehensive guidelines for implementing the Sepay payment gateway integration in Go projects following Clean Architecture, Domain-Driven Design (DDD), and Table-Driven Testing principles.

## What is Sepay?

Sepay is a Vietnamese fintech platform that provides payment gateway services through bank transfers. It connects directly with various Vietnamese banks via Open Banking APIs to offer:

- Real-time transaction monitoring
- Automatic payment confirmation through webhooks
- Support for VietQR code generation
- Lower transaction fees compared to traditional payment gateways

## Integration Overview

The Sepay integration workflow follows these steps:

1. Customer places an order in your application
2. Your application generates a VietQR payment code
3. Customer scans QR code with their banking app
4. Customer completes bank transfer
5. Bank notifies Sepay of the transaction
6. Sepay sends a webhook to your application
7. Your application verifies the payment and updates order status

## Clean Architecture Implementation

For a maintainable, testable, and framework-independent implementation, adhere to Clean Architecture principles:

### Project Structure

```
sepay-integration/
├── cmd/
│   └── api/
│       └── main.go                # Application entry point
├── internal/
│   ├── domain/                    # Core domain layer (entities & business rules)
│   │   ├── entity/
│   │   │   ├── order.go
│   │   │   ├── payment.go
│   │   │   └── transaction.go
│   │   ├── repository/            # Repository interfaces
│   │   │   ├── order_repository.go
│   │   │   └── transaction_repository.go
│   │   ├── service/               # Domain service interfaces
│   │   │   └── payment_service.go
│   │   └── vo/                    # Value objects
│   │       ├── money.go
│   │       └── qrcode.go
│   ├── usecase/                   # Application use cases
│   │   ├── create_payment.go
│   │   ├── process_webhook.go
│   │   └── update_transaction.go
│   ├── adapter/                   # Interface adapters
│   │   ├── api/                   # REST API
│   │   │   ├── handler/
│   │   │   │   ├── order_handler.go
│   │   │   │   └── webhook_handler.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth_middleware.go
│   │   │   │   └── logging_middleware.go
│   │   │   └── response/
│   │   │       └── response.go
│   │   ├── repository/            # Repository implementations
│   │   │   ├── mysql/
│   │   │   │   ├── mysql_order_repository.go
│   │   │   │   └── mysql_transaction_repository.go
│   │   │   └── sqlite/
│   │   │       └── sqlite_repository.go
│   │   └── service/               # Service implementations
│   │       └── sepay/
│   │           └── payment_service_impl.go
│   └── infrastructure/            # External frameworks & drivers
│       ├── config/
│       │   └── config.go
│       ├── persistence/
│       │   └── mysql.go
│       ├── qrcode/
│       │   └── vietqr_generator.go
│       └── http/
│           └── server.go
├── pkg/                           # Shared packages
│   ├── logger/
│   │   └── logger.go
│   ├── validator/
│   │   └── validator.go
│   └── errors/
│       └── errors.go
└── test/                          # Integration & E2E tests
    ├── integration/
    │   └── webhook_test.go
    └── mocks/
        └── mock_repositories.go
```

### Layer Responsibilities

1. **Domain Layer** - Contains business entities, value objects, domain services, and repository interfaces
2. **Use Case Layer** - Contains application logic that orchestrates domain objects
3. **Adapter Layer** - Contains implementations of interfaces defined in inner layers
4. **Infrastructure Layer** - Contains code that interacts with external systems

## Domain-Driven Design Implementation

Apply Domain-Driven Design principles to model a rich, expressive domain model:

### Ubiquitous Language

Define a common vocabulary shared by developers and domain experts:

- **Order** - A customer purchase that requires payment
- **Transaction** - A record of payment for an order
- **Payment** - The process of transferring funds
- **QR Code** - A visual representation of payment instructions

### Entity Implementation

Entities have identity and lifecycle:

```go
// internal/domain/entity/transaction.go
package entity

import (
	"errors"
	"time"
	
	"your-module/internal/domain/vo"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// Transaction is a rich domain entity with identity and business rules
type Transaction struct {
	ID               string
	OrderID          string
	Amount           vo.Money
	Status           TransactionStatus
	PaymentMethod    string
	PaymentReference string
	BankName         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// MarkAsCompleted validates and changes transaction status to completed
func (t *Transaction) MarkAsCompleted(bankReference string) error {
	if t.Status == TransactionStatusCompleted {
		return errors.New("transaction already completed")
	}
	
	if bankReference == "" {
		return errors.New("bank reference cannot be empty")
	}
	
	t.Status = TransactionStatusCompleted
	t.PaymentReference = bankReference
	t.UpdatedAt = time.Time{}
	
	return nil
}

// VerifyAmount checks if the received amount matches the expected amount
func (t *Transaction) VerifyAmount(receivedAmount vo.Money) bool {
	return t.Amount.Equals(receivedAmount)
}
```

### Value Objects

Value objects represent concepts with no identity:

```go
// internal/domain/vo/money.go
package vo

import "errors"

// Money is a value object representing an amount with currency
type Money struct {
	Amount   int64
	Currency string
}

// NewMoney creates a Money value object with validation
func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	return Money{Amount: amount, Currency: currency}, nil
}

// Equals compares two Money value objects
func (m Money) Equals(other Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}

// Add returns a new Money value object with the sum
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot add different currencies")
	}
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}
```

### Domain Services

Domain services implement business logic that doesn't fit naturally in entities:

```go
// internal/domain/service/payment_service.go
package service

import (
	"context"
	
	"your-module/internal/domain/entity"
	"your-module/internal/domain/vo"
)

// PaymentService defines operations for payment processing
type PaymentService interface {
	GeneratePaymentQR(ctx context.Context, orderID string, amount vo.Money, description string) (vo.QRCode, error)
	VerifyPaymentWebhook(ctx context.Context, payload []byte, signature string) (bool, error)
	GetTransactionStatus(ctx context.Context, transactionID string) (entity.TransactionStatus, error)
}
```

### Repositories

Repositories provide data access abstraction:

```go
// internal/domain/repository/transaction_repository.go
package repository

import (
	"context"
	
	"your-module/internal/domain/entity"
)

// TransactionRepository defines data access operations for transactions
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	FindPendingTransactionsOlderThan(ctx context.Context, age time.Duration) ([]*entity.Transaction, error)
}
```

## Use Case Implementation

Use cases orchestrate domain entities and services:

```go
// internal/usecase/process_webhook.go
package usecase

import (
	"context"
	"time"
	
	"your-module/internal/domain/entity"
	"your-module/internal/domain/repository"
	"your-module/internal/domain/service"
	"your-module/internal/domain/vo"
)

// WebhookData represents the incoming webhook payload
type WebhookData struct {
	ID                int64  `json:"id"`
	Gateway           string `json:"gateway"`
	TransactionDate   string `json:"transactionDate"`
	AccountNumber     string `json:"accountNumber"`
	Amount            int64  `json:"amount"`
	Description       string `json:"description"` // Contains OrderID
	CustomerInfo      string `json:"customerInfo"`
	CreditAmount      int64  `json:"creditAmount"`
	DebitAmount       int64  `json:"debitAmount"`
	Fee               int64  `json:"fee"`
	BankTransactionID string `json:"bankTransactionId"`
	WebhookURL        string `json:"webhookUrl"`
}

// ProcessWebhookUseCase handles incoming payment webhooks
type ProcessWebhookUseCase struct {
	transactionRepo repository.TransactionRepository
	orderRepo       repository.OrderRepository
	paymentService  service.PaymentService
	logger          Logger
}

// Execute processes the webhook data
func (uc *ProcessWebhookUseCase) Execute(ctx context.Context, data WebhookData) error {
	// Extract order ID from description (domain logic)
	orderID := data.Description
	
	// Convert amount to Money value object
	amount, err := vo.NewMoney(data.Amount, "VND")
	if err != nil {
		uc.logger.Error("Invalid amount in webhook", err)
		return err
	}
	
	// Get related transaction
	transaction, err := uc.transactionRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		uc.logger.Error("Failed to find transaction", err, "orderID", orderID)
		return err
	}
	
	// Verify amount matches expected transaction amount
	if !transaction.VerifyAmount(amount) {
		uc.logger.Error("Amount mismatch", nil, "expected", transaction.Amount, "received", amount)
		return errors.New("amount mismatch")
	}
	
	// Update transaction status using domain entity method
	if err := transaction.MarkAsCompleted(data.BankTransactionID); err != nil {
		uc.logger.Error("Failed to mark transaction as completed", err)
		return err
	}
	
	// Save updated transaction
	if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
		uc.logger.Error("Failed to update transaction", err)
		return err
	}
	
	// Update order
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		uc.logger.Error("Failed to find order", err, "orderID", orderID)
		return err
	}
	
	// Update order status
	order.Status = "paid"
	order.UpdatedAt = time.Now()
	
	return uc.orderRepo.Update(ctx, order)
}
```

## Table-Driven Testing

Always use table-driven testing to ensure comprehensive test coverage with minimal code duplication:

### Domain Entity Tests

```go
// internal/domain/entity/transaction_test.go
package entity_test

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	"your-module/internal/domain/entity"
	"your-module/internal/domain/vo"
)

func TestTransaction_MarkAsCompleted(t *testing.T) {
	// Table-driven test cases
	tests := []struct {
		name           string
		transaction    entity.Transaction
		bankReference  string
		expectedError  bool
		expectedStatus entity.TransactionStatus
	}{
		{
			name: "valid transaction",
			transaction: entity.Transaction{
				ID:     "tx1",
				Status: entity.TransactionStatusPending,
			},
			bankReference:  "BR12345",
			expectedError:  false,
			expectedStatus: entity.TransactionStatusCompleted,
		},
		{
			name: "already completed transaction",
			transaction: entity.Transaction{
				ID:     "tx2",
				Status: entity.TransactionStatusCompleted,
			},
			bankReference:  "BR12345",
			expectedError:  true,
			expectedStatus: entity.TransactionStatusCompleted,
		},
		{
			name: "empty bank reference",
			transaction: entity.Transaction{
				ID:     "tx3",
				Status: entity.TransactionStatusPending,
			},
			bankReference:  "",
			expectedError:  true,
			expectedStatus: entity.TransactionStatusPending,
		},
	}
	
	// Execute all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the transaction to test with
			tx := tt.transaction
			
			// Execute the method being tested
			err := tx.MarkAsCompleted(tt.bankReference)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			// Verify transaction status
			assert.Equal(t, tt.expectedStatus, tx.Status)
			
			// If completed, verify the bank reference was set
			if !tt.expectedError {
				assert.Equal(t, tt.bankReference, tx.PaymentReference)
			}
		})
	}
}
```

### Use Case Tests

```go
// internal/usecase/process_webhook_test.go
package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"your-module/internal/domain/entity"
	"your-module/internal/domain/vo"
	"your-module/test/mocks"
	"your-module/internal/usecase"
)

func TestProcessWebhookUseCase(t *testing.T) {
	// Table-driven test cases
	tests := []struct {
		name          string
		webhookData   usecase.WebhookData
		setupMocks    func(*mocks.TransactionRepository, *mocks.OrderRepository)
		expectedError bool
	}{
		{
			name: "successful processing",
			webhookData: usecase.WebhookData{
				ID:                12345,
				Gateway:           "Vietcombank",
				TransactionDate:   "2025-05-02 14:30:00",
				AccountNumber:     "0123456789",
				Amount:            150000,
				Description:       "ORDER123",
				CustomerInfo:      "Test Customer",
				CreditAmount:      150000,
				DebitAmount:       0,
				Fee:               0,
				BankTransactionID: "FT12345678",
				WebhookURL:        "https://api.yourdomain.com/sepay/webhook",
			},
			setupMocks: func(tr *mocks.TransactionRepository, or *mocks.OrderRepository) {
				// Create Money value object for testing
				amount, _ := vo.NewMoney(150000, "VND")
				
				// Mock transaction repository behavior
				transaction := &entity.Transaction{
					ID:      "txn-123",
					OrderID: "ORDER123",
					Amount:  amount,
					Status:  entity.TransactionStatusPending,
				}
				
				tr.On("GetByOrderID", mock.Anything, "ORDER123").Return(transaction, nil)
				tr.On("Update", mock.Anything, mock.MatchedBy(func(t *entity.Transaction) bool {
					return t.Status == entity.TransactionStatusCompleted &&
						t.PaymentReference == "FT12345678"
				})).Return(nil)
				
				// Mock order repository behavior
				order := &entity.Order{
					ID:     "ORDER123",
					Status: "pending",
				}
				
				or.On("GetByID", mock.Anything, "ORDER123").Return(order, nil)
				or.On("Update", mock.Anything, mock.MatchedBy(func(o *entity.Order) bool {
					return o.Status == "paid"
				})).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "transaction not found",
			webhookData: usecase.WebhookData{
				Description: "ORDER999",
				Amount:      150000,
			},
			setupMocks: func(tr *mocks.TransactionRepository, or *mocks.OrderRepository) {
				tr.On("GetByOrderID", mock.Anything, "ORDER999").Return(nil, errors.New("transaction not found"))
			},
			expectedError: true,
		},
		{
			name: "amount mismatch",
			webhookData: usecase.WebhookData{
				Description: "ORDER123",
				Amount:      200000, // Different from expected amount
			},
			setupMocks: func(tr *mocks.TransactionRepository, or *mocks.OrderRepository) {
				expectedAmount, _ := vo.NewMoney(150000, "VND")
				
				transaction := &entity.Transaction{
					ID:      "txn-123",
					OrderID: "ORDER123",
					Amount:  expectedAmount,
					Status:  entity.TransactionStatusPending,
				}
				
				tr.On("GetByOrderID", mock.Anything, "ORDER123").Return(transaction, nil)
			},
			expectedError: true,
		},
	}
	
	// Execute all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			transactionRepo := new(mocks.TransactionRepository)
			orderRepo := new(mocks.OrderRepository)
			logger := new(mocks.Logger)
			
			// Setup mocks
			tt.setupMocks(transactionRepo, orderRepo)
			
			// Create use case with mocked dependencies
			uc := usecase.NewProcessWebhookUseCase(transactionRepo, orderRepo, logger)
			
			// Execute use case
			err := uc.Execute(context.Background(), tt.webhookData)
			
			// Assert result
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			// Verify expectations on mocks
			transactionRepo.AssertExpectations(t)
			orderRepo.AssertExpectations(t)
		})
	}
}
```

### Handler Tests

```go
// internal/adapter/api/handler/webhook_handler_test.go
package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"your-module/internal/adapter/api/handler"
	"your-module/internal/usecase"
	"your-module/test/mocks"
)

func TestWebhookHandler_Handle(t *testing.T) {
	// Table-driven test cases
	tests := []struct {
		name           string
		requestMethod  string
		requestBody    interface{}
		setupMocks     func(*mocks.ProcessWebhookUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "valid webhook",
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"id":                12345,
				"gateway":           "Vietcombank",
				"transactionDate":   "2025-05-02 14:30:00",
				"accountNumber":     "0123456789",
				"amount":            150000,
				"description":       "ORDER123",
				"customerInfo":      "Test Customer",
				"creditAmount":      150000,
				"debitAmount":       0,
				"fee":               0,
				"bankTransactionId": "FT12345678",
				"webhookUrl":        "https://api.yourdomain.com/sepay/webhook",
			},
			setupMocks: func(uc *mocks.ProcessWebhookUseCase) {
				uc.On("Execute", mock.Anything, mock.MatchedBy(func(data usecase.WebhookData) bool {
					return data.Description == "ORDER123" && data.Amount == 150000
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}`,
		},
		{
			name:          "invalid method",
			requestMethod: http.MethodGet,
			requestBody:   map[string]interface{}{},
			setupMocks:    func(*mocks.ProcessWebhookUseCase) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:          "invalid JSON",
			requestMethod: http.MethodPost,
			requestBody:   "invalid json",
			setupMocks:    func(*mocks.ProcessWebhookUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:          "use case error",
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"description": "ORDER123",
				"amount":      150000,
			},
			setupMocks: func(uc *mocks.ProcessWebhookUseCase) {
				uc.On("Execute", mock.Anything, mock.Anything).Return(errors.New("use case error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}
	
	// Execute all test cases
	for _, tt := = tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock use case
			useCase := new(mocks.ProcessWebhookUseCase)
			
			// Setup mocks
			tt.setupMocks(useCase)
			
			// Create handler
			handler := handler.NewWebhookHandler(useCase)
			
			// Create test HTTP request
			var reqBody []byte
			var err error
			
			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}
			
			req := httptest.NewRequest(tt.requestMethod, "/sepay/webhook", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Handle the request
			handler.Handle(rr, req)
			
			// Assert response
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// If expecting specific body, check it
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}
			
			// Verify mock expectations
			useCase.AssertExpectations(t)
		})
	}
}
```

## Configuration and Dependencies

Use dependency injection for better testability:

```go
// cmd/api/main.go
package main

import (
	"context"
	"log"
	
	"your-module/internal/adapter/api/handler"
	"your-module/internal/adapter/repository/mysql"
	"your-module/internal/adapter/service/sepay"
	"your-module/internal/infrastructure/config"
	"your-module/internal/infrastructure/http"
	"your-module/internal/infrastructure/persistence"
	"your-module/internal/usecase"
	"your-module/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Set up logger
	appLogger := logger.NewLogger(cfg.LogLevel)
	
	// Set up database connection
	db, err := persistence.NewMySQLConnection(cfg.Database)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()
	
	// Initialize repositories
	orderRepo := mysql.NewOrderRepository(db)
	transactionRepo := mysql.NewTransactionRepository(db)
	
	// Initialize services
	paymentService := sepay.NewPaymentService(
		cfg.Sepay.APIKey,
		cfg.Sepay.APISecret,
		cfg.Sepay.BankID,
		cfg.Sepay.AccountNumber,
		cfg.Sepay.AccountName,
	)
	
	// Initialize use cases
	createPaymentUC := usecase.NewCreatePaymentUseCase(
		orderRepo,
		transactionRepo,
		paymentService,
		appLogger,
	)
	
	processWebhookUC := usecase.NewProcessWebhookUseCase(
		transactionRepo,
		orderRepo,
		paymentService,
		appLogger,
	)
	
	// Initialize handlers
	orderHandler := handler.NewOrderHandler(createPaymentUC, appLogger)
	webhookHandler := handler.NewWebhookHandler(processWebhookUC, appLogger)
	
	// Initialize HTTP server
	server := http.NewServer(cfg.Server, appLogger)
	
	// Register routes
	server.RegisterRoute("POST", "/orders/create", orderHandler.Create)
	server.RegisterRoute("POST", "/sepay/webhook", webhookHandler.Handle)
	
	// Start server
	if err := server.Start(context.Background()); err != nil {
		appLogger.Fatal("Server failed", err)
	}
}
```

## Testing Best Practices

1. **Test Domain Logic Thoroughly** - Most tests should focus on domain entities and use cases
2. **Mock External Dependencies** - Use interfaces and mocks for external services
3. **Use Table-Driven Tests** - Always use table-driven tests for comprehensive test coverage
4. **Test Edge Cases** - Include tests for error conditions and edge cases
5. **Separate Unit and Integration Tests** - Keep unit tests fast and integration tests comprehensive

## Security Best Practices

1. **Webhook Authentication** - Implement signature verification using HMAC-SHA256
2. **Data Validation** - Validate all incoming data against domain constraints
3. **HTTPS** - Use secure connections for all API endpoints
4. **Idempotent Processing** - Store webhook IDs to prevent duplicate processing
5. **Transaction Reconciliation** - Implement a background job to check for missed transactions
6. **IP Filtering** - Restrict webhook access to Sepay's IP addresses

## Error Handling

Implement consistent error handling throughout the application:

```go
// pkg/errors/errors.go
package errors

import "fmt"

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Err     error
}

// Error returns the error message
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewPaymentError creates a payment-specific domain error
func NewPaymentError(message string, err error) *DomainError {
	return &DomainError{
		Code:    "PAYMENT_ERROR",
		Message: message,
		Err:     err,
	}
}

// NewWebhookError creates a webhook-specific domain error
func NewWebhookError(message string, err error) *DomainError {
	return &DomainError{
		Code:    "WEBHOOK_ERROR",
		Message: message,
		Err:     err,
	}
}
```

## Reference Resources

- Sepay Website: [https://sepay.vn](https://sepay.vn)
- Sepay Dashboard: [https://my.sepay.vn](https://my.sepay.vn)
- QR Generator: [https://qr.sepay.vn](https://qr.sepay.vn)
- VietQR Official Website: [https://vietqr.net](https://vietqr.net)
- API Endpoint: [https://partner-api.sepay.vn/merchant/v1](https://partner-api.sepay.vn/merchant/v1)
- Clean Architecture: [https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- Domain-Driven Design: [https://martinfowler.com/tags/domain%20driven%20design.html](https://martinfowler.com/tags/domain%20driven%20design.html)
- Table-Driven Testing in Go: [https://github.com/golang/go/wiki/TableDrivenTests](https://github.com/golang/go/wiki/TableDrivenTests)