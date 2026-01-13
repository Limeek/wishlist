package repository

import (
	"context"
	"wishlist-auth/internal/entity"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, refreshToken *entity.RefreshToken) error
	Delete(ctx context.Context, userId uuid.UUID) error
	FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error)
}

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
