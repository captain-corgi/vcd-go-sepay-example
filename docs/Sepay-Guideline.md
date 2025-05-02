# Sepay Payment Gateway Integration with Golang

## Table of Contents
- [Introduction](#introduction)
- [How Sepay Works](#how-sepay-works)
- [Setting Up a Sepay Account](#setting-up-a-sepay-account)
- [Implementation Steps](#implementation-steps)
  - [1. Creating Payment QR Codes](#1-creating-payment-qr-codes)
  - [2. Setting Up Webhook Endpoint](#2-setting-up-webhook-endpoint)
  - [3. Processing Webhook Data](#3-processing-webhook-data)
  - [4. Updating Transaction Status](#4-updating-transaction-status)
- [Clean Architecture Implementation](#clean-architecture-implementation)
  - [Project Structure](#project-structure)
  - [Domain Layer](#domain-layer)
  - [Use Case Layer](#use-case-layer)
  - [Interface Adapters Layer](#interface-adapters-layer)
  - [Frameworks & Drivers Layer](#frameworks--drivers-layer)
- [Domain-Driven Design Approach](#domain-driven-design-approach)
  - [Core Domain Entities](#core-domain-entities)
  - [Value Objects](#value-objects)
  - [Domain Services](#domain-services)
  - [Repositories](#repositories)
- [Code Examples](#code-examples)
  - [Basic Project Structure](#basic-project-structure)
  - [Creating a QR Code Generator](#creating-a-qr-code-generator)
  - [Setting Up a Webhook Server](#setting-up-a-webhook-server)
  - [Processing Webhook Payload](#processing-webhook-payload)
  - [Database Integration](#database-integration)
- [Testing and Debugging](#testing-and-debugging)
  - [Table-Driven Testing](#table-driven-testing)
  - [Mocking External Dependencies](#mocking-external-dependencies)
  - [Integration Testing](#integration-testing)
- [Security Considerations](#security-considerations)
- [Advanced Configuration](#advanced-configuration)
- [Reference](#reference)

## Introduction

Sepay is a Vietnamese fintech platform that provides payment gateway solutions through bank transfers. It connects directly with various Vietnamese banks via Open Banking APIs, allowing for real-time transaction monitoring and automatic payment confirmation.

The main benefits of using Sepay include:

- Direct bank API integration for real-time transaction monitoring
- Automatic payment confirmation through webhooks
- Support for VietQR code generation
- Connection with multiple Vietnamese banks
- Lower transaction fees compared to traditional payment gateways

This guide covers how to integrate Sepay's payment solution with a Golang application.

## How Sepay Works

The payment flow through Sepay works as follows:

1. **Order Placement**: Customer places an order on your website/application.
2. **QR Code Generation**: Your application displays a payment QR code generated via VietQR.
3. **Customer Payment**: Customer scans the QR code using their banking app and completes the transfer.
4. **Bank Notification**: The bank processes the payment and notifies Sepay.
5. **Webhook Notification**: Sepay sends a webhook to your application with the transaction details.
6. **Order Update**: Your application processes the webhook, verifies the payment, and updates the order status.

## Setting Up a Sepay Account

1. Register for an account at [my.sepay.vn](https://my.sepay.vn)
2. Complete the verification process
3. Add your bank account information
4. Set up webhook configurations in your Sepay dashboard

## Implementation Steps

### 1. Creating Payment QR Codes

There are two ways to generate QR codes for payments:

- **Using Sepay's QR Generator**: Visit [qr.sepay.vn](https://qr.sepay.vn) to generate static or dynamic QR codes.
- **Using VietQR Standard**: Implement the VietQR standard directly in your application. 

For a Golang application, you can implement the VietQR standard using libraries like `github.com/skip2/go-qrcode` to generate QR codes.

### 2. Setting Up Webhook Endpoint

Create an endpoint in your Golang application to receive webhook notifications from Sepay. This will typically be a POST endpoint that processes the incoming transaction data.

### 3. Processing Webhook Data

When Sepay sends a webhook, it will include the following information in a JSON payload:

```json
{
  "id": 92704,                              // Transaction ID on Sepay
  "gateway": "Vietcombank",                 // Bank name
  "transactionDate": "2024-07-25 14:02:37", // Transaction timestamp
  "accountNumber": "0123499999",            // Bank account number
  "amount": 1000000,                        // Transaction amount (in VND)
  "description": "ORDER123",                // Transaction description
  "customerInfo": "Customer Name",          // Customer information
  "creditAmount": 1000000,                  // Credit amount
  "debitAmount": 0,                         // Debit amount
  "fee": 0,                                 // Transaction fee
  "bankTransactionId": "FT22722398374",     // Bank's transaction ID
  "webhookUrl": "https://your-webhook-url"  // Your webhook URL
}
```

Your endpoint should:
1. Receive and parse this JSON data
2. Validate the transaction (match amount, order reference, etc.)
3. Update your database accordingly

### 4. Updating Transaction Status

After validating the payment information, update your order/transaction status in your database to reflect the successful payment.

## Clean Architecture Implementation

To build a maintainable and testable Sepay integration, we recommend following Clean Architecture principles. This approach separates your code into concentric layers, with domain entities at the core and external frameworks/drivers at the outermost layer.

### Project Structure

Here's a recommended project structure that follows Clean Architecture:

```
sepay-integration/
├── cmd/
│   └── server/
│       └── main.go                  # Application entry point
├── internal/
│   ├── domain/                      # Domain Layer
│   │   ├── entity/
│   │   │   ├── order.go
│   │   │   ├── payment.go
│   │   │   └── transaction.go
│   │   ├── repository/
│   │   │   ├── order_repository.go     # Repository interfaces
│   │   │   └── transaction_repository.go
│   │   └── service/
│   │       └── payment_service.go      # Domain service interfaces
│   ├── usecase/                     # Use Case Layer
│   │   ├── create_payment.go
│   │   ├── process_webhook.go
│   │   └── update_transaction.go
│   ├── adapter/                     # Interface Adapters Layer
│   │   ├── api/
│   │   │   ├── handler/
│   │   │   │   ├── order_handler.go
│   │   │   │   └── webhook_handler.go
│   │   │   └── middleware/
│   │   │       └── auth_middleware.go
│   │   ├── repository/
│   │   │   ├── mysql_order_repository.go
│   │   │   └── mysql_transaction_repository.go
│   │   └── qrcode/
│   │       └── vietqr_generator.go
│   └── infrastructure/              # Frameworks & Drivers Layer
│       ├── config/
│       │   └── config.go
│       ├── persistence/
│       │   └── mysql.go
│       └── sepay/
│           └── client.go
└── pkg/                             # Shared packages
    ├── logger/
    │   └── logger.go
    └── errors/
        └── errors.go
```

### Domain Layer

The domain layer contains the business logic and rules of your application. It should be independent of external concerns.

```go
// internal/domain/entity/transaction.go
package entity

import "time"

type TransactionStatus string

const (
    TransactionStatusPending   TransactionStatus = "pending"
    TransactionStatusCompleted TransactionStatus = "completed"
    TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
    ID               string
    OrderID          string
    Amount           int64
    Status           TransactionStatus
    PaymentMethod    string
    PaymentReference string
    BankName         string
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

// internal/domain/repository/transaction_repository.go
package repository

import (
    "context"
    
    "your-module/internal/domain/entity"
)

type TransactionRepository interface {
    Create(ctx context.Context, transaction *entity.Transaction) error
    GetByID(ctx context.Context, id string) (*entity.Transaction, error)
    GetByOrderID(ctx context.Context, orderID string) (*entity.Transaction, error)
    Update(ctx context.Context, transaction *entity.Transaction) error
}
```

### Use Case Layer

The use case layer contains application-specific business rules and orchestrates the flow of data to and from entities.

```go
// internal/usecase/process_webhook.go
package usecase

import (
    "context"
    
    "your-module/internal/domain/entity"
    "your-module/internal/domain/repository"
)

type WebhookData struct {
    ID                int64  
    Gateway           string
    TransactionDate   string
    AccountNumber     string
    Amount            int64
    Description       string // Contains OrderID
    BankTransactionID string
}

type ProcessWebhookUseCase struct {
    transactionRepo repository.TransactionRepository
    orderRepo       repository.OrderRepository
}

func NewProcessWebhookUseCase(
    transactionRepo repository.TransactionRepository,
    orderRepo repository.OrderRepository,
) *ProcessWebhookUseCase {
    return &ProcessWebhookUseCase{
        transactionRepo: transactionRepo,
        orderRepo:       orderRepo,
    }
}

func (uc *ProcessWebhookUseCase) Execute(ctx context.Context, data WebhookData) error {
    // Extract orderID from description
    orderID := data.Description
    
    // Find related transaction
    transaction, err := uc.transactionRepo.GetByOrderID(ctx, orderID)
    if (err != nil) {
        return err
    }
    
    // Update transaction status
    transaction.Status = entity.TransactionStatusCompleted
    transaction.PaymentReference = data.BankTransactionID
    transaction.BankName = data.Gateway
    transaction.UpdatedAt = time.Now()
    
    // Save updated transaction
    if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
        return err
    }
    
    // Update order status
    order, err := uc.orderRepo.GetByID(ctx, orderID)
    if err != nil {
        return err
    }
    
    order.Status = "paid"
    order.UpdatedAt = time.Now()
    
    return uc.orderRepo.Update(ctx, order)
}
```

### Interface Adapters Layer

This layer contains adapters that convert data from the format most convenient for use cases and entities to the format most convenient for external services.

```go
// internal/adapter/api/handler/webhook_handler.go
package handler

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    
    "your-module/internal/usecase"
    "your-module/pkg/errors"
    "your-module/pkg/logger"
)

type WebhookHandler struct {
    processWebhookUseCase *usecase.ProcessWebhookUseCase
    logger                logger.Logger
}

func NewWebhookHandler(
    processWebhookUseCase *usecase.ProcessWebhookUseCase,
    logger logger.Logger,
) *WebhookHandler {
    return &WebhookHandler{
        processWebhookUseCase: processWebhookUseCase,
        logger:                logger,
    }
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Read request body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        h.logger.Error("Error reading webhook body", err)
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    // Parse webhook data
    var webhookData usecase.WebhookData
    if err := json.Unmarshal(body, &webhookData); err != nil {
        h.logger.Error("Error parsing webhook data", err)
        http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
        return
    }
    
    // Process webhook
    if err := h.processWebhookUseCase.Execute(r.Context(), webhookData); err != nil {
        h.logger.Error("Error processing webhook", err)
        http.Error(w, "Failed to process webhook", http.StatusInternalServerError)
        return
    }
    
    // Respond with success
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
```

### Frameworks & Drivers Layer

This outermost layer contains frameworks and tools like databases, web frameworks, etc.

```go
// internal/infrastructure/persistence/mysql.go
package persistence

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
    "your-module/internal/infrastructure/config"
)

func NewMySQLConnection(cfg *config.Config) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", 
        cfg.Database.User, 
        cfg.Database.Password, 
        cfg.Database.Host, 
        cfg.Database.Port, 
        cfg.Database.Name,
    )
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database connection: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
    db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
    
    return db, nil
}
```

## Domain-Driven Design Approach

Following Domain-Driven Design (DDD) principles helps create a more expressive and business-aligned codebase.

### Core Domain Entities

In Sepay integration, the core domain entities include:

1. **Order** - Represents a customer order that needs payment
2. **Transaction** - Represents a payment transaction
3. **Payment** - Represents the payment process

### Value Objects

Use value objects for immutable concepts that don't have an identity:

```go
// internal/domain/entity/value_objects.go
package entity

type Money struct {
    Amount   int64
    Currency string
}

func NewMoney(amount int64, currency string) Money {
    return Money{
        Amount:   amount,
        Currency: currency,
    }
}

type QRCode struct {
    Content string
    Size    int
}

func NewQRCode(content string, size int) QRCode {
    return QRCode{
        Content: content,
        Size:    size,
    }
}
```

### Domain Services

Some business logic doesn't naturally fit within a single entity. Use domain services for these operations:

```go
// internal/domain/service/payment_service.go
package service

import (
    "context"
    
    "your-module/internal/domain/entity"
)

type PaymentService interface {
    GeneratePaymentQR(ctx context.Context, order entity.Order) (entity.QRCode, error)
    VerifyPaymentWebhook(ctx context.Context, webhook entity.WebhookPayload) (bool, error)
}
```

### Repositories

Use repositories to abstract data access:

```go
// internal/domain/repository/order_repository.go
package repository

import (
    "context"
    
    "your-module/internal/domain/entity"
)

type OrderRepository interface {
    Create(ctx context.Context, order *entity.Order) error
    GetByID(ctx context.Context, id string) (*entity.Order, error)
    Update(ctx context.Context, order *entity.Order) error
    List(ctx context.Context, limit, offset int) ([]*entity.Order, error)
}
```

## Code Examples

### Basic Project Structure

Here's a suggested project structure for your Sepay integration with Golang:

```
sepay-integration/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   ├── order.go
│   │   └── webhook.go
│   ├── model/
│   │   ├── order.go
│   │   └── transaction.go
│   └── service/
│       ├── order_service.go
│       ├── payment_service.go
│       └── qr_service.go
├── pkg/
│   └── qrcode/
│       └── vietqr.go
└── go.mod
```

### Creating a QR Code Generator

```go
// pkg/qrcode/vietqr.go
package qrcode

import (
    "fmt"
    "github.com/skip2/go-qrcode"
)

// VietQRData represents the data needed to generate a VietQR code
type VietQRData struct {
    AccountNumber string
    AccountName   string
    BankID        string // Bank BIN code
    Amount        int64
    Description   string
}

// GenerateVietQR creates a VietQR code based on the provided data
func GenerateVietQR(data VietQRData) ([]byte, error) {
    // VietQR format: bankid|account|amount|description
    qrContent := fmt.Sprintf("VietQR|%s|%s|%d|%s", 
        data.BankID,
        data.AccountNumber,
        data.Amount,
        data.Description)
    
    // Generate QR code
    qrCode, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
    if err != nil {
        return nil, fmt.Errorf("failed to generate QR code: %w", err)
    }
    
    return qrCode, nil
}

// SaveQRCodeToFile saves the QR code to a file
func SaveQRCodeToFile(data VietQRData, filepath string) error {
    qrContent := fmt.Sprintf("VietQR|%s|%s|%d|%s", 
        data.BankID,
        data.AccountNumber,
        data.Amount,
        data.Description)
    
    err := qrcode.WriteFile(qrContent, qrcode.Medium, 256, filepath)
    if err != nil {
        return fmt.Errorf("failed to save QR code to file: %w", err)
    }
    
    return nil
}
```

### Setting Up a Webhook Server

```go
// cmd/server/main.go
package main

import (
    "log"
    "net/http"

    "your-module/internal/config"
    "your-module/internal/handler"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize handlers
    webhookHandler := handler.NewWebhookHandler(cfg)
    orderHandler := handler.NewOrderHandler(cfg)

    // Set up routes
    http.HandleFunc("/sepay/webhook", webhookHandler.Handle)
    http.HandleFunc("/orders/create", orderHandler.Create)
    http.HandleFunc("/orders/status", orderHandler.Status)

    // Start server
    serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
    log.Printf("Starting server on %s", serverAddr)
    if err := http.ListenAndServe(serverAddr, nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

### Processing Webhook Payload

```go
// internal/handler/webhook.go
package handler

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    
    "your-module/internal/config"
    "your-module/internal/model"
    "your-module/internal/service"
)

// WebhookHandler processes webhooks from Sepay
type WebhookHandler struct {
    cfg *config.Config
    orderService *service.OrderService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(cfg *config.Config) *WebhookHandler {
    return &WebhookHandler{
        cfg: cfg,
        orderService: service.NewOrderService(),
    }
}

// Handle processes incoming webhook requests
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Read request body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Printf("Error reading request body: %v", err)
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    // Parse webhook payload
    var webhookData model.SepayWebhook
    if err := json.Unmarshal(body, &webhookData); err != nil {
        log.Printf("Error unmarshaling webhook data: %v", err)
        http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
        return
    }
    
    // Validate webhook (optional: add signature verification)
    
    // Extract order ID from description
    // Assuming description format: "ORDER123"
    orderID := webhookData.Description
    
    // Update order status
    if err := h.orderService.UpdateOrderStatus(orderID, "paid"); err != nil {
        log.Printf("Error updating order status: %v", err)
        http.Error(w, "Failed to process webhook", http.StatusInternalServerError)
        return
    }
    
    // Respond with success
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
```

### Database Integration

```go
// internal/model/transaction.go
package model

import "time"

// SepayWebhook represents the webhook data structure from Sepay
type SepayWebhook struct {
    ID                int64     `json:"id"`
    Gateway           string    `json:"gateway"`
    TransactionDate   string    `json:"transactionDate"`
    AccountNumber     string    `json:"accountNumber"`
    Amount            int64     `json:"amount"`
    Description       string    `json:"description"`
    CustomerInfo      string    `json:"customerInfo"`
    CreditAmount      int64     `json:"creditAmount"`
    DebitAmount       int64     `json:"debitAmount"`
    Fee               int64     `json:"fee"`
    BankTransactionID string    `json:"bankTransactionId"`
    WebhookURL        string    `json:"webhookUrl"`
}

// Transaction represents a payment transaction in your system
type Transaction struct {
    ID               string    `json:"id"`
    OrderID          string    `json:"order_id"`
    Amount           int64     `json:"amount"`
    Status           string    `json:"status"`
    PaymentMethod    string    `json:"payment_method"`
    PaymentReference string    `json:"payment_reference"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

// internal/service/order_service.go
package service

import (
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/go-sql-driver/mysql"
    "your-module/internal/model"
)

// OrderService handles order-related operations
type OrderService struct {
    db *sql.DB
}

// NewOrderService creates a new order service
func NewOrderService() *OrderService {
    // Initialize database connection
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to database: %v", err))
    }
    
    return &OrderService{db: db}
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(orderID, status string) error {
    query := `UPDATE orders SET status = ?, updated_at = ? WHERE id = ?`
    _, err := s.db.Exec(query, status, time.Now(), orderID)
    if err != nil {
        return fmt.Errorf("failed to update order status: %w", err)
    }
    
    return nil
}

// RecordTransaction records a new transaction
func (s *OrderService) RecordTransaction(webhook model.SepayWebhook) error {
    // Parse transaction date
    transactionDate, err := time.Parse("2006-01-02 15:04:05", webhook.TransactionDate)
    if err != nil {
        return fmt.Errorf("invalid transaction date format: %w", err)
    }
    
    // Create transaction record
    query := `INSERT INTO transactions 
              (order_id, amount, status, payment_method, payment_reference, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
    
    _, err = s.db.Exec(query, 
        webhook.Description,        // order_id (from description field)
        webhook.Amount,             // amount
        "completed",                // status
        webhook.Gateway,            // payment_method
        webhook.BankTransactionID,  // payment_reference
        transactionDate,            // created_at
        time.Now(),                 // updated_at
    )
    
    if err != nil {
        return fmt.Errorf("failed to record transaction: %w", err)
    }
    
    return nil
}
```

## Testing and Debugging

### Local Testing

1. Use a tool like ngrok to expose your local webhook endpoint to the internet
2. Configure your Sepay account to send webhooks to your exposed endpoint
3. Make test transactions using the Sepay sandbox environment
4. Check the webhook logs in your Sepay dashboard at `my.sepay.vn` > `Transactions` > `View Transaction Details` > `View Webhooks`

### Debugging Tips

- Always log incoming webhook payloads for debugging purposes
- Implement error handling for all transaction processing steps
- Check Sepay's transaction logs for confirmation of webhook delivery
- If a webhook fails, you can trigger a redelivery from the Sepay dashboard

### Table-Driven Testing

Table-driven tests are highly effective for testing multiple scenarios with minimal code duplication. Here's how to implement them for Sepay integration:

```go
// internal/adapter/qrcode/vietqr_generator_test.go
package qrcode_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "your-module/internal/adapter/qrcode"
    "your-module/internal/domain/entity"
)

func TestGenerateVietQR(t *testing.T) {
    tests := []struct {
        name        string
        input       qrcode.VietQRData
        wantContent string
        wantErr     bool
    }{
        {
            name: "valid data with amount",
            input: qrcode.VietQRData{
                BankID:        "970436",
                AccountNumber: "1234567890",
                AccountName:   "Test User",
                Amount:        100000,
                Description:   "ORDER123",
            },
            wantContent: "VietQR|970436|1234567890|100000|ORDER123",
            wantErr:     false,
        },
        {
            name: "valid data without amount",
            input: qrcode.VietQRData{
                BankID:        "970436",
                AccountNumber: "1234567890",
                AccountName:   "Test User",
                Amount:        0,
                Description:   "ORDER123",
            },
            wantContent: "VietQR|970436|1234567890|0|ORDER123",
            wantErr:     false,
        },
        {
            name: "missing bank ID",
            input: qrcode.VietQRData{
                BankID:        "",
                AccountNumber: "1234567890",
                AccountName:   "Test User",
                Amount:        100000,
                Description:   "ORDER123",
            },
            wantContent: "",
            wantErr:     true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            generator := qrcode.NewVietQRGenerator()
            got, err := generator.Generate(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.wantContent, got.Content)
        })
    }
}
```

### Mocking External Dependencies

Use interfaces and mocks for external dependencies to facilitate testing:

```go
// internal/usecase/process_webhook_test.go
package usecase_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "your-module/internal/domain/entity"
    "your-module/internal/domain/repository/mocks"
    "your-module/internal/usecase"
)

func TestProcessWebhookUseCase(t *testing.T) {
    tests := []struct {
        name           string
        webhookData    usecase.WebhookData
        setupMocks     func(*mocks.TransactionRepository, *mocks.OrderRepository)
        expectedError  bool
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
                BankTransactionID: "FT12345678",
            },
            setupMocks: func(tr *mocks.TransactionRepository, or *mocks.OrderRepository) {
                // Setup transaction repository mock
                tr.On("GetByOrderID", mock.Anything, "ORDER123").Return(&entity.Transaction{
                    ID:      "txn-123",
                    OrderID: "ORDER123",
                    Status:  entity.TransactionStatusPending,
                }, nil)
                tr.On("Update", mock.Anything, mock.AnythingOfType("*entity.Transaction")).Return(nil)
                
                // Setup order repository mock
                or.On("GetByID", mock.Anything, "ORDER123").Return(&entity.Order{
                    ID:     "ORDER123",
                    Status: "pending",
                }, nil)
                or.On("Update", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)
            },
            expectedError: false,
        },
        {
            name: "transaction not found",
            webhookData: usecase.WebhookData{
                Description: "ORDER999",
            },
            setupMocks: func(tr *mocks.TransactionRepository, or *mocks.OrderRepository) {
                tr.On("GetByOrderID", mock.Anything, "ORDER999").Return(nil, errors.New("transaction not found"))
            },
            expectedError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mocks
            transactionRepo := new(mocks.TransactionRepository)
            orderRepo := new(mocks.OrderRepository)
            
            // Setup mocks
            tt.setupMocks(transactionRepo, orderRepo)
            
            // Create use case
            uc := usecase.NewProcessWebhookUseCase(transactionRepo, orderRepo)
            
            // Execute use case
            err := uc.Execute(context.Background(), tt.webhookData)
            
            // Assert result
            if tt.expectedError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            
            // Verify mocks
            transactionRepo.AssertExpectations(t)
            orderRepo.AssertExpectations(t)
        })
    }
}
```

### Integration Testing

For integration tests with Sepay, use a controlled environment:

```go
// tests/integration/webhook_test.go
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "your-module/internal/adapter/api/handler"
)

func TestWebhookIntegration(t *testing.T) {
    // Skip if not running integration tests
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Create a test server with your webhook handler
    server := setupTestServer()
    
    tests := []struct {
        name         string
        webhookJSON  string
        setupData    func()
        wantStatus   int
        assertResult func(t *testing.T)
    }{
        {
            name: "valid webhook updates order status",
            webhookJSON: `{
                "id": 12345,
                "gateway": "Vietcombank",
                "transactionDate": "2025-05-02 14:30:00",
                "accountNumber": "0123456789",
                "amount": 150000,
                "description": "ORDER123",
                "customerInfo": "Test Customer",
                "creditAmount": 150000,
                "debitAmount": 0,
                "fee": 0,
                "bankTransactionId": "FT12345678",
                "webhookUrl": "https://api.yourdomain.com/sepay/webhook"
            }`,
            setupData: func() {
                // Insert test order and transaction to your test database
                // ...
            },
            wantStatus: http.StatusOK,
            assertResult: func(t *testing.T) {
                // Verify order status was updated in the database
                // ...
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Set up test data
            tt.setupData()
            
            // Make the request
            req := httptest.NewRequest("POST", "/sepay/webhook", bytes.NewBufferString(tt.webhookJSON))
            req.Header.Set("Content-Type", "application/json")
            
            // Record the response
            res := httptest.NewRecorder()
            server.ServeHTTP(res, req)
            
            // Check status code
            assert.Equal(t, tt.wantStatus, res.Code)
            
            // Run additional assertions
            tt.assertResult(t)
        })
    }
}

func setupTestServer() http.Handler {
    // Set up test dependencies and create the HTTP handler
    // ...
}
```

## Security Considerations

1. **Authenticate Incoming Webhooks**: Implement API key authentication or signature verification for webhooks

2. **HTTPS**: Always use HTTPS for your webhook endpoints

3. **Idempotency**: Ensure your webhook processing is idempotent to handle potential duplicate webhook deliveries

4. **Data Validation**: Always validate incoming webhook data, especially transaction amounts and order references

5. **Access Control**: Limit access to your webhook endpoint to Sepay's IP addresses if possible

## Advanced Configuration

### API Key Authentication

To add API key authentication to your webhook endpoint:

```go
// Inside your webhook handler
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
    // Check for API key in headers
    apiKey := r.Header.Get("X-Api-Key")
    if apiKey != h.cfg.SepayAPIKey {
        log.Printf("Invalid API key: %s", apiKey)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Rest of your handler code...
}
```

### Transaction Reconciliation

Implement a periodic reconciliation service to catch any missed transactions:

```go
func startReconciliationService(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            performReconciliation()
        }
    }
}

func performReconciliation() {
    // Get pending orders from your database
    
    // For each pending order:
    // 1. Check its status with Sepay API
    // 2. Update local database if necessary
}
```

## Reference

- Sepay Website: [https://sepay.vn](https://sepay.vn)
- Sepay Dashboard: [https://my.sepay.vn](https://my.sepay.vn)
- QR Generator: [https://qr.sepay.vn](https://qr.sepay.vn)
- VietQR Official Website: [https://vietqr.net](https://vietqr.net)

---

This document was created on May 2, 2025, and represents the latest integration guidelines for Sepay payment gateway with Golang. Please check Sepay's official documentation for any updates or changes to their API.