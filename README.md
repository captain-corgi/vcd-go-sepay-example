# Sepay Payment Gateway Integration

A robust Go backend service for integrating with Sepay payment gateway, providing real-time transaction monitoring and automatic payment confirmation through webhooks.

## Features

- VietQR code generation for bank transfers
- Webhook endpoint for receiving payment notifications
- MySQL database integration for order and transaction storage
- Clean Architecture implementation with DDD principles
- Secure authentication with API key validation
- Comprehensive test suite
- CI/CD pipeline setup with GitHub Actions

## Prerequisites

- Go 1.24 or higher
- MySQL 5.7 or higher
- Docker (optional, for containerization)

## Getting Started

### Installation

1. Clone the repository:

```bash
git clone https://github.com/captain-corgi/vcd-go-sepay-example.git
cd vcd-go-sepay-example
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
# Server configuration
export SERVER_PORT=8080
export SERVER_READ_TIMEOUT=10
export SERVER_WRITE_TIMEOUT=10
export SERVER_SHUTDOWN_TIMEOUT=5

# Database configuration
export DB_DRIVER=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=sepay
export DB_MAX_OPEN_CONNS=10
export DB_MAX_IDLE_CONNS=5

# Sepay configuration
export SEPAY_API_KEY=your_api_key_here
export SEPAY_BANK_ID=your_bank_id_here
export SEPAY_ACCOUNT_NUMBER=your_account_number_here
export SEPAY_ACCOUNT_NAME=your_account_name_here
export SEPAY_WEBHOOK_SECRET=your_webhook_secret_here
export SEPAY_WEBHOOK_BASE_URL=https://api.example.com
```

4. Create database tables:

```sql
CREATE TABLE orders (
    id VARCHAR(36) PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    description TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE transactions (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    payment_reference VARCHAR(100),
    bank_name VARCHAR(100),
    description TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE INDEX idx_transactions_order_id ON transactions(order_id);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
```

### Running the Application

1. Build the application:

```bash
go build -o sepay-service ./cmd/server
```

2. Run the server:

```bash
./sepay-service
```

The server will start at http://localhost:8080 (or the port specified in your environment variables).

### Docker Deployment

1. Build the Docker image:

```bash
docker build -t sepay-integration .
```

2. Run the container:

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e SEPAY_API_KEY=your_api_key_here \
  -e SEPAY_BANK_ID=your_bank_id_here \
  -e SEPAY_ACCOUNT_NUMBER=your_account_number_here \
  -e SEPAY_ACCOUNT_NAME=your_account_name_here \
  sepay-integration
```

## API Documentation

See [API Documentation](docs/API-Documentation.md) for detailed information about the available endpoints.

## Sepay Integration Guide

See [Sepay Guideline](docs/Sepay-Guideline.md) for comprehensive information about integrating with Sepay payment gateway.

## Architecture

This project follows Clean Architecture principles with a domain-driven design approach:

- **Domain Layer**: Core business logic and entities
- **Use Case Layer**: Application-specific business rules
- **Interface Adapters Layer**: Adapters for external services and repositories
- **Frameworks & Drivers Layer**: External frameworks and tools

```
sepay-integration/
├── cmd/
│   └── server/
│       └── main.go                  # Application entry point
├── internal/
│   ├── domain/                      # Domain Layer
│   │   ├── entity/
│   │   ├── repository/
│   │   └── service/
│   ├── usecase/                     # Use Case Layer
│   ├── adapter/                     # Interface Adapters Layer
│   │   ├── api/
│   │   ├── repository/
│   │   └── qrcode/
│   └── infrastructure/              # Frameworks & Drivers Layer
│       ├── config/
│       ├── persistence/
│       └── sepay/
└── pkg/                             # Shared packages
    ├── logger/
    └── errors/
```

## Testing

Run the tests:

```bash
go test ./...
```

For test coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin feature/my-new-feature`)
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgments

- [Sepay](https://sepay.vn) for the payment gateway integration
- [Echo Framework](https://echo.labstack.com/) for the web framework
- [Go-QRCode](https://github.com/skip2/go-qrcode) for QR code generation
