package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type PutItemInWishlistUseCase struct {
	itemRepo     repository.ItemRepository
	wishlistRepo repository.WishlistRepository
}

func NewPutItemInWishListUseCase(itemRepo repository.ItemRepository, wishlistRepo repository.WishlistRepository) *PutItemInWishlistUseCase {
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
		return entity.ErrForbidden
	}

	return uc.itemRepo.PutItemInWishlist(ctx, itemId, wishlistId)
}
