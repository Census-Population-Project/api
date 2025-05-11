package users

import (
	"time"

	"github.com/google/uuid"
)

var Roles = []string{"agent", "administrator"}

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Email       string    `json:"email" db:"email"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Role        string    `json:"role" db:"role"`
	DefaultUser *bool     `json:"default_user,omitempty" db:"default_user"`
}

type UserAuth struct {
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	Password           string    `json:"password" db:"password"`
	LastLogin          time.Time `json:"last_login" db:"last_login"`
	LastPasswordChange time.Time `json:"last_password_change" db:"last_password_change"`
}
