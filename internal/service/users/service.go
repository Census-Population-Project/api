package users

import (
	"github.com/Census-Population-Project/API/internal/service/api/tools"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	InitDefaultUser() error

	CreateUser(email, password, firstName, lastName, role string, defaultUser bool) (*User, error)

	GetUsers(limit, offset int) ([]User, error)
	GetUserByID(id uuid.UUID) (*User, error)

	UpdateUserByID(
		id uuid.UUID,
		email tools.Optional[string], password tools.Optional[string],
		firstName tools.Optional[string], lastName tools.Optional[string],
		role tools.Optional[string],
		isPatch bool,
	) (*User, error)
}
