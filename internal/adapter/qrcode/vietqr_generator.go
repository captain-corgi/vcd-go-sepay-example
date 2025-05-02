package qrcode

import (
	"fmt"

	"github.com/captain-corgi/vcd-go-sepay-example/internal/domain/entity"
	goqrcode "github.com/skip2/go-qrcode"
)

// VietQRGenerator generates QR codes according to VietQR standard
type VietQRGenerator struct {
	defaultSize int
}

// NewVietQRGenerator creates a new VietQR generator
func NewVietQRGenerator(defaultSize int) *VietQRGenerator {
	return &VietQRGenerator{
		defaultSize: defaultSize,
	}
}

// Generate creates a QR code based on the provided VietQR data
func (g *VietQRGenerator) Generate(data entity.VietQRData) (*entity.QRCode, error) {
	// Validate input data
	if err := g.validateData(data); err != nil {
		return nil, err
	}

	// Format according to VietQR standard: bankid|account|amount|description
	content := fmt.Sprintf("VietQR|%s|%s|%d|%s",
		data.BankID,
		data.AccountNumber,
		data.Amount,
		data.Description,
	)

	// Generate QR code image
	qrBytes, err := goqrcode.Encode(content, goqrcode.Medium, g.defaultSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Return QR code data
	return &entity.QRCode{
		Content: content,
		Size:    g.defaultSize,
		Image:   qrBytes,
	}, nil
}

// SaveToFile saves a QR code to a file
func (g *VietQRGenerator) SaveToFile(data entity.VietQRData, filepath string) error {
	// Validate input data
	if err := g.validateData(data); err != nil {
		return err
	}

	// Format according to VietQR standard
	content := fmt.Sprintf("VietQR|%s|%s|%d|%s",
		data.BankID,
		data.AccountNumber,
		data.Amount,
		data.Description,
	)

	// Write QR code to file
	err := goqrcode.WriteFile(content, goqrcode.Medium, g.defaultSize, filepath)
	if err != nil {
		return fmt.Errorf("failed to save QR code to file: %w", err)
	}

	return nil
}

// validateData checks if the required fields are present
func (g *VietQRGenerator) validateData(data entity.VietQRData) error {
	if data.BankID == "" {
		return fmt.Errorf("bank ID is required")
	}
	if data.AccountNumber == "" {
		return fmt.Errorf("account number is required")
	}
	return nil
}
