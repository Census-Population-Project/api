package users

import (
	"context"
	"errors"

	"github.com/Census-Population-Project/API/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
)

type CRUDInterface interface {
	CreateUser(email, passwordHash, firstName, lastName, role string, defaultUser bool) (*User, error)

	SelectUsers(limit, offset int) ([]User, error)
	SelectUserByEmail(email string) (*User, error)
	SelectUserByID(id uuid.UUID) (*User, error)
	SelectUserAuthByEmail(email string) (*UserAuth, error)
	SelectDefaultUser() (*User, error)

	UpdateUserByID(
		id uuid.UUID,
		email string, password string,
		firstName string, lastName string,
		role string,
	) (*User, error)
	UpdateLastLoginByID(id uuid.UUID) error
}

type CRUDUsers struct {
	DataBase *database.DataBase
	Logger   *logrus.Logger
}

func (s *CRUDUsers) CreateUser(
	email, password, firstName, lastName, role string,
	defaultUser bool,
) (*User, error) {
	tx, err := s.DataBase.DBPool.Begin(context.Background())
	if err != nil {
		s.Logger.Error("Failed to begin transaction: ", err)
		return nil, err
	}
	defer tx.Rollback(context.Background())

	query := `INSERT INTO auth.users (email, first_name, last_name, role, default_user)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, email, first_name, last_name, role, default_user`
	row, err := tx.Query(context.Background(), query, email, firstName, lastName, role, defaultUser)
	if err != nil {
		s.Logger.Error("Failed to create user: ", err)
		return nil, err
	}
	defer row.Close()

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[User])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, NewUserAlreadyExistsError()
			}
		}
		s.Logger.Error("Failed to create user: ", err)
		return nil, err
	}

	authQuery := `INSERT INTO auth.user_auth (user_id, password)
	VALUES ($1, $2)`

	_, err = tx.Exec(context.Background(), authQuery, user.ID, password)
	if err != nil {
		s.Logger.Error("Failed to create user auth: ", err)
		return nil, err
	}

	if err = tx.Commit(context.Background()); err != nil {
		s.Logger.Error("Failed to commit transaction: ", err)
		return nil, err
	}

	return &user, nil
}

func (s *CRUDUsers) SelectUsers(limit, offset int) ([]User, error) {
	query := `SELECT id, email, first_name, last_name, role, NULL AS default_user
	FROM auth.users
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2`
	rows, err := s.DataBase.DBPool.Query(context.Background(), query, limit, offset)
	if err != nil {
		s.Logger.Error("Failed to get users: ", err)
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		s.Logger.Error("Failed to collect users: ", err)
		return nil, err
	}

	return users, nil
}

func (s *CRUDUsers) SelectUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, first_name, last_name, role, default_user
	FROM auth.users
	WHERE email = $1`
	row, err := s.DataBase.DBPool.Query(context.Background(), query, email)
	if err != nil {
		s.Logger.Error("Failed to get user by email: ", err)
		return nil, err
	}
	defer row.Close()

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewUserNotFoundError()
		}
		s.Logger.Error("Failed to collect user by email: ", err)
		return nil, err
	}

	return &user, nil
}

func (s *CRUDUsers) SelectUserByID(id uuid.UUID) (*User, error) {
	query := `SELECT id, email, first_name, last_name, role, default_user
	FROM auth.users
	WHERE id = $1`
	row, err := s.DataBase.DBPool.Query(context.Background(), query, id)
	if err != nil {
		s.Logger.Error("Failed to get user by ID: ", err)
		return nil, err
	}
	defer row.Close()

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewUserNotFoundError()
		}
		s.Logger.Error("Failed to collect user by ID: ", err)
		return nil, err
	}

	return &user, nil
}

func (s *CRUDUsers) SelectUserAuthByEmail(email string) (*UserAuth, error) {
	query := `SELECT ua.user_id, ua.password, ua.last_login, ua.last_password_change
	FROM auth.users u
	JOIN auth.user_auth ua ON u.id = ua.user_id
	WHERE email = $1`
	row, err := s.DataBase.DBPool.Query(context.Background(), query, email)
	if err != nil {
		s.Logger.Error("Failed to get user auth by email: ", err)
		return nil, NewUserNotFoundError()
	}
	defer row.Close()

	userAuth, err := pgx.CollectOneRow(row, pgx.RowToStructByName[UserAuth])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewUserNotFoundError()
		}
		s.Logger.Error("Failed to collect user auth by email: ", err)
		return nil, err
	}

	return &userAuth, nil
}

func (s *CRUDUsers) SelectDefaultUser() (*User, error) {
	query := `SELECT id, email, first_name, last_name, role, default_user
	FROM auth.users
	WHERE default_user = true`
	row, err := s.DataBase.DBPool.Query(context.Background(), query)
	if err != nil {
		s.Logger.Error("Failed to get default user: ", err)
		return nil, err
	}
	defer row.Close()

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewUserNotFoundError()
		}
		s.Logger.Error("Failed to collect default user: ", err)
		return nil, err
	}

	return &user, nil
}

func (s *CRUDUsers) UpdateUserByID(
	id uuid.UUID,
	email string, password string,
	firstName string, lastName string,
	role string,
) (*User, error) {
	tx, err := s.DataBase.DBPool.Begin(context.Background())
	if err != nil {
		s.Logger.Error("Failed to begin transaction: ", err)
		return nil, err
	}
	defer tx.Rollback(context.Background())

	query := `UPDATE auth.users
	SET email = $1, first_name = $2, last_name = $3, role = $4
	WHERE id = $5
	RETURNING id, email, first_name, last_name, role, default_user`
	row, err := tx.Query(context.Background(), query, email, firstName, lastName, role, id)
	if err != nil {
		s.Logger.Error("Failed to update user: ", err)
		return nil, err
	}
	defer row.Close()

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewUserNotFoundError()
		}
		s.Logger.Error("Failed to update user: ", err)
		return nil, err
	}

	authQuery := `UPDATE auth.user_auth
	SET password = $1
	WHERE user_id = $2`
	_, err = tx.Exec(context.Background(), authQuery, password, id)
	if err != nil {
		s.Logger.Error("Failed to update user auth: ", err)
		return nil, err
	}

	if err = tx.Commit(context.Background()); err != nil {
		s.Logger.Error("Failed to commit transaction: ", err)
		return nil, err
	}

	return &user, nil
}

func (s *CRUDUsers) UpdateLastLoginByID(id uuid.UUID) error {
	query := `UPDATE auth.user_auth
	SET last_login = NOW()
	WHERE user_id = $1`
	_, err := s.DataBase.DBPool.Exec(context.Background(), query, id)
	if err != nil {
		s.Logger.Error("Failed to update last login: ", err)
		return err
	}

	return nil
}

func NewUsersCRUD(db *database.DataBase, log *logrus.Logger) *CRUDUsers {
	return &CRUDUsers{
		DataBase: db,
		Logger:   log,
	}
}
