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

type UpdateItemRequest struct {
	Name        string
	Price       float64
	ImageUrl    string
	Link        string
	Description string
}

type UpdateItemHandler struct {
	updateItemUseCase *usecase.UpdateItemUseCase
}

func NewUpdateItemHandler(updateUseCase *usecase.UpdateItemUseCase) *UpdateItemHandler {
	return &UpdateItemHandler{
		updateItemUseCase: updateUseCase,
	}
}

func (h *UpdateItemHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Post("/my/items/{id}", h.update)
}

func (h *UpdateItemHandler) update(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Default().Printf("UpdateItem: malformed userId")
		return
	}

	itemIdStr := chi.URLParam(r, "id")
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("UpdateItem: invalid item id: %s", itemIdStr)
	}

	var request UpdateItemRequest

	err = json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("UpdateItem: invalid request: %s", r.Body)
		return
	}

	err = h.updateItemUseCase.Update(r.Context(), userId, &entity.Item{
		Id:          itemId,
		Name:        request.Name,
		Price:       request.Price,
		ImageUrl:    request.ImageUrl,
		Link:        request.Link,
		Description: request.Description,
	})

	switch {
	case errors.Is(err, entity.ErrItemNotFound):
		w.WriteHeader(http.StatusNotFound)
		log.Default().Printf("UpdateItem: item not found: %v", err)
		return
	case errors.Is(err, entity.ErrForbidden):
		w.WriteHeader(http.StatusForbidden)
		log.Default().Printf("UpdateItem: forbidden: %v", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("UpdateItem: internal server error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
