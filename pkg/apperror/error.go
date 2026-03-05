package apperror

import (
	"errors"
	"net/http"
)

type AppError struct {
	Err     error
	Message string
	Code    int
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(message string, code int) *AppError {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

func Wrap(err error, message string, code int) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// Common errors
var (
	ErrNotFound       = New("resource not found", http.StatusNotFound)
	ErrInternalServer = New("internal server error", http.StatusInternalServerError)
	ErrBadRequest     = New("bad request", http.StatusBadRequest)
	ErrUnauthorized   = New("unauthorized", http.StatusUnauthorized)
	ErrForbidden      = New("forbidden", http.StatusForbidden)
	ErrValidation     = New("validation failed", http.StatusUnprocessableEntity)
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, &target)
}

// Convert converts any error to AppError
func Convert(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return Wrap(err, "internal server error", http.StatusInternalServerError)
}
