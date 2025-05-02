package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Sepay    SepayConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port            int
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver       string
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
}

// SepayConfig holds Sepay-specific configuration
type SepayConfig struct {
	APIKey         string
	BankID         string
	AccountNumber  string
	AccountName    string
	WebhookSecret  string
	WebhookBaseURL string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))
	shutdownTimeout, _ := strconv.Atoi(getEnv("SERVER_SHUTDOWN_TIMEOUT", "5"))
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "10"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))

	return &Config{
		Server: ServerConfig{
			Port:            serverPort,
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			ShutdownTimeout: shutdownTimeout,
		},
		Database: DatabaseConfig{
			Driver:       getEnv("DB_DRIVER", "mysql"),
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         dbPort,
			User:         getEnv("DB_USER", "root"),
			Password:     getEnv("DB_PASSWORD", ""),
			Name:         getEnv("DB_NAME", "sepay"),
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
		},
		Sepay: SepayConfig{
			APIKey:         getEnv("SEPAY_API_KEY", ""),
			BankID:         getEnv("SEPAY_BANK_ID", ""),
			AccountNumber:  getEnv("SEPAY_ACCOUNT_NUMBER", ""),
			AccountName:    getEnv("SEPAY_ACCOUNT_NAME", ""),
			WebhookSecret:  getEnv("SEPAY_WEBHOOK_SECRET", ""),
			WebhookBaseURL: getEnv("SEPAY_WEBHOOK_BASE_URL", "https://api.example.com"),
		},
	}, nil
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
