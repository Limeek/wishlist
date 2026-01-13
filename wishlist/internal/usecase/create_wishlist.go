package usecase

import (
	"context"
	"wishlist/internal/entity"
	"wishlist/internal/repository"

	"github.com/google/uuid"
)

type CreateWishlistUseCase struct {
	wishlistRepo repository.WishlistRepository
}

func NewCreateWishlistUseCase(wishlistRepo repository.WishlistRepository) *CreateWishlistUseCase {
	return &CreateWishlistUseCase{wishlistRepo: wishlistRepo}
}

func (uc *CreateWishlistUseCase) Create(ctx context.Context, wishlist *entity.Wishlist) error {
	wishlist.Id = uuid.New()
	return uc.wishlistRepo.Create(ctx, wishlist)
}
