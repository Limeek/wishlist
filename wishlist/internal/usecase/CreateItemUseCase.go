package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type CreateItemUseCase struct {
	itemRepo     domain.ItemRepository
	wishlistRepo domain.WishlistRepository
}

func NewCreateItemUseCase(itemRepo domain.ItemRepository, wishlistRepo domain.WishlistRepository) *CreateItemUseCase {
	return &CreateItemUseCase{itemRepo: itemRepo, wishlistRepo: wishlistRepo}
}

func (uc *CreateItemUseCase) Create(ctx context.Context,
	item *domain.Item,
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

	item.Id = uuid.New()

	return uc.itemRepo.CreateItem(ctx, item, wishlistId)
}
