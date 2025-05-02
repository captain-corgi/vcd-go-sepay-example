package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Environment constants
const (
	EnvDevelopment       = "DEV"
	EnvSystemIntegration = "SIT"
	EnvUserAcceptance    = "UAT"
	EnvNonFunctional     = "NFT"
	EnvProduction        = "PRD"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Sepay       SepayConfig
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

// Load reads configuration from environment variables and config files
func Load() (*Config, error) {
	// Initialize viper
	v := viper.New()

	// Determine which environment to use
	env := strings.ToUpper(os.Getenv("APP_ENV"))
	if env == "" {
		env = EnvDevelopment // Default to development if not specified
	}

	// Set up config paths
	configPath := getConfigPath()
	envFile := fmt.Sprintf(".env.%s", strings.ToLower(env))
	fullPath := filepath.Join(configPath, envFile)

	// Load environment variables from .env file
	err := loadEnvFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load environment file: %w", err)
	}

	// Set up viper defaults
	setDefaults(v)

	// Bind environment variables to viper
	bindEnvVariables(v)

	// Create configuration
	config := &Config{
		Environment: env,
		Server: ServerConfig{
			Port:            v.GetInt("server.port"),
			ReadTimeout:     v.GetInt("server.read_timeout"),
			WriteTimeout:    v.GetInt("server.write_timeout"),
			ShutdownTimeout: v.GetInt("server.shutdown_timeout"),
		},
		Database: DatabaseConfig{
			Driver:       v.GetString("database.driver"),
			Host:         v.GetString("database.host"),
			Port:         v.GetInt("database.port"),
			User:         v.GetString("database.user"),
			Password:     v.GetString("database.password"),
			Name:         v.GetString("database.name"),
			MaxOpenConns: v.GetInt("database.max_open_conns"),
			MaxIdleConns: v.GetInt("database.max_idle_conns"),
		},
		Sepay: SepayConfig{
			APIKey:         v.GetString("sepay.api_key"),
			BankID:         v.GetString("sepay.bank_id"),
			AccountNumber:  v.GetString("sepay.account_number"),
			AccountName:    v.GetString("sepay.account_name"),
			WebhookSecret:  v.GetString("sepay.webhook_secret"),
			WebhookBaseURL: v.GetString("sepay.webhook_base_url"),
		},
	}

	return config, nil
}

// getConfigPath returns the path to the config directory
func getConfigPath() string {
	// Try to find the config directory relative to the working directory
	baseDir, err := os.Getwd()
	if err != nil {
		return "config/env" // fallback to default path
	}

	// Check if the config directory exists
	configPath := filepath.Join(baseDir, "config", "env")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try parent directory (in case we're in a subdirectory)
		parentDir := filepath.Dir(baseDir)
		configPath = filepath.Join(parentDir, "config", "env")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return "config/env" // fallback to default path
		}
	}

	return configPath
}

// loadEnvFile loads the environment variables from the specified file
func loadEnvFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("environment file not found: %s", path)
	}

	return godotenv.Load(path)
}

// setDefaults sets default values for configuration options
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 10)
	v.SetDefault("server.write_timeout", 10)
	v.SetDefault("server.shutdown_timeout", 5)

	// Database defaults
	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.max_open_conns", 10)
	v.SetDefault("database.max_idle_conns", 5)
}

// bindEnvVariables binds environment variables to viper configuration
func bindEnvVariables(v *viper.Viper) {
	// Server configuration
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	v.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	v.BindEnv("server.shutdown_timeout", "SERVER_SHUTDOWN_TIMEOUT")

	// Database configuration
	v.BindEnv("database.driver", "DB_DRIVER")
	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.user", "DB_USER")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.name", "DB_NAME")
	v.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	v.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")

	// Sepay configuration
	v.BindEnv("sepay.api_key", "SEPAY_API_KEY")
	v.BindEnv("sepay.bank_id", "SEPAY_BANK_ID")
	v.BindEnv("sepay.account_number", "SEPAY_ACCOUNT_NUMBER")
	v.BindEnv("sepay.account_name", "SEPAY_ACCOUNT_NAME")
	v.BindEnv("sepay.webhook_secret", "SEPAY_WEBHOOK_SECRET")
	v.BindEnv("sepay.webhook_base_url", "SEPAY_WEBHOOK_BASE_URL")
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == EnvDevelopment
}

// IsSystemIntegration returns true if the application is running in system integration mode
func (c *Config) IsSystemIntegration() bool {
	return c.Environment == EnvSystemIntegration
}

// IsUserAcceptance returns true if the application is running in user acceptance mode
func (c *Config) IsUserAcceptance() bool {
	return c.Environment == EnvUserAcceptance
}

// IsNonFunctional returns true if the application is running in non-functional testing mode
func (c *Config) IsNonFunctional() bool {
	return c.Environment == EnvNonFunctional
}
