package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type ReserveItemUseCase struct {
	itemRepo repository.ItemRepository
}

func NewReserveItemUseCase(itemRepo repository.ItemRepository) *ReserveItemUseCase {
	return &ReserveItemUseCase{itemRepo: itemRepo}
}

func (uc *ReserveItemUseCase) Reserve(ctx context.Context, itemId uuid.UUID, userId uuid.UUID) error {
	ownerId, err := uc.itemRepo.GetItemOwnerId(ctx, itemId)

	if err != nil {
		return err
	}

	if ownerId == userId {
		return entity.ErrCantBeReservedByOwner
	}

	return uc.itemRepo.UpdateItemReserve(ctx, itemId, userId)
}
