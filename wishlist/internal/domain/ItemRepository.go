package domain

import (
	"context"

	"github.com/google/uuid"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, item *Item, wishlistId uuid.UUID) error
	DeleteItem(ctx context.Context, itemId uuid.UUID, wishlistId uuid.UUID) error
	UpdateItem(ctx context.Context, item *Item) error
	UpdateItemReserve(ctx context.Context, itemId uuid.UUID, reservedBy uuid.UUID) error
	PutItemInWishlist(ctx context.Context, itemId uuid.UUID, wishlistId uuid.UUID) error
	GetItemOwnerId(ctx context.Context, itemId uuid.UUID) (uuid.UUID, error)
}
