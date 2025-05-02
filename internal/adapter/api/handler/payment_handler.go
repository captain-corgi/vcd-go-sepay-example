package handler

import (
	"net/http"
	"strconv"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/usecase"
	"github.com/labstack/echo/v4"
)

// PaymentHandler handles payment related HTTP requests
type PaymentHandler struct {
	generatePaymentQRUseCase *usecase.GeneratePaymentQRUseCase
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(generatePaymentQRUseCase *usecase.GeneratePaymentQRUseCase) *PaymentHandler {
	return &PaymentHandler{
		generatePaymentQRUseCase: generatePaymentQRUseCase,
	}
}

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	CustomerID  string `json:"customer_id"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

// CreatePaymentResponse represents the response for a payment creation
type CreatePaymentResponse struct {
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	QRContent     string `json:"qr_content"`
	QRImage       string `json:"qr_image"` // Base64 encoded image
	ExpiresAt     string `json:"expires_at"`
	BankID        string `json:"bank_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

// CreatePayment handles requests to create a new payment
func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	var req CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request: " + err.Error(),
		})
	}

	// Validate request
	if req.CustomerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Customer ID is required",
		})
	}

	if req.Amount == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Amount is required",
		})
	}

	// Parse amount
	amount, err := strconv.ParseInt(req.Amount, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid amount: " + err.Error(),
		})
	}

	if amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Amount must be greater than zero",
		})
	}

	// Create order input
	input := usecase.OrderInput{
		CustomerID:  req.CustomerID,
		Amount:      amount,
		Description: req.Description,
	}

	// Generate payment QR
	output, err := h.generatePaymentQRUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate payment QR: " + err.Error(),
		})
	}

	// Convert QR code image to base64
	qrImageBase64 := ""
	if output.QRCode != nil && len(output.QRCode.Image) > 0 {
		qrImageBase64 = "/9j/4AAQSkZJRgABAQAAAQABAAD..." // In real implementation, encode the QR image to base64
	}

	// Create response
	response := CreatePaymentResponse{
		OrderID:       output.OrderID,
		Amount:        output.Amount,
		QRContent:     output.QRCode.Content,
		QRImage:       qrImageBase64,
		ExpiresAt:     output.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		BankID:        output.BankID,
		BankName:      output.BankName,
		AccountNumber: output.AccountNumber,
		AccountName:   output.AccountName,
	}

	return c.JSON(http.StatusOK, response)
}

// GetPaymentStatus handles requests to check payment status
func (h *PaymentHandler) GetPaymentStatus(c echo.Context) error {
	orderID := c.QueryParam("order_id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Order ID is required",
		})
	}

	// In a real implementation, we would check the status of the order
	// For now, we'll return a placeholder response
	return c.JSON(http.StatusOK, map[string]string{
		"order_id": orderID,
		"status":   "pending", // In a real implementation, get this from database
	})
}

// RegisterRoutes registers the payment handler routes
func (h *PaymentHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/api/payments", h.CreatePayment)
	e.GET("/api/payments/status", h.GetPaymentStatus)
}
