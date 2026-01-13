package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type DeleteItemUseCase struct {
	itemRepo     repository.ItemRepository
	wishlistRepo repository.WishlistRepository
}

func NewDeleteItemUseCase(itemRepo repository.ItemRepository, wishlistRepo repository.WishlistRepository) *DeleteItemUseCase {
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
		return entity.ErrForbidden
	}

	return uc.itemRepo.DeleteItem(ctx, itemId, wishlistId)
}
