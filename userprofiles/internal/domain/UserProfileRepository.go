package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserProfileRepository interface {
	Get(ctx context.Context, userId uuid.UUID) (*UserProfile, error)
	Create(ctx context.Context, userProfile UserProfile) error
}
