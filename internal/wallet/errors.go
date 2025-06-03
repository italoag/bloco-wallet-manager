package wallet

import (
	"context"
	"fmt"
)

// ErrorCode represents specific wallet error types
type ErrorCode string

const (
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound       ErrorCode = "WALLET_NOT_FOUND"
	ErrCodeAlreadyExists  ErrorCode = "WALLET_EXISTS"
	ErrCodeAuthentication ErrorCode = "AUTH_FAILED"
	ErrCodeNetwork        ErrorCode = "NETWORK_ERROR"
	ErrCodeCrypto         ErrorCode = "CRYPTO_ERROR"
	ErrCodeStorage        ErrorCode = "STORAGE_ERROR"
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
)

// WalletError represents a domain-specific error
type WalletError struct {
	Code          ErrorCode `json:"code"`
	Message       string    `json:"message"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	Cause         error     `json:"-"`
}

func (e *WalletError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *WalletError) Unwrap() error {
	return e.Cause
}

// Error constructor functions
func NewValidationError(message string) *WalletError {
	return &WalletError{
		Code:    ErrCodeValidation,
		Message: message,
	}
}

func NewValidationErrorWithCause(message string, cause error) *WalletError {
	return &WalletError{
		Code:    ErrCodeValidation,
		Message: message,
		Cause:   cause,
	}
}

func NewNotFoundError(message string) *WalletError {
	return &WalletError{
		Code:    ErrCodeNotFound,
		Message: message,
	}
}

func NewAlreadyExistsError(message string) *WalletError {
	return &WalletError{
		Code:    ErrCodeAlreadyExists,
		Message: message,
	}
}

func NewAuthenticationError(message string) *WalletError {
	return &WalletError{
		Code:    ErrCodeAuthentication,
		Message: message,
	}
}

func NewAuthenticationErrorWithCause(message string, cause error) *WalletError {
	return &WalletError{
		Code:    ErrCodeAuthentication,
		Message: message,
		Cause:   cause,
	}
}

func NewNetworkError(message string, cause error) *WalletError {
	return &WalletError{
		Code:    ErrCodeNetwork,
		Message: message,
		Cause:   cause,
	}
}

func NewCryptoError(message string, cause error) *WalletError {
	return &WalletError{
		Code:    ErrCodeCrypto,
		Message: message,
		Cause:   cause,
	}
}

func NewStorageError(message string, cause error) *WalletError {
	return &WalletError{
		Code:    ErrCodeStorage,
		Message: message,
		Cause:   cause,
	}
}

func NewInternalError(message string, cause error) *WalletError {
	return &WalletError{
		Code:    ErrCodeInternal,
		Message: message,
		Cause:   cause,
	}
}

// Helper to add correlation ID from context
func AddCorrelationID(ctx context.Context, err *WalletError) *WalletError {
	if correlationID, ok := ctx.Value("correlation_id").(string); ok {
		err.CorrelationID = correlationID
	}
	return err
}