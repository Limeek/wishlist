package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"wishlist/internal/entity"
	"wishlist/internal/usecase"
	"wishlist/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateItemRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"imageUrl"`
	Link        string  `json:"link"`
	Description string  `json:"description"`
}

type CreateItemHandler struct {
	createItemUseCase *usecase.CreateItemUseCase
}

func NewCreateItemHandler(createItemUseCase *usecase.CreateItemUseCase) *CreateItemHandler {
	return &CreateItemHandler{createItemUseCase: createItemUseCase}
}

func (h *CreateItemHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Post("/my/wishlist/{id}/items", h.createItem)
}

func (h *CreateItemHandler) createItem(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Default().Printf("CreateItem: malformed userId")
		return
	}

	wishlistIdStr := chi.URLParam(r, "id")
	wishlistId, err := uuid.Parse(wishlistIdStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("CreateItem: Invalid Wishlist Id: %v\n", err)
		return
	}

	var req CreateItemRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("CreateItem: Invalid request: %v\n", err)
		return
	}

	item := &entity.Item{
		Name:        req.Name,
		Price:       req.Price,
		ImageUrl:    req.ImageUrl,
		Link:        req.Link,
		Description: req.Description,
	}

	err = h.createItemUseCase.Create(r.Context(), item, wishlistId, userId)

	switch {
	case errors.Is(err, entity.ErrForbidden):
		w.WriteHeader(http.StatusForbidden)
		log.Default().Printf("CreateItem: Forbidden: %v\n", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("CreateItem: CreateItem Error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
