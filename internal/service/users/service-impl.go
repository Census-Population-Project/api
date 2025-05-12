package users

import (
	"errors"
	"slices"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/database"
	"github.com/Census-Population-Project/API/internal/service/api/tools"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Config *config.Config

	DB     *database.DataBase
	Logger *logrus.Logger

	CRUDUsers *CRUDUsers
}

func (s *Service) InitDefaultUser() error {
	user, err := s.CRUDUsers.SelectDefaultUser()
	if err != nil {
		var srvErr serviceerrors.ServiceError
		if !errors.As(err, &srvErr) {
			if !errors.Is(err, NewUserAlreadyExistsError()) {
				return err
			}
		}
	}

	if user != nil {
		return nil
	}

	if s.Config.DefaultUserEmail == "" {
		s.Config.DefaultUserEmail = "admin@example.com"
		s.Logger.Info("Default user email is empty, using default: ", s.Config.DefaultUserEmail)
	}
	if s.Config.DefaultUserPassword == "" {
		s.Config.DefaultUserPassword = "secure"
		s.Logger.Info("Default user password is empty, using default: ", s.Config.DefaultUserPassword)
	}

	user, err = s.CreateUser(s.Config.DefaultUserEmail, s.Config.DefaultUserPassword, "", "", "administrator", true)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateUser(email, password, firstName, lastName, role string, defaultUser bool) (*User, error) {
	if !slices.Contains(Roles, role) {
		return nil, NewUserRoleNotFoundError()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.CRUDUsers.CreateUser(email, string(hashedPassword), firstName, lastName, role, defaultUser)
	if err != nil {
		return nil, err
	}
	user.DefaultUser = nil

	return user, nil
}

func (s *Service) GetUsers(limit, offset int) ([]User, error) {
	return s.CRUDUsers.SelectUsers(limit, offset)
}

func (s *Service) GetUserByID(id uuid.UUID) (*User, error) {
	user, err := s.CRUDUsers.SelectUserByID(id)
	if err != nil {
		return nil, err
	}
	user.DefaultUser = nil

	return user, nil
}

func (s *Service) UpdateUserByID(
	id uuid.UUID,
	email tools.Optional[string], password tools.Optional[string],
	firstName tools.Optional[string], lastName tools.Optional[string],
	role tools.Optional[string],
	isPatch bool,
) (*User, error) {
	user, err := s.CRUDUsers.SelectUserByID(id)
	if err != nil {
		return nil, err
	}

	userAuth, err := s.CRUDUsers.SelectUserAuthByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	updateEmail := tools.UpdateOptionalField(email, &user.Email, isPatch, true)
	updatePassword := tools.UpdateOptionalField(password, &userAuth.Password, isPatch, true)
	updateFirstName := tools.UpdateOptionalField(firstName, &user.FirstName, isPatch, true)
	updateLastName := tools.UpdateOptionalField(lastName, &user.LastName, isPatch, true)
	updateRole := tools.UpdateOptionalField(role, &user.Role, isPatch, true)

	user, err = s.CRUDUsers.UpdateUserByID(
		id,
		*updateEmail, *updatePassword,
		*updateFirstName, *updateLastName,
		*updateRole,
	)
	if err != nil {
		return nil, err
	}
	user.DefaultUser = nil

	return user, nil
}

func NewService(cfg *config.Config, db *database.DataBase, logger *logrus.Logger) *Service {
	return &Service{
		Config: cfg,

		DB:     db,
		Logger: logger,

		CRUDUsers: NewUsersCRUD(db, logger),
	}
}
