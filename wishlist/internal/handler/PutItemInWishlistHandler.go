package handler

import (
	"errors"
	"log"
	"net/http"
	"wishlist/internal/domain"
	"wishlist/internal/usecase"
	"wishlist/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PutItemInWishlistHandler struct {
	putItemInWishlistUseCase *usecase.PutItemInWishlistUseCase
}

func NewPutItemInWishlistHandler(putItemInWishlistUseCase *usecase.PutItemInWishlistUseCase) *PutItemInWishlistHandler {
	return &PutItemInWishlistHandler{
		putItemInWishlistUseCase: putItemInWishlistUseCase,
	}
}

func (h *PutItemInWishlistHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Put("/my/wishlist/{wishlistId}/items/{itemId}", h.putItem)
}

func (h *PutItemInWishlistHandler) putItem(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Default().Printf("PutItemInWishlist: malformed userId")
		return
	}

	wishlistIdStr := chi.URLParam(r, "wishlistId")
	wishlistId, err := uuid.Parse(wishlistIdStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("PutItemInWishlist: invalid wishlistId: %v", err)
		return
	}

	itemIdStr := chi.URLParam(r, "itemId")
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("PutItemInWishlist: invalid itemId: %v", err)
		return
	}

	err = h.putItemInWishlistUseCase.Put(r.Context(), itemId, wishlistId, userId)

	switch {
	case errors.Is(err, domain.ErrItemNotFound):
		w.WriteHeader(http.StatusNotFound)
		log.Default().Printf("PutItemInWishlist: item not found: %v", err)
		return
	case errors.Is(err, domain.ErrForbidden):
		w.WriteHeader(http.StatusForbidden)
		log.Default().Printf("PutItemInWishlist: forbidden: %v", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("PutItemInWishlist: internal server error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
