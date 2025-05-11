package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/database"
	"github.com/Census-Population-Project/API/internal/service/users"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Config *config.Config

	DB     *database.DataBase
	RDB    *redis.Client
	Logger *logrus.Logger

	CRUDUsers *users.CRUDUsers
}

func (s *Service) CompareHashAndPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Service) GenerationAccessToken(userId uuid.UUID, role string, accessJti, refreshJti uuid.UUID) (string, error) {
	var expTime = time.Now().Add(30 * time.Minute).Unix()

	claims := jwt.MapClaims{
		"sub":         userId,
		"type":        "access",
		"jti":         accessJti.String(),
		"refresh_jti": refreshJti.String(),
		"nbf":         time.Now().Unix(),
		"exp":         expTime,
		"role":        role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signedToken, err := token.SignedString(s.Config.Secure.PrivateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *Service) GenerationRefreshToken(userId uuid.UUID, accessJti, refreshJti uuid.UUID) (string, error) {
	var expTime = time.Now().Add(7 * 24 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"sub":        userId,
		"type":       "refresh",
		"jti":        refreshJti.String(),
		"access_jti": accessJti.String(),
		"nbf":        time.Now().Unix(),
		"exp":        expTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signedToken, err := token.SignedString(s.Config.Secure.PrivateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *Service) LoginUser(email, password string) (*Tokens, error) {
	user, err := s.CRUDUsers.SelectUserByEmail(email)
	if err != nil {
		return nil, err
	}

	userAuth, err := s.CRUDUsers.SelectUserAuthByEmail(email)
	if err != nil {
		return nil, err
	}

	if !s.CompareHashAndPassword(userAuth.Password, password) {
		return nil, NewInvalidCredentialsError()
	}

	accessJti := uuid.New()
	refreshJti := uuid.New()

	accessToken, err := s.GenerationAccessToken(user.ID, user.Role, accessJti, refreshJti)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerationRefreshToken(user.ID, accessJti, refreshJti)
	if err != nil {
		return nil, err
	}

	accessTokenKey := fmt.Sprintf("access_token:%s:%s", user.ID.String(), accessJti.String())
	refreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", user.ID.String(), refreshJti.String())

	rdbPipeline := s.RDB.Pipeline()
	rdbPipeline.Set(context.Background(), accessTokenKey, accessToken, 30*time.Minute)
	rdbPipeline.Set(context.Background(), refreshTokenKey, refreshToken, 7*24*time.Hour)
	_, err = rdbPipeline.Exec(context.Background())
	if err != nil {
		return nil, err
	}

	go func() {
		_ = s.CRUDUsers.UpdateLastLoginByID(user.ID)
	}()

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*Tokens, error) { // TODO: Implement
	return nil, nil
}

func (s *Service) LogoutUser(refreshToken string) error { // TODO: Implement
	return nil
}

func (s *Service) ValidateUserToken(accessToken string) (bool, *jwt.MapClaims, error) {
	if s.Config.Secure.PublicKey != nil {
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			return s.Config.Secure.PublicKey, nil
		})

		if err != nil || !token.Valid {
			return false, nil, NewInvalidOrExpiredTokenError()
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return false, nil, NewInvalidOrExpiredTokenError()
		}

		if !s.IsAvailableToken(&claims) {
			return false, nil, NewInvalidOrExpiredTokenError()
		}

		return true, &claims, nil
	}
	return false, nil, nil
}

func (s *Service) IsAvailableToken(claims *jwt.MapClaims) bool {
	jtiVal, ok := (*claims)["jti"]
	if !ok {
		return false
	}
	jti, ok := jtiVal.(string)
	if !ok {
		return false
	}

	userIdVal, ok := (*claims)["sub"]
	if !ok {
		return false
	}
	userId, ok := userIdVal.(string)
	if !ok {
		return false
	}

	_, err := s.RDB.Get(context.Background(), fmt.Sprintf("access_token:%s:%s", userId, jti)).Result()
	if err != nil {
		return false
	}
	return true
}

func NewService(cfg *config.Config, db *database.DataBase, rdb *redis.Client, logger *logrus.Logger) *Service {
	return &Service{
		Config: cfg,

		DB:     db,
		RDB:    rdb,
		Logger: logger,

		CRUDUsers: users.NewUsersCRUD(db, logger),
	}
}
