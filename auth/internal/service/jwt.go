package service

import (
	"os"
	"time"
	"wishlist-auth/internal/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtTokenTTL = time.Hour * 24
var jwtRefreshTokenTTL = time.Hour * 24 * 7

type JwtService struct {
}

func NewJwtService() *JwtService {
	return &JwtService{}
}

func (service *JwtService) GenerateToken(userId uuid.UUID) (string, error) {
	tokenClaims := jwt.MapClaims{
		"sub": userId.String(),
		"exp": time.Now().Add(jwtTokenTTL).Unix(),
		"iat": time.Now().Unix(),
		"iss": "wishlist-auth",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (service *JwtService) GenerateRefreshToken(userId uuid.UUID) (*entity.RefreshToken, error) {
	expirationTime := time.Now().Add(jwtRefreshTokenTTL)

	tokenClaims := jwt.MapClaims{
		"sub": userId.String(),
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
		"iss": "wishlist-auth",
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))

	if err != nil {
		return nil, err
	}

	return &entity.RefreshToken{
		Token:     token,
		UserId:    userId,
		ExpiresAt: expirationTime,
	}, nil
}
