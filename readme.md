# User Service

Микросервис для управления пользователями.  
Позволяет регистрировать пользователей, аутентифицироваться через email/password или Telegram, и получать Access Token для дальнейшей работы с системой.

---

## 📦 Технологии

- **Golang 1.22+**
- **PostgreSQL** — хранение пользователей
- **Redis** — хранение refresh-токенов через `auth-service`
- **JWT** — генерация Access токенов
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

```json
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

```json
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

## 🚀 Планы на будущее

- Добавить refresh токены для продления сессий
- Реализовать /me для получения профиля пользователя
- Расширить роли (admin, moderator)
- Вынести конфиги Redis/JWT в централизованный config
- Добавить поддержку OAuth (Google, GitHub)
- Подключить трейсинг и метрики (Prometheus, OpenTelemetry)
