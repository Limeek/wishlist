package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"wishlist/internal/entity"
	"wishlist/internal/handler/dto"
	"wishlist/internal/handler/mapper"
	"wishlist/internal/usecase"
	"wishlist/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GetMyWishlistsResponse struct {
	Wishlists []dto.WishlistDTO `json:"wishlists"`
}

type GetMyWishlistsHandler struct {
	getWishlistsUseCase *usecase.GetWishlistsUseCase
}

func NewGetMyWishlistsHandler(getWishlistsUseCase *usecase.GetWishlistsUseCase) *GetMyWishlistsHandler {
	return &GetMyWishlistsHandler{
		getWishlistsUseCase: getWishlistsUseCase,
	}
}

func (h *GetMyWishlistsHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Get("/my/wishlist", h.getWishlists)
}

func (h *GetMyWishlistsHandler) getWishlists(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("GetMyWishlists: invalid userId")
		return
	}

	wishlists, err := h.getWishlistsUseCase.Get(r.Context(), userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("GetMyWishlists: error getting wishlists: %v", err)
		return
	}

	response := h.createResponse(wishlists)

	err = json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("GetMyWishlists: error encoding response: %v", err)
		return
	}
}

func (h *GetMyWishlistsHandler) createResponse(wishlists []entity.Wishlist) GetMyWishlistsResponse {
	wishlistResponseElements := make([]dto.WishlistDTO, len(wishlists))

	for i, wishlist := range wishlists {
		wishlistResponseElements[i] = mapper.ToWishlistDto(wishlist)
	}

	return GetMyWishlistsResponse{
		Wishlists: wishlistResponseElements,
	}
}
