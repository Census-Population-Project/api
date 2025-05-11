package users

import (
	"net/http"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"
)

// UserCreationError is an error that is returned when user creation fails.
type userCreationError struct {
	message string
}

func (e *userCreationError) ErrorStatusCode() int { return http.StatusInternalServerError }

func (e *userCreationError) Error() string { return e.message }

func NewUserCreationError(message string) serviceerrors.ServiceError {
	return &userCreationError{message: message}
}

// UserNotFoundError is an error that is returned when a user is not found.
type userNotFoundError struct{}

func (e *userNotFoundError) ErrorStatusCode() int { return http.StatusNotFound }

func (*userNotFoundError) Error() string { return "User not found" }

func NewUserNotFoundError() serviceerrors.ServiceError {
	return &userNotFoundError{}
}

// UserAlreadyExistsError is an error that is returned when a user with the same email already exists.
type userAlreadyExistsError struct{}

func (e *userAlreadyExistsError) ErrorStatusCode() int { return http.StatusConflict }

func (*userAlreadyExistsError) Error() string { return "User with this email already exists" }

func NewUserAlreadyExistsError() serviceerrors.ServiceError {
	return &userAlreadyExistsError{}
}

// UserRoleNotFoundError is an error that is returned when a user role is not found.
type userRoleNotFoundError struct{}

func (e *userRoleNotFoundError) ErrorStatusCode() int { return http.StatusBadRequest }

func (*userRoleNotFoundError) Error() string { return "User role not found" }

func NewUserRoleNotFoundError() serviceerrors.ServiceError {
	return &userRoleNotFoundError{}
}
