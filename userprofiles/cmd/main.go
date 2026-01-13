package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"userprofiles/internal/consumer"
	"userprofiles/internal/handler"
	"userprofiles/internal/repository"
	"userprofiles/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	sqlDb, err := repository.NewPostgresConnection()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	userProfileRepo := repository.NewUserProfileRepositoryPostgres(sqlDb)

	createProfileUsecase := usecase.NewCreateUserProfileUseCase(userProfileRepo)
	getProfileUsecase := usecase.NewGetUserProfileUseCase(userProfileRepo)

	getUserProfileInternalHandler := handler.NewGetUserProfileInternalHandler(getProfileUsecase)
	myUserProfileHandler := handler.NewMyUserProfileHandler(getProfileUsecase)

	userRegisteredConsumer := consumer.NewUserRegisteredEventConsumer(createProfileUsecase)
	userRegisteredConsumer.Start(ctx)

	router := chi.NewRouter()
	getUserProfileInternalHandler.RegisterRoutes(router)
	myUserProfileHandler.RegisterRoutes(router)

	err = http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
