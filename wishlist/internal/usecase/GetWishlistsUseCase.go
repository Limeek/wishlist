package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type GetWishlistsUseCase struct {
	wishlistRepo domain.WishlistRepository
}

func NewGetWishlistsUseCase(wishlistRepo domain.WishlistRepository) *GetWishlistsUseCase {
	return &GetWishlistsUseCase{
		wishlistRepo: wishlistRepo,
	}
}

func (uc *GetWishlistsUseCase) Get(ctx context.Context, ownerId uuid.UUID) ([]domain.Wishlist, error) {
	return uc.wishlistRepo.GetWishlists(ctx, ownerId)
}
