package usecase

import (
	"context"
	"errors"
	"regexp"
	"time"
	"wishlist-auth/internal/domain"
	"wishlist-auth/internal/publisher"
	"wishlist-auth/internal/service"

	"github.com/google/uuid"
)

type RegisterUseCase struct {
	authRepo                     domain.UserRepository
	refreshTokenRepo             domain.RefreshTokenRepository
	passwordHashService          *service.PasswordHashService
	jwtTokenService              *service.JwtTokenService
	userRegisteredEventPublisher *publisher.UserRegisteredEventPublisher
}

func NewRegisterUseCase(
	authRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	passwordHashService *service.PasswordHashService,
	tokenService *service.JwtTokenService,
	userRegisteredEventPublisher *publisher.UserRegisteredEventPublisher,
) *RegisterUseCase {
	return &RegisterUseCase{
		authRepo:                     authRepo,
		refreshTokenRepo:             refreshTokenRepo,
		passwordHashService:          passwordHashService,
		jwtTokenService:              tokenService,
		userRegisteredEventPublisher: userRegisteredEventPublisher,
	}
}

func (uc *RegisterUseCase) RegisterUser(ctx context.Context, user *domain.User, details *domain.UserDetails) (*domain.UserTokens, error) {

	err := validateUser(user)

	if err != nil {
		return nil, err
	}

	err = validateUserDetails(details)

	if err != nil {
		return nil, err
	}

	existingUser, err := uc.authRepo.GetByEmail(ctx, user.Email)

	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	passwordHash, err := uc.passwordHashService.HashPassword(user.Password)

	if err != nil {
		return nil, err
	}

	user.ID = uuid.New()

	err = uc.authRepo.Save(ctx, &domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: passwordHash,
	})

	if err != nil {
		return nil, err
	}

	token, err := uc.jwtTokenService.GenerateToken(user.ID)

	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtTokenService.GenerateRefreshToken(user.ID)

	if err != nil {
		return nil, err
	}

	err = uc.refreshTokenRepo.Save(ctx, refreshToken)

	if err != nil {
		return nil, err
	}

	tokens := &domain.UserTokens{
		Token:        token,
		RefreshToken: refreshToken.Token,
	}

	err = uc.userRegisteredEventPublisher.PublishUserRegisteredEvent(ctx, publisher.UserRegisteredEvent{
		UserId:    user.ID.String(),
		Email:     user.Email,
		Phone:     details.Phone,
		FirstName: details.FirstName,
		LastName:  details.LastName,
		AvatarURL: details.AvatarURL,
		BirthDate: details.BirthDate,
	})

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func validateUser(user *domain.User) error {
	valid, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, user.Email)
	if err != nil {
		return err
	}
	if !valid {
		return domain.ErrEmailInvalidFormat
	}

	if user.Password == "" {
		return domain.ErrPasswordInvalidFormat
	}

	if len(user.Password) < 6 {
		return domain.ErrPasswordInvalidFormat
	}

	return nil
}

func validateUserDetails(details *domain.UserDetails) error {
	valid, err := regexp.MatchString(`^(8|\+7)(\s|\(|-)?(\d{3})(\s|\)|-)?(\d{3})(\s|-)?(\d{2})(\s|-)?(\d{2})$`, details.Phone)

	if err != nil {
		return err
	}
	if !valid {
		return domain.ErrPhoneInvalidFormat
	}

	_, err = time.Parse(time.DateOnly, details.BirthDate)
	if err != nil {
		return domain.ErrBirthDateInvalidFormat
	}

	return nil
}
