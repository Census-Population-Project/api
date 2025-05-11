package errors

import "net/http"

type ServiceError interface {
	error
	ErrorStatusCode() int
}

// ForbiddenError is an error that is returned when a user service is not found
type forbiddenError struct{}

func (e *forbiddenError) ErrorStatusCode() int { return http.StatusForbidden }

func (*forbiddenError) Error() string { return "Forbidden" }

func NewForbiddenError() ServiceError {
	return &forbiddenError{}
}
