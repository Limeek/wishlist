package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type GetWishlistsUseCase struct {
	wishlistRepo repository.WishlistRepository
}

func NewGetWishlistsUseCase(wishlistRepo repository.WishlistRepository) *GetWishlistsUseCase {
	return &GetWishlistsUseCase{
		wishlistRepo: wishlistRepo,
	}
}

func (uc *GetWishlistsUseCase) Get(ctx context.Context, ownerId uuid.UUID) ([]entity.Wishlist, error) {
	return uc.wishlistRepo.GetWishlists(ctx, ownerId)
}
