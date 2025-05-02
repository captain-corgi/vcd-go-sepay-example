package middleware

import (
	"net/http"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/infrastructure/config"
	"github.com/labstack/echo/v4"
)

// APIKeyAuth is middleware for validating API keys
type APIKeyAuth struct {
	config *config.Config
}

// NewAPIKeyAuth creates a new API key authentication middleware
func NewAPIKeyAuth(config *config.Config) *APIKeyAuth {
	return &APIKeyAuth{
		config: config,
	}
}

// Middleware returns an Echo middleware function for API key authentication
func (m *APIKeyAuth) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get API key from header
			apiKey := c.Request().Header.Get("X-API-Key")

			// Skip auth for webhook endpoint (webhook has its own auth check)
			if c.Path() == "/api/sepay/webhook" {
				return next(c)
			}

			// Check if API key is valid
			if apiKey != m.config.Sepay.APIKey {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: invalid API key",
				})
			}

			// Continue with the next handler
			return next(c)
		}
	}
}
