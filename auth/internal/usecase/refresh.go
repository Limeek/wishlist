package usecase

import (
	"context"
	"time"
	"wishlist-auth/internal/entity"
	"wishlist-auth/internal/repository"
	"wishlist-auth/internal/service"
)

type RefreshUseCase struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	jwtService             *service.JwtService
}

func NewRefreshUseCase(
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	jwtService *service.JwtService,
) *RefreshUseCase {
	return &RefreshUseCase{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		jwtService:             jwtService,
	}
}

func (uc *RefreshUseCase) Refresh(ctx context.Context, refreshToken string) (*entity.UserTokens, error) {
	token, err := uc.refreshTokenRepository.FindByToken(ctx, refreshToken)

	if err != nil {
		return nil, err
	}

	isExpired := token.ExpiresAt.UTC().Before(time.Now().UTC())

	if isExpired {
		return nil, entity.ErrRefreshTokenIsNotValid
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

	return &entity.UserTokens{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
