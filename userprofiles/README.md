# User Profile Service

Отвечает за публичный профиль пользователя.

## API

- GET /internal/userprofile - внутренняя ручка для получения юзер профиля
- GET /userprofile/my - ручка для получения своего профиля юзера (по JWT)

## Данные

PostgreSQL:
- user_profiles


## Kafka

consumer: user-registered
При получении события создает профиль пользователя.

## Запуск
Создать .env файл по примеру из .env.example

docker compose up --build