# Wishlist Service

Реализует логику wishlist’ов и айтемов внутри.

## Функционал

- создать wishlist
- добавить item
- добавить существующий item в свой wishlist 
- удалить item
- получить свои wishlist’ы
- получить wishlist другого пользователя
- получить вишлист с айтемами

## API

- POST /my/wishlist/{id}/items - создание айтема в вишлисте
- POST /my/wishlist - создание вишлиста 
- DELETE /my/wishlist/{wishlistId}/items/{itemId} - удаление айтема из вишлиста (если нигде не остался, то и удаление айтема)
- DELETE /my/wishlist/{id} - удаление вишлиста 
- GET /my/wishlist - получение вишлистов авторизованного пользователя 
- GET /wishlist/{id} - получение вишлиста с айтемами 
- GET /users/{userId}/wishlist - получение вишлистов пользователя 
- PUT /my/wishlist/{wishlistId}/items/{itemId} - добавить айтем в другой вишлист
- POST /items/{id}/reserve - зарезервировать айтем
- POST /my/items/{id} - апдейт айтема 
- POST /my/wishlist/{id} - апдейт вишлиста 

## Данные

PostgreSQL:
- `wishlists`
- `items`
- `wishlists_items` (join table)

## Запуск
Создать .env файл по примеру из .env.example

docker compose up --build