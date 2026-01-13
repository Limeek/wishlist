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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type LoginHandler struct {
	useCase *usecase.LoginUseCase
}

func NewLoginHandler(useCase *usecase.LoginUseCase) *LoginHandler {
	return &LoginHandler{useCase: useCase}
}

func (h *LoginHandler) RegisterRoutes(router chi.Router) {
	router.Post("/auth/login", h.Login)
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Login: error decoding login request: %v", err)
		return
	}

	tokens, err := h.useCase.Login(r.Context(), req.Email, req.Password)

	switch {
	case errors.Is(err, entity.ErrUserNotFound):
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Login: %v", err)
		return
	case errors.Is(err, entity.ErrUserInvalidPassword):
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Login: %v", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("Login: internal server error: %v", err)
		return
	}

	response := LoginResponse{
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("Login: internal server error: %v", err)
		return
	}
}
