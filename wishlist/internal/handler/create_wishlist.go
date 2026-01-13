package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"wishlist/internal/entity"
	"wishlist/internal/usecase"
	"wishlist/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateWishlistRequest struct {
	Name      string `json:"name"`
	ExpiresAt string `json:"expiresAt"`
}

type CreateWishlistHandler struct {
	createWishlistUseCase *usecase.CreateWishlistUseCase
}

func NewCreateWishlistHandler(createWishlistUseCase *usecase.CreateWishlistUseCase) *CreateWishlistHandler {
	return &CreateWishlistHandler{
		createWishlistUseCase: createWishlistUseCase,
	}
}

func (h *CreateWishlistHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Post("/my/wishlist", h.createWishlist)
}

func (h *CreateWishlistHandler) createWishlist(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Default().Printf("CreateWishlist: malformed userId")
		return
	}

	var req CreateWishlistRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("CreateWishlist: failed to decode request body: %v", err)
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("CreateWishlist: failed to parse expiresAt: %v", err)
		return
	}

	wishlist := &entity.Wishlist{
		Name:      req.Name,
		ExpiresAt: expiresAt,
		OwnerId:   userId,
	}

	err = h.createWishlistUseCase.Create(r.Context(), wishlist)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("CreateWishlist: failed to create wishlist: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
