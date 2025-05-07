# User Service

CRUD-сервис для работы с даанными пользователя.

## Эндпоинты

### GET /users/me

Получить профиль текущего аутентифицированного пользователя.

### GET /users/{id}

Получить профиль пользователя по ID.

---

## Архитектура

- `cmd/main.go` — старт приложения  
- `internal/config` — загрузка env- и валидация  
- `internal/delivery/httpdelivery` — HTTP-хендлеры  
- `internal/service` — бизнес-логика (GetByID)  
- `internal/repository/postgres` — взаимодействие с БД  
- `internal/model` — DTO  
- `migrations` — SQL-миграции  

---

## Поднимаем

```bash
export APP_PORT=":8081"
export DB_HOST="localhost"
export DB_USER="..."
export DB_PASS="..."
export DB_NAME="userdb"
# и др.
go run cmd/main.go
