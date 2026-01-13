package handler

import (
	"errors"
	"log"
	"net/http"
	"wishlist/internal/entity"
	"wishlist/internal/usecase"
	"wishlist/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DeleteItemHandler struct {
	deleteItemUseCase *usecase.DeleteItemUseCase
}

func NewDeleteItemHandler(deleteItemUseCase *usecase.DeleteItemUseCase) *DeleteItemHandler {
	return &DeleteItemHandler{deleteItemUseCase: deleteItemUseCase}
}

func (h *DeleteItemHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Delete("/my/wishlist/{wishlistId}/items/{itemId}", h.deleteItem)
}

func (h *DeleteItemHandler) deleteItem(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Default().Printf("DeleteItem: malformed userId")

		return
	}

	wishlistIdStr := chi.URLParam(r, "wishlistId")
	itemIdStr := chi.URLParam(r, "itemId")

	wishlistId, err := uuid.Parse(wishlistIdStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("DeleteItem: invalid wishlist id %v", err)
		return
	}

	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("DeleteItem: invalid item id %v", err)
		return
	}

	err = h.deleteItemUseCase.Delete(r.Context(), itemId, wishlistId, userId)

	switch {
	case errors.Is(err, entity.ErrWishlistNotFound):
		w.WriteHeader(http.StatusNotFound)
		log.Default().Printf("DeleteItem: Wishlist not found %v", err)
		return
	case errors.Is(err, entity.ErrItemNotFound):
		w.WriteHeader(http.StatusNotFound)
		log.Default().Printf("DeleteItem: item not found %v", err)
		return
	case errors.Is(err, entity.ErrForbidden):
		w.WriteHeader(http.StatusForbidden)
		log.Default().Printf("DeleteItem: forbidden %v", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("DeleteItem: failed to delete wishlist: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
