package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type GetWishlistUseCase struct {
	wishlistRepo domain.WishlistRepository
}

func NewGetWishlistUseCase(wishlistRepo domain.WishlistRepository) *GetWishlistUseCase {
	return &GetWishlistUseCase{
		wishlistRepo: wishlistRepo,
	}
}

func (uc *GetWishlistUseCase) Get(ctx context.Context, wishlistId uuid.UUID) (*domain.Wishlist, error) {
	return uc.wishlistRepo.GetWishlist(ctx, wishlistId)
}
