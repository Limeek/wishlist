package main

import (
	"log"
	"net/http"
	"wishlist/internal/handler"
	"wishlist/internal/repository"
	"wishlist/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := repository.NewPostgresConnection()

	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	wishlistRepo := repository.NewWishlistRepositoryPostgres(db)
	itemRepo := repository.NewItemRepositoryPostgres(db)

	createItemUseCase := usecase.NewCreateItemUseCase(itemRepo, wishlistRepo)
	createWishlistUseCase := usecase.NewCreateWishlistUseCase(wishlistRepo)
	deleteItemUseCase := usecase.NewDeleteItemUseCase(itemRepo, wishlistRepo)
	deleteWishlistUseCase := usecase.NewDeleteWishlistUseCase(wishlistRepo)
	getWishlistUseCase := usecase.NewGetWishlistUseCase(wishlistRepo)
	getWishlistsUseCase := usecase.NewGetWishlistsUseCase(wishlistRepo)
	putItemInWishlistUseCase := usecase.NewPutItemInWishListUseCase(itemRepo, wishlistRepo)
	reserveItemUseCase := usecase.NewReserveItemUseCase(itemRepo)
	updateItemUseCase := usecase.NewUpdateItemUseCase(itemRepo)
	updateWishListUseCase := usecase.NewUpdateWishlistUseCase(wishlistRepo)

	createItemHandler := handler.NewCreateItemHandler(createItemUseCase)
	createWishlistHandler := handler.NewCreateWishlistHandler(createWishlistUseCase)
	deleteItemHandler := handler.NewDeleteItemHandler(deleteItemUseCase)
	deleteWishlistHandler := handler.NewDeleteWishlistHandler(deleteWishlistUseCase)
	getMyWishlistsHandler := handler.NewGetMyWishlistsHandler(getWishlistsUseCase)
	getWishlistsHandler := handler.NewGetWishlistsHandler(getWishlistsUseCase)
	getWishlistHandler := handler.NewGetWishlistHandler(getWishlistUseCase)
	putItemInWishlistHandler := handler.NewPutItemInWishlistHandler(putItemInWishlistUseCase)
	reserveItemHandler := handler.NewReserveItemHandler(reserveItemUseCase)
	updateItemHandler := handler.NewUpdateItemHandler(updateItemUseCase)
	updateWishlistHandler := handler.NewUpdateWishlistHandler(updateWishListUseCase)

	r := chi.NewRouter()

	createItemHandler.RegisterRoutes(r)
	createWishlistHandler.RegisterRoutes(r)
	deleteItemHandler.RegisterRoutes(r)
	deleteWishlistHandler.RegisterRoutes(r)
	getMyWishlistsHandler.RegisterRoutes(r)
	getWishlistsHandler.RegisterRoutes(r)
	getWishlistHandler.RegisterRoutes(r)
	putItemInWishlistHandler.RegisterRoutes(r)
	reserveItemHandler.RegisterRoutes(r)
	updateItemHandler.RegisterRoutes(r)
	updateWishlistHandler.RegisterRoutes(r)

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatalf("failed to run http server: %v", err)
	}
}
