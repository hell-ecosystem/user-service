# User Service

CRUD-сервис для работы с профилем пользователя.

---

## Оглавление

- [Описание](#описание)  
- [Формат ответов API](#формат-ответов-api)  
- [Эндпоинты](#эндпоинты)  
  - [GET /users/me](#get-usersme)  
  - [GET /users/{id}](#get-usersid)  
- [Архитектура](#архитектура)  
- [Миграции](#миграции)  
- [Запуск](#запуск)  
- [Переменные окружения](#переменные-окружения)  

---

## Описание

User Service — это лёгкий HTTP-сервис на Go, который предоставляет доступ к данным пользователей:

- хранит их в PostgreSQL,
- позволяет получать профиль текущего пользователя по JWT,
- возвращает данные в едином JSON-формате.

---

## Формат ответов API

Все ответы приходят в «конверте»:

### Успех

```json
{
  "success": true,
  "data": { /* произвольная полезная нагрузка */ }
}
```

### Ошибка

```json
{
  "success": false,
  "error": {
    "code":    "MACHINE_CODE",
    "message": "человекочитаемое сообщение"
  }
}
```

## Эндпоинты

### GET /users/me

Получить профиль текущего аутентифицированного пользователя.

**Пример запроса:**

```bash
curl -H "Authorization: Bearer <access_token>" \
     http://localhost:8081/users/me
```

**Ответы:**

- 200 OK

  ```json
  {
    "success": true,
    "data": {
      "id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
      "email": "user@example.com",
      "telegram_id": 123456789,
      "created_at": "2025-05-07T12:34:56Z"
    }
  }
  ```

- 401 Unauthorized

  ```json
  {
    "success": false,
    "error": {
      "code":    "UNAUTHORIZED",
      "message": "неавторизованный запрос"
    }
  }
  ```

- 404 Not Found

  ```json
  {
    "success": false,
    "error": {
      "code":    "USER_NOT_FOUND",
      "message": "пользователь не найден"
    }
  }
  ```
  
### GET /users/{id}

Получить профиль пользователя по ID.

**Пример запроса:**

```bash
curl http://localhost:8081/users/d290f1ee-6c54-4b01-90e6-d701748f0851
```

**Ответы:**

- 200 OK

  ```json
  {
    "success": true,
    "data": {
      "id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
      "email": "user@example.com",
      "telegram_id": null,
      "created_at": "2025-05-07T12:34:56Z"
    }
  }
  ```

- 404 Not Found

  ```json
  {
    "success": false,
    "error": {
      "code":    "USER_NOT_FOUND",
      "message": "пользователь не найден"
    }
  }
  ```

- 500 Internal Server Error

  ```json
  {
    "success": false,
    "error": {
      "code":    "INTERNAL_ERROR",
      "message": "внутренняя ошибка сервера"
    }
  }
  ```

## Архитектура

```text
cmd/
  main.go                     — точка входа
internal/
  config/
    config.go                 — загрузка и валидация env
    helper.go                 — утилиты (таймауты, DSN)
  delivery/
    httpdelivery/
      handler.go              — HTTP-хендлеры и маршруты
      response.go             — общий формат ответов
  service/
    user_service.go           — бизнес-логика
  repository/
    postgres/
      repository.go           — запросы к PostgreSQL
  model/
    user.go                   — DTO пользователя
migrations/
  1_create_base_entities.up.sql
  1_create_base_entities.down.sql
```

## Миграции

- up: `migrations/1_create_base_entities.up.sql`

  ```sql
  CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE,
    telegram_id BIGINT UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );
  ```

- down: `migrations/1_create_base_entities.down.sql`

  ```sql
  DROP TABLE IF EXISTS users;
  ```

## Запуск

- *Клонировать репозиторий и перейти в папку:*

  ```bash
  git clone https://github.com/hell-ecosystem/user-service.git
  cd user-service
  ```

- *Установить зависимости:*

  ```bash
  go mod download
  ```

- *Задать переменные окружения:*

  ```bash
  export APP_PORT=":8081"
  export DB_HOST="localhost"
  export DB_PORT="5432"
  export DB_USER="postgres"
  export DB_PASS="password"
  export DB_NAME="userdb"
  ```

- *Запустить миграции:*

  ```bash
  go run cmd/service/main.go migrate
  ```

- *Запустить сервис:*

  ```bash
  go run cmd/service/main.go serve
  ```

## Переменные окружения

| Переменная             | Описание                        | По умолчанию |
| ---------------------- | ------------------------------- | ------------ |
| `APP_PORT`             | Порт HTTP-сервера               | `:8081`      |
| `DB_HOST`              | Хост PostgreSQL                 | —            |
| `DB_PORT`              | Порт PostgreSQL                 | `5432`       |
| `DB_USER`              | Пользователь PostgreSQL         | —            |
| `DB_PASS`              | Пароль PostgreSQL               | —            |
| `DB_NAME`              | База данных PostgreSQL          | —            |
| `DB_SSLMODE`           | SSL-режим (`disable`,`require`) | `disable`    |
| `APP_READ_TIMEOUT`     | Read timeout, сек               | `10`         |
| `APP_WRITE_TIMEOUT`    | Write timeout, сек              | `10`         |
| `APP_IDLE_TIMEOUT`     | Idle timeout, сек               | `120`        |
| `DB_MAX_OPEN_CONNS`    | Max open connections            | `100`        |
| `DB_MAX_IDLE_CONNS`    | Max idle connections            | `20`         |
| `DB_CONN_MAX_LIFETIME` | Conn max lifetime, сек          | `3600`       |
