package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"wishlist-auth/internal/domain"
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
		writeError(http.StatusBadRequest, err, w)
		log.Default().Printf("[ERROR] decoding login request: %v", err)
		return
	}

	tokens, err := h.useCase.Login(r.Context(), req.Email, req.Password)

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		writeError(http.StatusBadRequest, err, w)
		return
	case errors.Is(err, domain.ErrUserInvalidPassword):
		writeError(http.StatusBadRequest, err, w)
		return
	default:
		if err != nil {
			writeError(http.StatusInternalServerError, err, w)
			return
		}
	}

	response := LoginResponse{
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		writeError(http.StatusInternalServerError, err, w)
		log.Default().Printf("[ERROR] encoding response: %v", err)
		return
	}
}
