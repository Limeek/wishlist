package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type PutItemInWishlistUseCase struct {
	itemRepo     domain.ItemRepository
	wishlistRepo domain.WishlistRepository
}

func NewPutItemInWishListUseCase(itemRepo domain.ItemRepository, wishlistRepo domain.WishlistRepository) *PutItemInWishlistUseCase {
	return &PutItemInWishlistUseCase{itemRepo: itemRepo, wishlistRepo: wishlistRepo}
}

func (uc *PutItemInWishlistUseCase) Put(
	ctx context.Context,
	itemId uuid.UUID,
	wishlistId uuid.UUID,
	userId uuid.UUID,
) error {
	ownerId, err := uc.wishlistRepo.GetWishlistOwner(ctx, wishlistId)

	if err != nil {
		return err
	}

	if ownerId != userId {
		return domain.ErrForbidden
	}

	return uc.itemRepo.PutItemInWishlist(ctx, itemId, wishlistId)
}
