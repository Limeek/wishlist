package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"wishlist/internal/domain"
	"wishlist/internal/handler/dto"
	"wishlist/internal/handler/mapper"
	"wishlist/internal/usecase"
	"wishlist/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GetWishlistsResponse struct {
	Wishlists []dto.WishlistDTO `json:"wishlists"`
}

type GetWishlistsHandler struct {
	getWishlistsUseCase *usecase.GetWishlistsUseCase
}

func NewGetWishlistsHandler(getWishlistsUseCase *usecase.GetWishlistsUseCase) *GetWishlistsHandler {
	return &GetWishlistsHandler{
		getWishlistsUseCase: getWishlistsUseCase,
	}
}

func (h *GetWishlistsHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Get("/users/{userId}/wishlist", h.getWishlists)
}

func (h *GetWishlistsHandler) getWishlists(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("GetWishlists: invalid userId: %v", err)
		return
	}

	wishlists, err := h.getWishlistsUseCase.Get(r.Context(), userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("GetWishlists: error getting wishlists: %v", err)
		return
	}

	response := h.createResponse(wishlists)

	err = json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("GetWishlists: error encoding response: %v", err)
		return
	}
}

func (h *GetWishlistsHandler) createResponse(wishlists []domain.Wishlist) GetWishlistsResponse {
	wishlistResponseElements := make([]dto.WishlistDTO, len(wishlists))

	for i, wishlist := range wishlists {
		wishlistResponseElements[i] = mapper.ToWishlistDto(wishlist)
	}

	return GetWishlistsResponse{
		Wishlists: wishlistResponseElements,
	}
}
