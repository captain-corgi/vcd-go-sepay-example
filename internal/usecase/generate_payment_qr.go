package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/repository"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/infrastructure/config"
	"github.com/google/uuid"
)

// QRCodeGenerator defines the interface for generating QR codes
type QRCodeGenerator interface {
	Generate(data entity.VietQRData) (*entity.QRCode, error)
}

// GeneratePaymentQRUseCase handles generation of payment QR codes
type GeneratePaymentQRUseCase struct {
	orderRepo       repository.OrderRepository
	transactionRepo repository.TransactionRepository
	qrGenerator     QRCodeGenerator
	config          *config.Config
}

// NewGeneratePaymentQRUseCase creates a new payment QR code generation use case
func NewGeneratePaymentQRUseCase(
	orderRepo repository.OrderRepository,
	transactionRepo repository.TransactionRepository,
	qrGenerator QRCodeGenerator,
	config *config.Config,
) *GeneratePaymentQRUseCase {
	return &GeneratePaymentQRUseCase{
		orderRepo:       orderRepo,
		transactionRepo: transactionRepo,
		qrGenerator:     qrGenerator,
		config:          config,
	}
}

// OrderInput represents input data for creating an order
type OrderInput struct {
	CustomerID  string
	Amount      int64
	Description string
}

// PaymentQROutput represents the output data after QR code generation
type PaymentQROutput struct {
	OrderID       string
	QRCode        *entity.QRCode
	Amount        int64
	ExpiresAt     time.Time
	BankID        string
	BankName      string
	AccountNumber string
	AccountName   string
}

// Execute generates a payment QR code for an order
func (uc *GeneratePaymentQRUseCase) Execute(ctx context.Context, input OrderInput) (*PaymentQROutput, error) {
	// Generate unique order ID
	orderID := fmt.Sprintf("ord_%s", uuid.New().String()[:8])

	// Create new order
	order := &entity.Order{
		ID:          orderID,
		CustomerID:  input.CustomerID,
		Amount:      input.Amount,
		Status:      entity.OrderStatusPending,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save order to database
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create pending transaction
	transaction := &entity.Transaction{
		ID:            fmt.Sprintf("txn_%s", uuid.New().String()[:8]),
		OrderID:       orderID,
		Amount:        input.Amount,
		Status:        entity.TransactionStatusPending,
		PaymentMethod: "bank_transfer",
		Description:   orderID, // Include order ID in description for webhook matching
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Generate VietQR data
	vietQRData := entity.VietQRData{
		BankID:        uc.config.Sepay.BankID,
		AccountNumber: uc.config.Sepay.AccountNumber,
		AccountName:   uc.config.Sepay.AccountName,
		Amount:        input.Amount,
		Description:   orderID, // Use order ID as description for matching with webhook
	}

	// Generate QR code
	qrCode, err := uc.qrGenerator.Generate(vietQRData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Construct response
	expiresAt := time.Now().Add(24 * time.Hour) // QR code expires in 24 hours
	output := &PaymentQROutput{
		OrderID:       orderID,
		QRCode:        qrCode,
		Amount:        input.Amount,
		ExpiresAt:     expiresAt,
		BankID:        uc.config.Sepay.BankID,
		AccountNumber: uc.config.Sepay.AccountNumber,
		AccountName:   uc.config.Sepay.AccountName,
	}

	return output, nil
}
