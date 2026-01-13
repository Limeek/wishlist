package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type UpdateWishlistUseCase struct {
	wishlistRepo repository.WishlistRepository
}

func NewUpdateWishlistUseCase(wishlistRepo repository.WishlistRepository) *UpdateWishlistUseCase {
	return &UpdateWishlistUseCase{wishlistRepo: wishlistRepo}
}

func (uc *UpdateWishlistUseCase) Update(ctx context.Context, userId uuid.UUID, wishlist *entity.Wishlist) error {
	ownerId, err := uc.wishlistRepo.GetWishlistOwner(ctx, wishlist.Id)

	if err != nil {
		return err
	}

	if ownerId != userId {
		return entity.ErrForbidden
	}

	return uc.wishlistRepo.Update(ctx, wishlist)
}
