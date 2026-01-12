package usecase

import (
	"context"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type CreateWishlistUseCase struct {
	wishlistRepo domain.WishlistRepository
}

func NewCreateWishlistUseCase(wishlistRepo domain.WishlistRepository) *CreateWishlistUseCase {
	return &CreateWishlistUseCase{wishlistRepo: wishlistRepo}
}

func (uc *CreateWishlistUseCase) Create(ctx context.Context, wishlist *domain.Wishlist) error {
	wishlist.Id = uuid.New()
	return uc.wishlistRepo.Create(ctx, wishlist)
}
