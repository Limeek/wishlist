package repository

import (
	"context"
	"userprofiles/internal/entity"

	"github.com/google/uuid"
)

type UserProfileRepository interface {
	Get(ctx context.Context, userId uuid.UUID) (*entity.UserProfile, error)
	Create(ctx context.Context, userProfile entity.UserProfile) error
}
