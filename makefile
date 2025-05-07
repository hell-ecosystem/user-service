# -------------------------------------------------------------------
# Environment variables (can be overridden when calling make):
# -------------------------------------------------------------------
export APP_PORT               ?= :8081
export LOG_LEVEL              ?= INFO
export LOG_FORMAT             ?= text

export DB_HOST                ?= localhost
export DB_PORT                ?= 5432
export DB_USER                ?= postgres
export DB_PASS                ?= postgres
export DB_NAME                ?= postgres_user_service
export DB_SSLMODE             ?= disable

export APP_READ_TIMEOUT       ?= 10
export APP_WRITE_TIMEOUT      ?= 10
export APP_IDLE_TIMEOUT       ?= 120

export DB_MAX_OPEN_CONNS      ?= 100
export DB_MAX_IDLE_CONNS      ?= 20
export DB_CONN_MAX_LIFETIME   ?= 3600

# -------------------------------------------------------------------
# Targets
# -------------------------------------------------------------------
.PHONY: serve migrate

# Запустить HTTP-сервер
serve:
	@echo "Starting HTTP server on $$APP_PORT..."
	go run cmd/service/main.go serve

# Применить все миграции (директория migrations/, DSN из конфига)
migrate:
	@echo "Running migrations against postgres://$$DB_USER:$$DB_PASS@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE"
	go run cmd/service/main.go migrate
