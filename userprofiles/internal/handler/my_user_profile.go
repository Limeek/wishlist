package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"userprofiles/internal/entity"
	"userprofiles/internal/usecase"
	"userprofiles/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MyUserProfileHandler struct {
	getUserProfileUseCase *usecase.GetUserProfileUseCase
}

func NewMyUserProfileHandler(getUserProfileUseCase *usecase.GetUserProfileUseCase) *MyUserProfileHandler {
	return &MyUserProfileHandler{
		getUserProfileUseCase: getUserProfileUseCase,
	}
}

func (h *MyUserProfileHandler) RegisterRoutes(router chi.Router) {
	router.With(auth.JWTMiddleware()).Get("/userprofile/my", h.getUserProfile)
}

func (h *MyUserProfileHandler) getUserProfile(w http.ResponseWriter, req *http.Request) {
	userId, ok := req.Context().Value("userId").(uuid.UUID)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("userId is required"))
		if err != nil {
			return
		}
		return
	}

	userProfile, err := h.getUserProfileUseCase.GetUserProfile(req.Context(), userId)

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
