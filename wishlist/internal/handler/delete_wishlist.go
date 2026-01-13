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

type DeleteWishlistHandler struct {
	deleteWishlistUseCase *usecase.DeleteWishlistUseCase
}

func NewDeleteWishlistHandler(deleteWishlistUseCase *usecase.DeleteWishlistUseCase) *DeleteWishlistHandler {
	return &DeleteWishlistHandler{
		deleteWishlistUseCase: deleteWishlistUseCase,
	}
}

func (h *DeleteWishlistHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Delete("/my/wishlist/{id}", h.delete)
}

func (h *DeleteWishlistHandler) delete(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	wishlistIdStr := chi.URLParam(r, "id")
	wishlistId, err := uuid.Parse(wishlistIdStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("DeleteWishlist: invalid wishlistId: %v", err)
		return
	}

	err = h.deleteWishlistUseCase.Delete(r.Context(), wishlistId, userId)

	switch {
	case errors.Is(err, entity.ErrForbidden):
		w.WriteHeader(http.StatusForbidden)
		log.Default().Printf("DeleteWishlist: forbidden: %v", err)
		return
	case errors.Is(err, entity.ErrWishlistNotFound):
		w.WriteHeader(http.StatusNotFound)
		log.Default().Printf("DeleteWishlist: not found: %v", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("DeleteWishlist: failed to delete wishlist: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
