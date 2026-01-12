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

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phoneNumber"`
	BirthDate string `json:"birthDate"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
}

type RegisterResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterHandler struct {
	registerUseCase *usecase.RegisterUseCase
}

func NewRegisterHandler(registerUseCase *usecase.RegisterUseCase) *RegisterHandler {
	return &RegisterHandler{
		registerUseCase: registerUseCase,
	}
}

func (h *RegisterHandler) RegisterRoutes(router chi.Router) {
	router.Post("/auth/register", h.Register)
}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeError(http.StatusBadRequest, err, w)
		log.Default().Printf("[ERROR] decoding register request: %v", err)
		return
	}

	user := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	details := domain.UserDetails{
		Phone:     req.Phone,
		BirthDate: req.BirthDate,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		AvatarURL: req.AvatarURL,
	}

	tokens, err := h.registerUseCase.RegisterUser(r.Context(), &user, &details)

	switch {
	case errors.Is(err, domain.ErrEmailInvalidFormat):
		writeError(http.StatusBadRequest, err, w)
		log.Default().Printf("[ERROR] registering user: %v", err)
		return
	case errors.Is(err, domain.ErrPasswordInvalidFormat):
		writeError(http.StatusBadRequest, err, w)
		log.Default().Printf("[ERROR] registering user: %v", err)
		return
	case errors.Is(err, domain.ErrPhoneInvalidFormat):
		writeError(http.StatusBadRequest, err, w)
		log.Default().Printf("[ERROR] registering user: %v", err)
		return
	case errors.Is(err, domain.ErrUserAlreadyExists):
		writeError(http.StatusBadRequest, err, w)
		log.Default().Printf("[ERROR] registering user: %v", err)
		return
	case err != nil:
		writeError(http.StatusInternalServerError, err, w)
		log.Default().Printf("[ERROR] registering user: %v", err)
		return
	}

	resp := RegisterResponse{
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		writeError(http.StatusInternalServerError, err, w)
		log.Default().Printf("[ERROR] encoding response: %v", err)
		return
	}
}
