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

func (s *Service) RefreshToken(refreshToken string) (*Tokens, error) {
	ok, claims, err := s.ValidateUserToken(refreshToken, "refresh")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, NewInvalidOrExpiredTokenError()
	}

	accessJti := uuid.New()
	refreshJti := uuid.New()

	exists, err := s.RDB.Exists(context.Background(), fmt.Sprintf("blacklisted:%s", (*claims)["jti"].(string))).Result()
	if err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, NewInvalidOrExpiredTokenError()
	}

	userIdStr, ok := (*claims)["sub"].(string)
	if !ok {
		return nil, NewInvalidOrExpiredTokenError()
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil, NewInvalidOrExpiredTokenError()
	}

	user, err := s.CRUDUsers.SelectUserByID(userId)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := s.GenerationAccessToken(user.ID, user.Role, accessJti, refreshJti)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.GenerationRefreshToken(user.ID, accessJti, refreshJti)
	if err != nil {
		return nil, err
	}

	newAccessTokenKey := fmt.Sprintf("access_token:%s:%s", user.ID.String(), accessJti.String())
	newRefreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", user.ID.String(), refreshJti.String())

	oldAccessTokenKey := fmt.Sprintf("access_token:%s:%s", user.ID.String(), (*claims)["access_jti"].(string))
	oldRefreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", user.ID.String(), (*claims)["jti"].(string))

	blacklistedAccessKey := fmt.Sprintf("blacklisted:%s", (*claims)["access_jti"].(string))
	blacklistedRefreshKey := fmt.Sprintf("blacklisted:%s", (*claims)["jti"].(string))

	rdbPipeline := s.RDB.Pipeline()
	rdbPipeline.Set(context.Background(), newAccessTokenKey, newAccessToken, 30*time.Minute)
	rdbPipeline.Set(context.Background(), newRefreshTokenKey, newRefreshToken, 7*24*time.Hour)
	rdbPipeline.Del(context.Background(), oldAccessTokenKey)
	rdbPipeline.Del(context.Background(), oldRefreshTokenKey)
	rdbPipeline.Set(context.Background(), blacklistedAccessKey, true, 60*time.Minute)
	rdbPipeline.Set(context.Background(), blacklistedRefreshKey, true, 8*24*time.Hour)
	_, err = rdbPipeline.Exec(context.Background())
	if err != nil {
		return nil, err
	}

	go func() {
		_ = s.CRUDUsers.UpdateLastLoginByID(user.ID)
	}()

	return &Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) LogoutUser(accessTokenClaims *jwt.MapClaims) error {
	accessJti := (*accessTokenClaims)["jti"].(string)
	refreshJti := (*accessTokenClaims)["refresh_jti"].(string)
	userId := (*accessTokenClaims)["sub"].(string)

	accessTokenKey := fmt.Sprintf("access_token:%s:%s", userId, accessJti)
	refreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", userId, refreshJti)
	blacklistedAccessKey := fmt.Sprintf("blacklisted:%s", accessJti)
	blacklistedRefreshKey := fmt.Sprintf("blacklisted:%s", refreshJti)

	rdbPipeline := s.RDB.Pipeline()
	rdbPipeline.Del(context.Background(), accessTokenKey)
	rdbPipeline.Del(context.Background(), refreshTokenKey)
	rdbPipeline.Set(context.Background(), blacklistedAccessKey, true, 60*time.Minute)
	rdbPipeline.Set(context.Background(), blacklistedRefreshKey, true, 8*24*time.Hour)
	_, err := rdbPipeline.Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateUserToken(token string, tokenType string) (bool, *jwt.MapClaims, error) {
	if s.Config.Secure.PublicKey != nil {
		token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return s.Config.Secure.PublicKey, nil
		})

		if err != nil || !token.Valid {
			return false, nil, NewInvalidOrExpiredTokenError()
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return false, nil, NewInvalidOrExpiredTokenError()
		}

		if !s.IsAvailableToken(tokenType, &claims) {
			return false, nil, NewInvalidOrExpiredTokenError()
		}

		return true, &claims, nil
	}
	return false, nil, nil
}

func (s *Service) IsAvailableToken(tokenType string, claims *jwt.MapClaims) bool {
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

	var tokenKey string
	switch tokenType {
	case "access":
		tokenKey = fmt.Sprintf("access_token:%s:%s", userId, jti)
	case "refresh":
		tokenKey = fmt.Sprintf("refresh_token:%s:%s", userId, jti)
	default:
		return false
	}
	_, err := s.RDB.Get(context.Background(), tokenKey).Result()
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
