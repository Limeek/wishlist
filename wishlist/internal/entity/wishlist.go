package entity

import (
	"time"

	"github.com/google/uuid"
)

type Wishlist struct {
	Id        uuid.UUID
	OwnerId   uuid.UUID
	Name      string
	ExpiresAt time.Time
	Items     []Item
}
