package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	CompareHashAndPassword(hash, password string) bool
	GenerationAccessToken(userId uuid.UUID, role string, accessJti, refreshJti uuid.UUID) (string, error)
	GenerationRefreshToken(userId uuid.UUID, accessJti, refreshJti uuid.UUID) (string, error)

	LoginUser(email, password string) (*Tokens, error)
	RefreshToken(refreshToken string) (*Tokens, error)
	LogoutUser(refreshToken string) error

	ValidateUserToken(accessToken string) (bool, *jwt.MapClaims, error)
	IsAvailableToken(claims *jwt.MapClaims) bool
}
