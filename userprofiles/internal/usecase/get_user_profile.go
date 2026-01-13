package usecase

import (
	"context"
	"userprofiles/internal/entity"
	"userprofiles/internal/repository"

	"github.com/google/uuid"
)

type GetUserProfileUseCase struct {
	profilesRepo repository.UserProfileRepository
}

func NewGetUserProfileUseCase(profilesRepo repository.UserProfileRepository) *GetUserProfileUseCase {
	return &GetUserProfileUseCase{profilesRepo}
}

func (u *GetUserProfileUseCase) GetUserProfile(ctx context.Context, userId uuid.UUID) (*entity.UserProfile, error) {
	userProfile, err := u.profilesRepo.Get(ctx, userId)

	if err != nil {
		return nil, err
	}

	return userProfile, nil
}
