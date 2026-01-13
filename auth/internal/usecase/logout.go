package usecase

import (
	"context"
	"wishlist-auth/internal/repository"

	"github.com/google/uuid"
)

type LogoutUseCase struct {
	refreshTokenRepository repository.RefreshTokenRepository
}

func NewLogoutUseCase(refreshTokenRepository repository.RefreshTokenRepository) *LogoutUseCase {
	return &LogoutUseCase{
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (uc *LogoutUseCase) Logout(ctx context.Context, userId uuid.UUID) error {
	return uc.refreshTokenRepository.Delete(ctx, userId)
}
