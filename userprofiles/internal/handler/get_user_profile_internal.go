package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
	"userprofiles/internal/entity"
	"userprofiles/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GetUserProfileResponse struct {
	UserId    uuid.UUID `json:"userId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phoneNumber"`
	Avatar    string    `json:"avatar"`
	BirthDate string    `json:"birthDate"`
}

type GetUserProfileInternalHandler struct {
	getUserProfileUseCase *usecase.GetUserProfileUseCase
}

func NewGetUserProfileInternalHandler(getUserProfileUseCase *usecase.GetUserProfileUseCase) *GetUserProfileInternalHandler {
	return &GetUserProfileInternalHandler{
		getUserProfileUseCase: getUserProfileUseCase,
	}
}

func (h *GetUserProfileInternalHandler) RegisterRoutes(router chi.Router) {
	router.Get("/internal/userprofile", h.getUserProfile)
}

func (h *GetUserProfileInternalHandler) getUserProfile(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("X-Internal-Token") != os.Getenv("INTERNAL_SERVICE_SECRET") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userId, err := uuid.Parse(r.URL.Query().Get("userId"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("missing user id"))
		if err != nil {
			return
		}
		return
	}

	userProfile, err := h.getUserProfileUseCase.GetUserProfile(r.Context(), userId)

	if errors.Is(err, entity.ErrUserNotFound) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("user not found"))
		if err != nil {
			return
		}
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	response := GetUserProfileResponse{
		UserId:    userProfile.UserId,
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
		Email:     userProfile.Email,
		Phone:     userProfile.Phone,
		Avatar:    userProfile.Avatar,
		BirthDate: userProfile.BirthDate.Format(time.DateOnly),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		return
	}
}
