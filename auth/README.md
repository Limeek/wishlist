# Auth Service

Отвечает за:
- регистрацию пользователя
- логин
- выдачу JWT
- рефреш JWT 

При регистрации публикует событие user-registered в Kafka для user-profile-service.

## API

- POST /auth/register  - регистрация пользователя

- POST /auth/login - авторизация пользователя

- POST /auth/logout - логаут пользователя 

- POST /auth/refresh - рефреш JWT 

## Данные

PostgreSQL:

- refresh_tokens
- users

## Kafka

topic: user-registered

При регистрации публикует событие user-registered в Kafka для user-profile-service.

## Запуск
Создать .env файл по примеру из .env.example

docker compose up --build