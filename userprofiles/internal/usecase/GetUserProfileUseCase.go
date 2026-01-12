package usecase

import (
	"context"
	"userprofiles/internal/domain"

	"github.com/google/uuid"
)

type GetUserProfileUseCase struct {
	profilesRepo domain.UserProfileRepository
}

func NewGetUserProfileUseCase(profilesRepo domain.UserProfileRepository) *GetUserProfileUseCase {
	return &GetUserProfileUseCase{profilesRepo}
}

func (u *GetUserProfileUseCase) GetUserProfile(ctx context.Context, userId uuid.UUID) (*domain.UserProfile, error) {
	userProfile, err := u.profilesRepo.Get(ctx, userId)

	if err != nil {
		return nil, err
	}

	return userProfile, nil
}
