package auth

import (
	"net/http"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"
)

// InvalidCredentialsError is an error that is returned when the user credentials are invalid.
type invalidCredentialsError struct{}

func (e *invalidCredentialsError) ErrorStatusCode() int { return http.StatusTeapot }

func (*invalidCredentialsError) Error() string { return "Invalid credentials" }

func NewInvalidCredentialsError() serviceerrors.ServiceError {
	return &invalidCredentialsError{}
}

// InvalidOrExpiredTokenError is returned when the access token is invalid
type invalidOrExpiredTokenError struct{}

func (e *invalidOrExpiredTokenError) ErrorStatusCode() int { return http.StatusUnauthorized }

func (*invalidOrExpiredTokenError) Error() string { return "Invalid or expired token" }

func NewInvalidOrExpiredTokenError() serviceerrors.ServiceError {
	return &invalidOrExpiredTokenError{}
}
