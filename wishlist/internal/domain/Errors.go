package domain

import "errors"

var ErrWishlistNotFound = errors.New("wishlist not found")
var ErrItemNotFound = errors.New("item not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")
var ErrCantBeReservedByOwner = errors.New("can't be reserved by owner")
