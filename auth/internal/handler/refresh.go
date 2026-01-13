package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"wishlist-auth/internal/entity"
	"wishlist-auth/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshHandler struct {
	refreshUseCase *usecase.RefreshUseCase
}

func NewRefreshHandler(refreshUseCase *usecase.RefreshUseCase) *RefreshHandler {
	return &RefreshHandler{
		refreshUseCase: refreshUseCase,
	}
}

func (h *RefreshHandler) RegisterRoutes(router chi.Router) {
	router.Post("/auth/refresh", h.RefreshAuth)
}

func (h *RefreshHandler) RefreshAuth(w http.ResponseWriter, r *http.Request) {

	var req RefreshRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Refresh: error decoding body: %v", err)
		return
	}

	tokens, err := h.refreshUseCase.Refresh(r.Context(), req.RefreshToken)

	if errors.Is(err, entity.ErrRefreshTokenIsNotValid) {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Refresh: %v", err)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("Refresh: internal server error: %v", err)
		return
	}

	response := RefreshResponse{
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("Refresh: internal server error: %v", err)
		return
	}
}
