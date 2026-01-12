package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type ReserveItemUseCase struct {
	itemRepo domain.ItemRepository
}

func NewReserveItemUseCase(itemRepo domain.ItemRepository) *ReserveItemUseCase {
	return &ReserveItemUseCase{itemRepo: itemRepo}
}

func (uc *ReserveItemUseCase) Reserve(ctx context.Context, itemId uuid.UUID, userId uuid.UUID) error {
	ownerId, err := uc.itemRepo.GetItemOwnerId(ctx, itemId)

	if err != nil {
		return err
	}

	if ownerId == userId {
		return domain.ErrCantBeReservedByOwner
	}

	return uc.itemRepo.UpdateItemReserve(ctx, itemId, userId)
}
