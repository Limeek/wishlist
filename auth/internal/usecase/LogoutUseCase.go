package usecase

import (
	"context"
	"wishlist-auth/internal/domain"

	"github.com/google/uuid"
)

type LogoutUseCase struct {
	refreshTokenRepository domain.RefreshTokenRepository
}

func NewLogoutUseCase(refreshTokenRepository domain.RefreshTokenRepository) *LogoutUseCase {
	return &LogoutUseCase{
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (uc *LogoutUseCase) Logout(ctx context.Context, userId uuid.UUID) error {
	return uc.refreshTokenRepository.Delete(ctx, userId)
}
