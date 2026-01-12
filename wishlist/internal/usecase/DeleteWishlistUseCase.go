package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type DeleteWishlistUseCase struct {
	wishlistRepo domain.WishlistRepository
}

func NewDeleteWishlistUseCase(wishlistRepo domain.WishlistRepository) *DeleteWishlistUseCase {
	return &DeleteWishlistUseCase{
		wishlistRepo: wishlistRepo,
	}
}

func (uc *DeleteWishlistUseCase) Delete(ctx context.Context, wishlistId uuid.UUID, userId uuid.UUID) error {
	ownerId, err := uc.wishlistRepo.GetWishlistOwner(ctx, wishlistId)

	if err != nil {
		return err
	}

	if ownerId != userId {
		return domain.ErrForbidden
	}

	return uc.wishlistRepo.Delete(ctx, wishlistId)
}
