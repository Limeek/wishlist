package mapper

import (
	"time"
	"wishlist/internal/domain"
	"wishlist/internal/handler/dto"
)

func ToWishlistDto(wishlist domain.Wishlist) dto.WishlistDTO {
	return dto.WishlistDTO{
		Id:        wishlist.Id,
		Name:      wishlist.Name,
		OwnerId:   wishlist.OwnerId,
		ExpiresAt: wishlist.ExpiresAt.Format(time.RFC3339),
	}
}
