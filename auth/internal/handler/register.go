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
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Register: error decoding register request: %v", err)
		return
	}

	user := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	details := entity.UserDetails{
		Phone:     req.Phone,
		BirthDate: req.BirthDate,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		AvatarURL: req.AvatarURL,
	}

	tokens, err := h.registerUseCase.RegisterUser(r.Context(), &user, &details)

	switch {
	case errors.Is(err, entity.ErrEmailInvalidFormat):
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Register: %v", err)
		return
	case errors.Is(err, entity.ErrPasswordInvalidFormat):
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Register: %v", err)
		return
	case errors.Is(err, entity.ErrPhoneInvalidFormat):
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Register: %v", err)
		return
	case errors.Is(err, entity.ErrUserAlreadyExists):
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("Register: %v", err)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("Register: internal server error: %v", err)
		return
	}

	resp := RegisterResponse{
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("Register: internal server error: %v", err)
		return
	}
}
