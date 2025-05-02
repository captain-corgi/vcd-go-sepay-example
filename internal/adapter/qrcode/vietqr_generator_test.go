package qrcode_test

import (
	"testing"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/adapter/qrcode"
	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestVietQRGenerator_Generate(t *testing.T) {
	tests := []struct {
		name           string
		qrData         entity.VietQRData
		expectedFormat string
		expectError    bool
	}{
		{
			name: "valid data with amount",
			qrData: entity.VietQRData{
				BankID:        "970436",
				AccountNumber: "1234567890",
				AccountName:   "John Doe",
				Amount:        1000000,
				Description:   "ORDER123",
			},
			expectedFormat: "VietQR|970436|1234567890|1000000|ORDER123",
			expectError:    false,
		},
		{
			name: "valid data without amount",
			qrData: entity.VietQRData{
				BankID:        "970436",
				AccountNumber: "1234567890",
				AccountName:   "John Doe",
				Amount:        0,
				Description:   "ORDER123",
			},
			expectedFormat: "VietQR|970436|1234567890|0|ORDER123",
			expectError:    false,
		},
		{
			name: "invalid data - missing bank ID",
			qrData: entity.VietQRData{
				BankID:        "",
				AccountNumber: "1234567890",
				AccountName:   "John Doe",
				Amount:        1000000,
				Description:   "ORDER123",
			},
			expectedFormat: "",
			expectError:    true,
		},
		{
			name: "invalid data - missing account number",
			qrData: entity.VietQRData{
				BankID:        "970436",
				AccountNumber: "",
				AccountName:   "John Doe",
				Amount:        1000000,
				Description:   "ORDER123",
			},
			expectedFormat: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create generator with default size
			generator := qrcode.NewVietQRGenerator(256)
			
			// Generate QR code
			qrCode, err := generator.Generate(tt.qrData)
			
			// Check if we expected an error
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, qrCode)
				return
			}
			
			// If no error expected, validate the QR code
			assert.NoError(t, err)
			assert.NotNil(t, qrCode)
			assert.Equal(t, tt.expectedFormat, qrCode.Content)
			assert.Equal(t, 256, qrCode.Size)
			assert.NotEmpty(t, qrCode.Image)
		})
	}
}

func TestVietQRGenerator_SaveToFile(t *testing.T) {
	// Skip in CI environments or when file writing is not possible
	t.Skip("Skipping test that writes to the file system")
	
	generator := qrcode.NewVietQRGenerator(256)
	
	qrData := entity.VietQRData{
		BankID:        "970436",
		AccountNumber: "1234567890",
		AccountName:   "John Doe",
		Amount:        1000000,
		Description:   "ORDER123",
	}
	
	// Test saving to a temporary file
	tempFile := "/tmp/test_vietqr.png"
	err := generator.SaveToFile(qrData, tempFile)
	assert.NoError(t, err)
	
	// Additional checks could verify file exists and has content
}