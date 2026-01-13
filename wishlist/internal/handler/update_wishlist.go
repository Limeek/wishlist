package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	"wishlist/internal/entity"
	"wishlist/internal/usecase"
	"wishlist/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UpdateWishlistRequest struct {
	Name      string
	ExpiresAt string
}

type UpdateWishlistHandler struct {
	updateWishlistUseCase *usecase.UpdateWishlistUseCase
}

func NewUpdateWishlistHandler(updateWishlistUseCase *usecase.UpdateWishlistUseCase) *UpdateWishlistHandler {
	return &UpdateWishlistHandler{
		updateWishlistUseCase: updateWishlistUseCase,
	}
}

func (h *UpdateWishlistHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Post("/my/wishlist/{id}", h.updateWishlist)
}

func (h *UpdateWishlistHandler) updateWishlist(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Default().Printf("UpdateWishlist: malformed userId")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("UpdateWishlist: failed to parse id: %v", err)
		return
	}

	var request UpdateWishlistRequest

	err = json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("UpdateWishlist: failed to parse request: %v", err)
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, request.ExpiresAt)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("UpdateWishlist: failed to parse request: %v", err)
		return
	}

	err = h.updateWishlistUseCase.Update(r.Context(), userId, &entity.Wishlist{
		Id:        id,
		Name:      request.Name,
		ExpiresAt: expiresAt,
	})

	switch {
	case errors.Is(err, entity.ErrWishlistNotFound):
		w.WriteHeader(http.StatusNotFound)
		log.Default().Printf("UpdateWishlist: Wishlist not found: %v", err)
		return
	case errors.Is(err, entity.ErrForbidden):
		w.WriteHeader(http.StatusForbidden)
		log.Default().Printf("UpdateWishlist: Forbidden: %v", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("UpdateWishlist: failed to update: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
