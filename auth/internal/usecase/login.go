package usecase

import (
	"context"
	"wishlist-auth/internal/entity"
	"wishlist-auth/internal/repository"
	"wishlist-auth/internal/service"
)

type LoginUseCase struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	tokenService           *service.JwtService
	passwordHashService    *service.PasswordHashService
}

func NewLoginUseCase(userRepository repository.UserRepository, refreshTokenRepository repository.RefreshTokenRepository, tokenService *service.JwtService, hashService *service.PasswordHashService) *LoginUseCase {
	return &LoginUseCase{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		tokenService:           tokenService,
		passwordHashService:    hashService,
	}
}

func (uc *LoginUseCase) Login(ctx context.Context, username string, password string) (*entity.UserTokens, error) {
	user, err := uc.userRepository.GetByEmail(ctx, username)

	if err != nil {
		return nil, err
	}

	isValidPassword := uc.passwordHashService.CheckPassword(password, user.Password)

	if !isValidPassword {
		return nil, entity.ErrUserInvalidPassword
	}

	token, err := uc.tokenService.GenerateToken(user.ID)

	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)

	if err != nil {
		return nil, err
	}

	err = uc.refreshTokenRepository.Save(ctx, refreshToken)

	if err != nil {
		return nil, err
	}

	return &entity.UserTokens{
		Token:        token,
		RefreshToken: refreshToken.Token,
	}, nil

}
