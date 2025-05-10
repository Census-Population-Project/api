package tools

import (
	"net/http"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"
)

// InvalidTimeFormatError is an error type for invalid time format
type invalidTimeFormatError struct{}

func (e *invalidTimeFormatError) ErrorStatusCode() int { return http.StatusBadRequest }

func (*invalidTimeFormatError) Error() string { return "Invalid time format" }

func NewInvalidTimeFormatError() serviceerrors.ServiceError {
	return &invalidTimeFormatError{}
}

// InvalidValueForFieldError is an error type for invalid value for field
type invalidValueForFieldError struct {
	Field string
}

func (e *invalidValueForFieldError) ErrorStatusCode() int { return http.StatusBadRequest }

func (e *invalidValueForFieldError) Error() string {
	return "Invalid value for field: " + e.Field
}

func NewInvalidValueForFieldError(field string) serviceerrors.ServiceError {
	return &invalidValueForFieldError{Field: field}
}

// UnknownFieldError is an error type for unknown field
type unknownFieldError struct {
	Field string
}

func (e *unknownFieldError) ErrorStatusCode() int { return http.StatusBadRequest }

func (e *unknownFieldError) Error() string {
	return "Unknown field: " + e.Field
}

func NewUnknownFieldError(field string) serviceerrors.ServiceError {
	return &unknownFieldError{Field: field}
}
