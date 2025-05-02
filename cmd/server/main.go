package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/adapter/api/handler"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/adapter/api/middleware"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/adapter/qrcode"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/adapter/repository"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/infrastructure/config"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/infrastructure/persistence"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/usecase"
	"github.com/captain-corgi/vcd-go-sepay-example/pkg/logger"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize logger
	log := logger.NewJSONLogger("sepay-integration")
	log.Info("Starting Sepay integration service")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", err)
	}

	// Initialize database
	log.Info("Connecting to database")
	db, err := persistence.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", err, map[string]interface{}{
			"db_host": cfg.Database.Host,
			"db_name": cfg.Database.Name,
		})
	}
	defer db.Close()
	log.Info("Database connection established")

	// Initialize repositories
	orderRepo := repository.NewMySQLOrderRepository(db)
	transactionRepo := repository.NewMySQLTransactionRepository(db)

	// Initialize QR code generator
	qrGenerator := qrcode.NewVietQRGenerator(256)

	// Initialize use cases
	generatePaymentQRUseCase := usecase.NewGeneratePaymentQRUseCase(orderRepo, transactionRepo, qrGenerator, cfg)
	processWebhookUseCase := usecase.NewProcessWebhookUseCase(orderRepo, transactionRepo)

	// Initialize handlers
	paymentHandler := handler.NewPaymentHandler(generatePaymentQRUseCase)
	webhookHandler := handler.NewWebhookHandler(processWebhookUseCase, cfg)

	// Initialize Echo server
	e := echo.New()

	// Add middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Add API key authentication middleware
	apiKeyAuth := middleware.NewAPIKeyAuth(cfg)
	e.Use(apiKeyAuth.Middleware())

	// Register routes
	paymentHandler.RegisterRoutes(e)
	webhookHandler.RegisterRoutes(e)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "UP",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Start server
	go func() {
		address := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Info("Starting HTTP server", map[string]interface{}{
			"port": cfg.Server.Port,
		})

		// Configure TLS if in production
		if os.Getenv("GO_ENV") == "production" {
			if err := e.StartTLS(address, "certs/server.crt", "certs/server.key"); err != nil && err != http.ErrServerClosed {
				log.Fatal("Failed to start HTTPS server", err)
			}
		} else {
			if err := e.Start(address); err != nil && err != http.ErrServerClosed {
				log.Fatal("Failed to start HTTP server", err)
			}
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed", err)
	}

	log.Info("Server gracefully stopped")
}
