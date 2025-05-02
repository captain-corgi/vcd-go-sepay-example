package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/repository"
)

// ProcessWebhookUseCase handles processing of Sepay webhook data
type ProcessWebhookUseCase struct {
	orderRepo       repository.OrderRepository
	transactionRepo repository.TransactionRepository
}

// NewProcessWebhookUseCase creates a new webhook processing use case
func NewProcessWebhookUseCase(orderRepo repository.OrderRepository, transactionRepo repository.TransactionRepository) *ProcessWebhookUseCase {
	return &ProcessWebhookUseCase{
		orderRepo:       orderRepo,
		transactionRepo: transactionRepo,
	}
}

// Execute processes a webhook payload and updates relevant entities
func (uc *ProcessWebhookUseCase) Execute(ctx context.Context, payload *entity.WebhookPayload) error {
	// Extract order ID from description
	orderID := payload.GetOrderID()
	if orderID == "" {
		return fmt.Errorf("invalid webhook: order ID not found in description")
	}

	// Get order by ID
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}

	// Verify payment amount
	if payload.Amount != order.Amount {
		return fmt.Errorf("payment amount mismatch: expected %d, got %d", order.Amount, payload.Amount)
	}

	// Find existing transaction or create a new one
	transaction, err := uc.transactionRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		// Create new transaction if not exists
		transaction = &entity.Transaction{
			ID:            fmt.Sprintf("txn_%d", payload.ID), // Use Sepay's transaction ID
			OrderID:       orderID,
			Amount:        payload.Amount,
			Status:        entity.TransactionStatusPending,
			PaymentMethod: "bank_transfer",
			Description:   payload.Description,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	}

	// Update transaction with webhook data
	transaction.Status = entity.TransactionStatusCompleted
	transaction.PaymentReference = payload.BankTransactionID
	transaction.BankName = payload.Gateway
	transaction.UpdatedAt = time.Now()

	if transaction.ID == "" {
		// Create new transaction
		err = uc.transactionRepo.Create(ctx, transaction)
	} else {
		// Update existing transaction
		err = uc.transactionRepo.Update(ctx, transaction)
	}

	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	// Update order status
	order.Status = entity.OrderStatusPaid
	order.UpdatedAt = time.Now()
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}
