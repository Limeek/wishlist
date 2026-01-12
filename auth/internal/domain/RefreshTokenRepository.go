package domain

import (
	"context"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, refreshToken *RefreshToken) error
	Delete(ctx context.Context, userId uuid.UUID) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
}
