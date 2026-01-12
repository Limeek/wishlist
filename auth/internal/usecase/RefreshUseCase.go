package usecase

import (
	"context"
	"time"
	"wishlist-auth/internal/domain"
	"wishlist-auth/internal/service"
)

type RefreshUseCase struct {
	userRepository         domain.UserRepository
	refreshTokenRepository domain.RefreshTokenRepository
	jwtService             *service.JwtTokenService
}

func NewRefreshUseCase(
	userRepository domain.UserRepository,
	refreshTokenRepository domain.RefreshTokenRepository,
	jwtService *service.JwtTokenService,
) *RefreshUseCase {
	return &RefreshUseCase{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		jwtService:             jwtService,
	}
}

func (uc *RefreshUseCase) Refresh(ctx context.Context, refreshToken string) (*domain.UserTokens, error) {
	token, err := uc.refreshTokenRepository.FindByToken(ctx, refreshToken)

	if err != nil {
		return nil, err
	}

	isExpired := token.ExpiresAt.UTC().Before(time.Now().UTC())

	if isExpired {
		return nil, domain.ErrRefreshTokenIsNotValid
	}

	accessToken, err := uc.jwtService.GenerateToken(token.UserId)

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(token.UserId)

	if err != nil {
		return nil, err
	}

	err = uc.refreshTokenRepository.Save(ctx, newRefreshToken)

	if err != nil {
		return nil, err
	}

	return &domain.UserTokens{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
