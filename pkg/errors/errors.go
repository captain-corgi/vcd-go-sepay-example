package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of an error
type ErrorType string

const (
	// NotFound represents a resource not found error
	NotFound ErrorType = "NOT_FOUND"
	// ValidationFailed represents a validation error
	ValidationFailed ErrorType = "VALIDATION_FAILED"
	// Unauthorized represents an authentication error
	Unauthorized ErrorType = "UNAUTHORIZED"
	// Forbidden represents an authorization error
	Forbidden ErrorType = "FORBIDDEN"
	// Internal represents an internal server error
	Internal ErrorType = "INTERNAL"
	// BadRequest represents a bad request error
	BadRequest ErrorType = "BAD_REQUEST"
)

// AppError is a custom error type that includes additional context
type AppError struct {
	Type      ErrorType `json:"type"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	OrigError error     `json:"-"`
	Code      int       `json:"-"`
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.OrigError != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Message, e.Details, e.OrigError.Error())
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Details)
}

// Unwrap returns the original error (implements go1.13+ error unwrapping)
func (e *AppError) Unwrap() error {
	return e.OrigError
}

// HTTPStatusCode returns an appropriate HTTP status code based on the error type
func (e *AppError) HTTPStatusCode() int {
	if e.Code != 0 {
		return e.Code
	}

	switch e.Type {
	case NotFound:
		return http.StatusNotFound
	case ValidationFailed:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// NewNotFound creates a new not found error
func NewNotFound(message, details string, origErr error) *AppError {
	return &AppError{
		Type:      NotFound,
		Message:   message,
		Details:   details,
		OrigError: origErr,
		Code:      http.StatusNotFound,
	}
}

// NewValidationFailed creates a new validation error
func NewValidationFailed(message, details string, origErr error) *AppError {
	return &AppError{
		Type:      ValidationFailed,
		Message:   message,
		Details:   details,
		OrigError: origErr,
		Code:      http.StatusBadRequest,
	}
}

// NewUnauthorized creates a new unauthorized error
func NewUnauthorized(message, details string, origErr error) *AppError {
	return &AppError{
		Type:      Unauthorized,
		Message:   message,
		Details:   details,
		OrigError: origErr,
		Code:      http.StatusUnauthorized,
	}
}

// NewForbidden creates a new forbidden error
func NewForbidden(message, details string, origErr error) *AppError {
	return &AppError{
		Type:      Forbidden,
		Message:   message,
		Details:   details,
		OrigError: origErr,
		Code:      http.StatusForbidden,
	}
}

// NewInternal creates a new internal server error
func NewInternal(message, details string, origErr error) *AppError {
	return &AppError{
		Type:      Internal,
		Message:   message,
		Details:   details,
		OrigError: origErr,
		Code:      http.StatusInternalServerError,
	}
}

// NewBadRequest creates a new bad request error
func NewBadRequest(message, details string, origErr error) *AppError {
	return &AppError{
		Type:      BadRequest,
		Message:   message,
		Details:   details,
		OrigError: origErr,
		Code:      http.StatusBadRequest,
	}
}
