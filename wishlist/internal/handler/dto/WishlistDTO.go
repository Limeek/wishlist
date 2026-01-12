package dto

import "github.com/google/uuid"

type WishlistDTO struct {
	Id        uuid.UUID `json:"id"`
	OwnerId   uuid.UUID `json:"ownerId"`
	Name      string    `json:"name"`
	ExpiresAt string    `json:"expiresAt"`
}
