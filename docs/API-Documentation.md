# Sepay Payment Gateway API Documentation

This document describes the API endpoints available for integrating with the Sepay payment gateway.

## Base URL

For development: `http://localhost:8080`  
For production: `https://api.your-domain.com`

## Authentication

Most API endpoints require authentication via an API key. Include this key in the `X-API-Key` header for requests.

```
X-API-Key: your_api_key_here
```

## Endpoints

### Create Payment

Generates a payment QR code for a new order.

**URL:** `/api/payments`

**Method:** `POST`

**Authentication Required:** Yes

**Request Body:**

```json
{
  "customer_id": "customer_123",
  "amount": "1000000",
  "description": "Payment for order #12345"
}
```

**Successful Response (200 OK):**

```json
{
  "order_id": "ord_abcd1234",
  "amount": 1000000,
  "qr_content": "VietQR|970436|1234567890|1000000|ord_abcd1234",
  "qr_image": "base64_encoded_image_data",
  "expires_at": "2025-05-03T12:34:56+07:00",
  "bank_id": "970436",
  "bank_name": "Vietcombank",
  "account_number": "1234567890",
  "account_name": "COMPANY NAME"
}
```

**Error Response (400 Bad Request):**

```json
{
  "error": "Amount is required"
}
```

### Check Payment Status

Checks the status of a payment.

**URL:** `/api/payments/status?order_id=ord_abcd1234`

**Method:** `GET`

**Authentication Required:** Yes

**Query Parameters:**

- `order_id`: The ID of the order to check (required)

**Successful Response (200 OK):**

```json
{
  "order_id": "ord_abcd1234",
  "status": "paid",
  "paid_at": "2025-05-02T15:30:45+07:00",
  "payment_reference": "FT12345678",
  "bank_name": "Vietcombank"
}
```

**Error Response (404 Not Found):**

```json
{
  "error": "Order not found"
}
```

### Sepay Webhook Endpoint

Endpoint that receives payment notifications from Sepay.

**URL:** `/api/sepay/webhook`

**Method:** `POST`

**Authentication Required:** Yes (via X-API-Key header)

**Request Body Example:**

```json
{
  "id": 92704,
  "gateway": "Vietcombank",
  "transactionDate": "2025-05-02 14:02:37",
  "accountNumber": "0123499999",
  "amount": 1000000,
  "description": "ord_abcd1234",
  "customerInfo": "Customer Name",
  "creditAmount": 1000000,
  "debitAmount": 0,
  "fee": 0,
  "bankTransactionId": "FT22722398374",
  "webhookUrl": "https://api.your-domain.com/api/sepay/webhook"
}
```

**Successful Response (200 OK):**

```json
{
  "status": "success"
}
```

**Error Response (401 Unauthorized):**

```json
{
  "error": "Unauthorized: invalid API key"
}
```

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - The request was successful |
| 400 | Bad Request - The request was invalid or cannot be served |
| 401 | Unauthorized - Authentication failed or not provided |
| 404 | Not Found - The requested resource could not be found |
| 500 | Internal Server Error - Something went wrong on the server |

## Integration Flow

1. **Create Payment**:
   - Call the Create Payment endpoint with customer details and amount
   - Receive a QR code and order ID
   - Display the QR code to the customer for scanning with their banking app

2. **Process Payment**:
   - Customer scans the QR code and completes the payment through their banking app
   - Sepay detects the payment and sends a notification to your webhook endpoint
   - Your system processes the webhook, updates the order status, and performs any additional business logic

3. **Verify Payment Status**:
   - Optionally, check the payment status using the Check Payment Status endpoint
   - Use this for reconciliation or if webhook delivery fails

## Testing

For testing purposes, use the following credentials:

- API Key: `test_api_key`
- Bank ID: `970436` (Vietcombank test account)
- Account Number: `1234567890`
- Account Name: `TEST ACCOUNT`

## Webhook Security

To ensure the security of webhook communications, we recommend:

1. Validating the API key in the X-API-Key header
2. Verifying that the payment amount matches your expected order amount
3. Idempotent processing to handle duplicate webhook deliveries
4. Responding quickly to webhook requests (under 5 seconds) to prevent timeouts

## Rate Limits

The API enforces the following rate limits:

- Create Payment: 60 requests per minute
- Check Payment Status: 120 requests per minute
- Webhook: No limit (to ensure all payment notifications are received)

## Support

For integration support, please contact:
- Email: support@your-company.com
- Phone: +84 123 456 789