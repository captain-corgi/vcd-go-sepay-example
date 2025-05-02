package entity

// QRCode represents a payment QR code
type QRCode struct {
	Content string `json:"content"`
	Size    int    `json:"size"`
	Image   []byte `json:"image,omitempty"`
}

// VietQRData represents the data needed to generate a VietQR code
type VietQRData struct {
	BankID        string `json:"bank_id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Amount        int64  `json:"amount"`
	Description   string `json:"description"`
}

// NewVietQRData creates a new VietQRData instance
func NewVietQRData(bankID, accountNumber, accountName string, amount int64, description string) VietQRData {
	return VietQRData{
		BankID:        bankID,
		AccountNumber: accountNumber,
		AccountName:   accountName,
		Amount:        amount,
		Description:   description,
	}
}
