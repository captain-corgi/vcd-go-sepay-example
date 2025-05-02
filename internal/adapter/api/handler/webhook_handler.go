package handler

import (
	"encoding/json"
	"net/http"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/infrastructure/config"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/usecase"
	"github.com/labstack/echo/v4"
)

// WebhookHandler handles webhook related HTTP requests
type WebhookHandler struct {
	processWebhookUseCase *usecase.ProcessWebhookUseCase
	config                *config.Config
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(processWebhookUseCase *usecase.ProcessWebhookUseCase, config *config.Config) *WebhookHandler {
	return &WebhookHandler{
		processWebhookUseCase: processWebhookUseCase,
		config:                config,
	}
}

// HandleWebhook handles incoming webhook requests from Sepay
func (h *WebhookHandler) HandleWebhook(c echo.Context) error {
	// Parse webhook data
	var webhookPayload entity.WebhookPayload
	if err := json.NewDecoder(c.Request().Body).Decode(&webhookPayload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid webhook payload: " + err.Error(),
		})
	}

	// Process webhook
	if err := h.processWebhookUseCase.Execute(c.Request().Context(), &webhookPayload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process webhook: " + err.Error(),
		})
	}

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

// RegisterRoutes registers the webhook handler routes
func (h *WebhookHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/api/sepay/webhook", h.HandleWebhook)
}
