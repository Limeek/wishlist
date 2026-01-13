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

type ReserveItemHandler struct {
	reserveItemUseCase *usecase.ReserveItemUseCase
}

func NewReserveItemHandler(reserveItemUseCase *usecase.ReserveItemUseCase) *ReserveItemHandler {
	return &ReserveItemHandler{reserveItemUseCase: reserveItemUseCase}
}

func (h *ReserveItemHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Post("/items/{id}/reserve", h.reserve)
}

func (h *ReserveItemHandler) reserve(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Default().Printf("PutItemInWishlist: malformed userId")
		return
	}

	itemIdStr := chi.URLParam(r, "id")
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("PutItemInWishlist: invalid item id: %v", err)
		return
	}

	err = h.reserveItemUseCase.Reserve(r.Context(), itemId, userId)

	switch {
	case errors.Is(err, entity.ErrItemNotFound):
		w.WriteHeader(http.StatusNotFound)
		log.Default().Printf("PutItemInWishlist: item not found: %v", err)
		return
	case errors.Is(err, entity.ErrCantBeReservedByOwner):
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("PutItemInWishlist: can't be reserved by owner: %v", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("PutItemInWishlist: internal server error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
