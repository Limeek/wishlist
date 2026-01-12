package domain

import (
	"context"

	"github.com/google/uuid"
)

type WishlistRepository interface {
	GetWishlists(ctx context.Context, ownerId uuid.UUID) ([]Wishlist, error)
	Create(ctx context.Context, wishlist *Wishlist) error
	Delete(ctx context.Context, wishlistId uuid.UUID) error
	Update(ctx context.Context, wishlist *Wishlist) error
	GetWishlist(ctx context.Context, id uuid.UUID) (*Wishlist, error)
	GetWishlistOwner(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}
