package middleware

import (
	"net/http"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"
)

// AuthorizationHeaderIsMissingError is an error that is returned when the Authorization header is missing
type authorizationHeaderIsMissingError struct{}

func (e *authorizationHeaderIsMissingError) ErrorStatusCode() int { return http.StatusUnauthorized }

func (*authorizationHeaderIsMissingError) Error() string { return "Authorization header is missing" }

func NewAuthorizationHeaderIsMissingError() serviceerrors.ServiceError {
	return &authorizationHeaderIsMissingError{}
}
