package repository

import (
	"context"
	"wishlist/internal/entity"

	"github.com/google/uuid"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, item *entity.Item, wishlistId uuid.UUID) error
	DeleteItem(ctx context.Context, itemId uuid.UUID, wishlistId uuid.UUID) error
	UpdateItem(ctx context.Context, item *entity.Item) error
	UpdateItemReserve(ctx context.Context, itemId uuid.UUID, reservedBy uuid.UUID) error
	PutItemInWishlist(ctx context.Context, itemId uuid.UUID, wishlistId uuid.UUID) error
	GetItemOwnerId(ctx context.Context, itemId uuid.UUID) (uuid.UUID, error)
}

type WishlistRepository interface {
	GetWishlists(ctx context.Context, ownerId uuid.UUID) ([]entity.Wishlist, error)
	Create(ctx context.Context, wishlist *entity.Wishlist) error
	Delete(ctx context.Context, wishlistId uuid.UUID) error
	Update(ctx context.Context, wishlist *entity.Wishlist) error
	GetWishlist(ctx context.Context, id uuid.UUID) (*entity.Wishlist, error)
	GetWishlistOwner(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}
