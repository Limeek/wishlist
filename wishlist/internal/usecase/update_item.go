package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type UpdateItemUseCase struct {
	itemRepo repository.ItemRepository
}

func NewUpdateItemUseCase(itemRepo repository.ItemRepository) *UpdateItemUseCase {
	return &UpdateItemUseCase{itemRepo: itemRepo}
}

func (uc *UpdateItemUseCase) Update(ctx context.Context, userId uuid.UUID, item *entity.Item) error {
	ownerId, err := uc.itemRepo.GetItemOwnerId(ctx, item.Id)

	if err != nil {
		return err
	}

	if ownerId != userId {
		return entity.ErrForbidden
	}

	return uc.itemRepo.UpdateItem(ctx, item)
}
