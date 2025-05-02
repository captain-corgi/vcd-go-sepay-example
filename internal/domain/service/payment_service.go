package service

import (
	"context"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
)

// PaymentService defines the interface for payment operations
type PaymentService interface {
	// GeneratePaymentQR generates a QR code for payment based on order data
	GeneratePaymentQR(ctx context.Context, order *entity.Order) (*entity.QRCode, error)

	// ProcessWebhook processes an incoming webhook from Sepay
	ProcessWebhook(ctx context.Context, payload *entity.WebhookPayload) error

	// VerifyPayment verifies if a payment is valid
	VerifyPayment(ctx context.Context, orderID string, amount int64) (bool, error)
}
