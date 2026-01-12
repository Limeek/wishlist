package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"wishlist-auth/internal/handler"
	"wishlist-auth/internal/publisher"
	"wishlist-auth/internal/repository"
	"wishlist-auth/internal/service"
	"wishlist-auth/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db, err := repository.NewPostgresConnection()

	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepositoryDb(db)

	passwordHashService := service.NewPasswordHashService()
	jwtTokenService := service.NewJwtTokenService()

	userRegisteredEventPublisher := publisher.NewUserRegisteredEventPublisher()

	registerUseCase := usecase.NewRegisterUseCase(userRepo, refreshTokenRepo, passwordHashService, jwtTokenService, userRegisteredEventPublisher)
	loginUseCase := usecase.NewLoginUseCase(userRepo, refreshTokenRepo, jwtTokenService, passwordHashService)
	logoutUseCase := usecase.NewLogoutUseCase(refreshTokenRepo)
	refreshUseCase := usecase.NewRefreshUseCase(userRepo, refreshTokenRepo, jwtTokenService)

	registerHandler := handler.NewRegisterHandler(registerUseCase)
	loginHandler := handler.NewLoginHandler(loginUseCase)
	logoutHandler := handler.NewLogoutHandler(logoutUseCase)
	refreshHandler := handler.NewRefreshHandler(refreshUseCase)

	r := chi.NewRouter()

	registerHandler.RegisterRoutes(r)
	loginHandler.RegisterRoutes(r)
	logoutHandler.RegisterRoutes(r)
	refreshHandler.RegisterRoutes(r)

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatalf("failed to run http server: %v", err)
	}

	<-ctx.Done()
	userRegisteredEventPublisher.Close()
}
