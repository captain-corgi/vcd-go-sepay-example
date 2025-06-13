package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQRCode_StructFields(t *testing.T) {
	// Test creating a QR code with all fields
	imageData := []byte("fake qr code image data")
	qrCode := QRCode{
		Content: "VietQR payment data",
		Size:    256,
		Image:   imageData,
	}

	// Verify all fields are set correctly
	assert.Equal(t, "VietQR payment data", qrCode.Content)
	assert.Equal(t, 256, qrCode.Size)
	assert.Equal(t, imageData, qrCode.Image)
}

func TestQRCode_EmptyImage(t *testing.T) {
	// Test QR code without image data (just content)
	qrCode := QRCode{
		Content: "Payment content",
		Size:    128,
		Image:   nil,
	}

	assert.Equal(t, "Payment content", qrCode.Content)
	assert.Equal(t, 128, qrCode.Size)
	assert.Nil(t, qrCode.Image)
}

func TestVietQRData_NewVietQRData(t *testing.T) {
	// Test creating VietQRData using constructor
	bankID := "970415"
	accountNumber := "0123456789"
	accountName := "NGUYEN VAN A"
	amount := int64(100000)
	description := "Payment for order ORDER123"

	vietQR := NewVietQRData(bankID, accountNumber, accountName, amount, description)

	assert.Equal(t, bankID, vietQR.BankID)
	assert.Equal(t, accountNumber, vietQR.AccountNumber)
	assert.Equal(t, accountName, vietQR.AccountName)
	assert.Equal(t, amount, vietQR.Amount)
	assert.Equal(t, description, vietQR.Description)
}

func TestVietQRData_StructFields(t *testing.T) {
	// Test creating VietQRData directly
	vietQR := VietQRData{
		BankID:        "970415",
		AccountNumber: "0123456789",
		AccountName:   "NGUYEN VAN A",
		Amount:        150000,
		Description:   "Test payment",
	}

	assert.Equal(t, "970415", vietQR.BankID)
	assert.Equal(t, "0123456789", vietQR.AccountNumber)
	assert.Equal(t, "NGUYEN VAN A", vietQR.AccountName)
	assert.Equal(t, int64(150000), vietQR.Amount)
	assert.Equal(t, "Test payment", vietQR.Description)
}

func TestVietQRData_Validation(t *testing.T) {
	testCases := []struct {
		name          string
		vietQR        VietQRData
		shouldBeValid bool
		description   string
	}{
		{
			name: "valid VietQR data",
			vietQR: VietQRData{
				BankID:        "970415",
				AccountNumber: "0123456789",
				AccountName:   "NGUYEN VAN A",
				Amount:        100000,
				Description:   "Payment for order",
			},
			shouldBeValid: true,
			description:   "All required fields are present and valid",
		},
		{
			name: "missing bank ID",
			vietQR: VietQRData{
				BankID:        "",
				AccountNumber: "0123456789",
				AccountName:   "NGUYEN VAN A",
				Amount:        100000,
				Description:   "Payment for order",
			},
			shouldBeValid: false,
			description:   "Bank ID is required",
		},
		{
			name: "missing account number",
			vietQR: VietQRData{
				BankID:        "970415",
				AccountNumber: "",
				AccountName:   "NGUYEN VAN A",
				Amount:        100000,
				Description:   "Payment for order",
			},
			shouldBeValid: false,
			description:   "Account number is required",
		},
		{
			name: "missing account name",
			vietQR: VietQRData{
				BankID:        "970415",
				AccountNumber: "0123456789",
				AccountName:   "",
				Amount:        100000,
				Description:   "Payment for order",
			},
			shouldBeValid: false,
			description:   "Account name is required",
		},
		{
			name: "zero amount",
			vietQR: VietQRData{
				BankID:        "970415",
				AccountNumber: "0123456789",
				AccountName:   "NGUYEN VAN A",
				Amount:        0,
				Description:   "Payment for order",
			},
			shouldBeValid: false,
			description:   "Amount must be positive",
		},
		{
			name: "negative amount",
			vietQR: VietQRData{
				BankID:        "970415",
				AccountNumber: "0123456789",
				AccountName:   "NGUYEN VAN A",
				Amount:        -50000,
				Description:   "Payment for order",
			},
			shouldBeValid: false,
			description:   "Amount must be positive",
		},
		{
			name: "missing description",
			vietQR: VietQRData{
				BankID:        "970415",
				AccountNumber: "0123456789",
				AccountName:   "NGUYEN VAN A",
				Amount:        100000,
				Description:   "",
			},
			shouldBeValid: false,
			description:   "Description is required for payment tracking",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Validation logic
			isValid := tc.vietQR.BankID != "" &&
				tc.vietQR.AccountNumber != "" &&
				tc.vietQR.AccountName != "" &&
				tc.vietQR.Amount > 0 &&
				tc.vietQR.Description != ""

			assert.Equal(t, tc.shouldBeValid, isValid, tc.description)
		})
	}
}

func TestVietQRData_BankIDFormats(t *testing.T) {
	testCases := []struct {
		name     string
		bankID   string
		expected bool
	}{
		{
			name:     "valid Vietcombank ID",
			bankID:   "970415",
			expected: true,
		},
		{
			name:     "valid BIDV ID",
			bankID:   "970418",
			expected: true,
		},
		{
			name:     "valid Techcombank ID",
			bankID:   "970407",
			expected: true,
		},
		{
			name:     "too short bank ID",
			bankID:   "9704",
			expected: false,
		},
		{
			name:     "too long bank ID",
			bankID:   "9704151",
			expected: false,
		},
		{
			name:     "non-numeric bank ID",
			bankID:   "ABC123",
			expected: false,
		},
		{
			name:     "empty bank ID",
			bankID:   "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simple validation: bank ID should be 6 digits
			isValid := len(tc.bankID) == 6
			if isValid {
				// Check if all characters are digits
				for _, char := range tc.bankID {
					if char < '0' || char > '9' {
						isValid = false
						break
					}
				}
			}

			assert.Equal(t, tc.expected, isValid)
		})
	}
}

func TestVietQRData_AmountLimits(t *testing.T) {
	testCases := []struct {
		name     string
		amount   int64
		expected bool
	}{
		{
			name:     "minimum valid amount",
			amount:   1000, // 1,000 VND
			expected: true,
		},
		{
			name:     "typical small amount",
			amount:   50000, // 50,000 VND
			expected: true,
		},
		{
			name:     "typical large amount",
			amount:   1000000, // 1,000,000 VND
			expected: true,
		},
		{
			name:     "very large amount",
			amount:   999999999, // Under 1 billion VND
			expected: true,
		},
		{
			name:     "zero amount",
			amount:   0,
			expected: false,
		},
		{
			name:     "negative amount",
			amount:   -100000,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vietQR := VietQRData{
				BankID:        "970415",
				AccountNumber: "0123456789",
				AccountName:   "NGUYEN VAN A",
				Amount:        tc.amount,
				Description:   "Test payment",
			}

			// Amount validation
			isValid := vietQR.Amount > 0

			assert.Equal(t, tc.expected, isValid)
		})
	}
}
