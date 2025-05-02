package entity

// WebhookPayload represents the incoming webhook data from Sepay
type WebhookPayload struct {
	ID                int64  `json:"id"`                // Transaction ID on Sepay
	Gateway           string `json:"gateway"`           // Bank name
	TransactionDate   string `json:"transactionDate"`   // Transaction timestamp
	AccountNumber     string `json:"accountNumber"`     // Bank account number
	Amount            int64  `json:"amount"`            // Transaction amount (in VND)
	Description       string `json:"description"`       // Transaction description (contains OrderID)
	CustomerInfo      string `json:"customerInfo"`      // Customer information
	CreditAmount      int64  `json:"creditAmount"`      // Credit amount
	DebitAmount       int64  `json:"debitAmount"`       // Debit amount
	Fee               int64  `json:"fee"`               // Transaction fee
	BankTransactionID string `json:"bankTransactionId"` // Bank's transaction ID
	WebhookURL        string `json:"webhookUrl"`        // Your webhook URL
}

// GetOrderID extracts the order ID from the description field
func (wp *WebhookPayload) GetOrderID() string {
	// In a real implementation, you might want to apply additional parsing logic
	// if the description has a specific format. For simplicity, we assume the
	// description directly contains the order ID.
	return wp.Description
}
