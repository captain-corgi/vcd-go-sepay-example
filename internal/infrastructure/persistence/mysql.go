package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/infrastructure/config"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLDB represents a MySQL database connection
type MySQLDB struct {
	DB *sql.DB
}

// NewMySQLConnection creates a new MySQL database connection
func NewMySQLConnection(cfg *config.Config) (*MySQLDB, error) {
	// Format: username:password@tcp(host:port)/dbname?parseTime=true&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
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

	// Set connection pool configuration
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &MySQLDB{DB: db}, nil
}

// Close closes the database connection
func (m *MySQLDB) Close() error {
	if m.DB != nil {
		return m.DB.Close()
	}
	return nil
}
