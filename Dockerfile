FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o sepay-service ./cmd/server

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/sepay-service .

# Create directory for certificates
RUN mkdir -p /app/certs

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./sepay-service"]