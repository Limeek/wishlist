package handler

import (
	"log"
	"net/http"
	"wishlist-auth/internal/usecase"
	"wishlist-auth/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type LogoutHandler struct {
	logoutUseCase *usecase.LogoutUseCase
}

func NewLogoutHandler(logoutUseCase *usecase.LogoutUseCase) *LogoutHandler {
	return &LogoutHandler{
		logoutUseCase: logoutUseCase,
	}
}

func (h *LogoutHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Post("/auth/logout", h.Logout)
}

func (h *LogoutHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Default().Printf("Logout: context missing user id")
		return
	}

	err := h.logoutUseCase.Logout(r.Context(), userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("Logout: internal server error: %v", err)
		return
	}
}
