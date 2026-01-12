package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type UpdateItemUseCase struct {
	itemRepo domain.ItemRepository
}

func NewUpdateItemUseCase(itemRepo domain.ItemRepository) *UpdateItemUseCase {
	return &UpdateItemUseCase{itemRepo: itemRepo}
}

func (uc *UpdateItemUseCase) Update(ctx context.Context, userId uuid.UUID, item *domain.Item) error {
	ownerId, err := uc.itemRepo.GetItemOwnerId(ctx, item.Id)

	if err != nil {
		return err
	}

	if ownerId != userId {
		return domain.ErrForbidden
	}

	return uc.itemRepo.UpdateItem(ctx, item)
}
