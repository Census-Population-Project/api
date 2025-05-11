package users

import "github.com/Census-Population-Project/API/internal/service/api/tools"

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

type UpdateUserRequest struct {
	Email     tools.Optional[string] `json:"email"`
	Password  tools.Optional[string] `json:"password"`
	FirstName tools.Optional[string] `json:"first_name"`
	LastName  tools.Optional[string] `json:"last_name"`
	Role      tools.Optional[string] `json:"role"`
}
