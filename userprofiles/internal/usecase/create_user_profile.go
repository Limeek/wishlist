package usecase

import (
	"context"
	"userprofiles/internal/entity"
	"userprofiles/internal/repository"
)

type CreateUserProfileUseCase struct {
	userProfileRepository repository.UserProfileRepository
}

func NewCreateUserProfileUseCase(userProfileRepo repository.UserProfileRepository) *CreateUserProfileUseCase {
	return &CreateUserProfileUseCase{
		userProfileRepository: userProfileRepo,
	}
}

func (uc *CreateUserProfileUseCase) CreateUserProfile(context context.Context, profile entity.UserProfile) error {
	err := uc.userProfileRepository.Create(context, profile)
	return err
}
