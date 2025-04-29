# User Service

Микросервис для управления пользователями.  
Позволяет регистрировать пользователей, аутентифицироваться через email/password или Telegram, и получать Access Token для дальнейшей работы с системой.

---

## 📦 Технологии

- **Golang 1.22+**
- **PostgreSQL** — хранение пользователей
- **Redis** — хранение refresh-токенов через `auth-service`
- **JWT** — генерация Access токенов
- **Чистая архитектура (Clean Architecture)** — чёткое разделение слоёв
- **Библиотека авторизации**: [`hell-ecosystem/auth-service`](https://github.com/hell-ecosystem/auth-service)

---

## 📚 Архитектура слоёв

```plaintext
cmd/main.go                  — запуск приложения
internal/config/             — загрузка конфигурации
internal/delivery/http/      — http-обработчики
internal/service/            — бизнес-логика пользователей
internal/repository/postgres/— работа с PostgreSQL
internal/model/              — модели DTO
migrations/                  — миграции БД
```

## 🔑 Форматы запросов
### `POST /register`
Описание: Регистрация нового пользователя.

**Request**:
```
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**:
`<access_token>`

### `POST /login`
Описание: Аутентификация по email и паролю.

**Request**:
```
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**:
`<access_token>`

### `POST /telegram?id={telegram_id}`
Описание: Аутентификация или создание нового пользователя через Telegram.

**Request**:

В query параметре: `id=123456789`

**Response**:
`<access_token>`
