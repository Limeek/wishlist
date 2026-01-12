package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type DeleteItemUseCase struct {
	itemRepo     domain.ItemRepository
	wishlistRepo domain.WishlistRepository
}

func NewDeleteItemUseCase(itemRepo domain.ItemRepository, wishlistRepo domain.WishlistRepository) *DeleteItemUseCase {
	return &DeleteItemUseCase{itemRepo: itemRepo, wishlistRepo: wishlistRepo}
}

func (uc *DeleteItemUseCase) Delete(
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

	return uc.itemRepo.DeleteItem(ctx, itemId, wishlistId)
}
