package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type DeleteWishlistUseCase struct {
	wishlistRepo repository.WishlistRepository
}

func NewDeleteWishlistUseCase(wishlistRepo repository.WishlistRepository) *DeleteWishlistUseCase {
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
		return entity.ErrForbidden
	}

	return uc.wishlistRepo.Delete(ctx, wishlistId)
}
