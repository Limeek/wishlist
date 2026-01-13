package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type CreateItemUseCase struct {
	itemRepo     repository.ItemRepository
	wishlistRepo repository.WishlistRepository
}

func NewCreateItemUseCase(itemRepo repository.ItemRepository, wishlistRepo repository.WishlistRepository) *CreateItemUseCase {
	return &CreateItemUseCase{itemRepo: itemRepo, wishlistRepo: wishlistRepo}
}

func (uc *CreateItemUseCase) Create(ctx context.Context,
	item *entity.Item,
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

	item.Id = uuid.New()

	return uc.itemRepo.CreateItem(ctx, item, wishlistId)
}
