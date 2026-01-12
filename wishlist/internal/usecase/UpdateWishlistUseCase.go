package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type UpdateWishlistUseCase struct {
	wishlistRepo domain.WishlistRepository
}

func NewUpdateWishlistUseCase(wishlistRepo domain.WishlistRepository) *UpdateWishlistUseCase {
	return &UpdateWishlistUseCase{wishlistRepo: wishlistRepo}
}

func (uc *UpdateWishlistUseCase) Update(ctx context.Context, userId uuid.UUID, wishlist *domain.Wishlist) error {
	ownerId, err := uc.wishlistRepo.GetWishlistOwner(ctx, wishlist.Id)

	if err != nil {
		return err
	}

	if ownerId != userId {
		return domain.ErrForbidden
	}

	return uc.wishlistRepo.Update(ctx, wishlist)
}
