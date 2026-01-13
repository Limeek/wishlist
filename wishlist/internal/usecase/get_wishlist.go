package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type GetWishlistUseCase struct {
	wishlistRepo repository.WishlistRepository
}

func NewGetWishlistUseCase(wishlistRepo repository.WishlistRepository) *GetWishlistUseCase {
	return &GetWishlistUseCase{
		wishlistRepo: wishlistRepo,
	}
}

func (uc *GetWishlistUseCase) Get(ctx context.Context, wishlistId uuid.UUID) (*entity.Wishlist, error) {
	return uc.wishlistRepo.GetWishlist(ctx, wishlistId)
}
