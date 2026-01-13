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

type GetWishlistResponse struct {
	Id        uuid.UUID      `json:"id"`
	OwnerId   uuid.UUID      `json:"ownerId"`
	Name      string         `json:"name"`
	ExpiresAt string         `json:"expiresAt"`
	Items     []ItemResponse `json:"items"`
}

type ItemResponse struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Price       float64    `json:"price"`
	ImageUrl    string     `json:"imageUrl"`
	Link        string     `json:"link"`
	Description string     `json:"description"`
	ReservedBy  *uuid.UUID `json:"reservedBy,omitempty"`
}

type GetWishlistHandler struct {
	getWishlistUseCase *usecase.GetWishlistUseCase
}

func NewGetWishlistHandler(getWishlistUseCase *usecase.GetWishlistUseCase) *GetWishlistHandler {
	return &GetWishlistHandler{
		getWishlistUseCase: getWishlistUseCase,
	}
}

func (h *GetWishlistHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Get("/wishlist/{id}", h.getWishlist)
}

func (h *GetWishlistHandler) getWishlist(w http.ResponseWriter, r *http.Request) {
	wishlistIdStr := chi.URLParam(r, "id")
	wishlistId, err := uuid.Parse(wishlistIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("GetWishlist: invalid userId: %v", err)
		return
	}

	wishlist, err := h.getWishlistUseCase.Get(r.Context(), wishlistId)

	switch {
	case errors.Is(err, entity.ErrWishlistNotFound):
		w.WriteHeader(http.StatusNotFound)
		log.Default().Printf("GetWishlist: wishlist not found: %v", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("GetWishlist: internal server error: %v", err)
		return
	}

	response := h.createResponse(wishlist)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("GetWishlist: encoding response: %v", err)
		return
	}
}

func (h *GetWishlistHandler) createResponse(wishlist *entity.Wishlist) GetWishlistResponse {
	items := make([]ItemResponse, len(wishlist.Items))
	for i, item := range wishlist.Items {
		items[i] = ItemResponse{
			Id:          item.Id,
			Name:        item.Name,
			Price:       item.Price,
			ImageUrl:    item.ImageUrl,
			Description: item.Description,
			ReservedBy:  item.ReservedBy,
		}
	}

	return GetWishlistResponse{
		Id:        wishlist.Id,
		OwnerId:   wishlist.OwnerId,
		Name:      wishlist.Name,
		ExpiresAt: wishlist.ExpiresAt.Format(time.RFC3339),
		Items:     items,
	}
}
