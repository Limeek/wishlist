package usecase

import (
	"context"
	"userprofiles/internal/domain"
)

type CreateUserProfileUseCase struct {
	userProfileRepository domain.UserProfileRepository
}

func NewCreateUserProfileUseCase(userProfileRepo domain.UserProfileRepository) *CreateUserProfileUseCase {
	return &CreateUserProfileUseCase{
		userProfileRepository: userProfileRepo,
	}
}

func (uc *CreateUserProfileUseCase) CreateUserProfile(context context.Context, profile domain.UserProfile) error {
	err := uc.userProfileRepository.Create(context, profile)
	return err
}
