package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/go-playground/validator/v10"

	"github.com/hell-ecosystem/user-service/internal/retry"
)

type Config struct {
	AppPort string `env:"APP_PORT" envDefault:":8080" validate:"required"`

	DBHost    string `env:"DB_HOST" validate:"required,hostname|ip"`
	DBPort    string `env:"DB_PORT" envDefault:"5432" validate:"required,numeric"`
	DBUser    string `env:"DB_USER" validate:"required"`
	DBPass    string `env:"DB_PASS" validate:"required"`
	DBName    string `env:"DB_NAME" validate:"required"`
	DBSSLMode string `env:"DB_SSLMODE" envDefault:"disable" validate:"required,oneof=disable require"`

	// timeouts in seconds
	ReadTimeoutSec  int `env:"APP_READ_TIMEOUT" envDefault:"10" validate:"gte=1"`
	WriteTimeoutSec int `env:"APP_WRITE_TIMEOUT" envDefault:"10" validate:"gte=1"`
	IdleTimeoutSec  int `env:"APP_IDLE_TIMEOUT" envDefault:"120" validate:"gte=10"`

	DBMaxOpenConns       int `env:"DB_MAX_OPEN_CONNS" envDefault:"100"`
	DBMaxIdleConns       int `env:"DB_MAX_IDLE_CONNS" envDefault:"20"`
	DBConnMaxLifetimeSec int `env:"DB_CONN_MAX_LIFETIME" envDefault:"3600"`
}

var (
	Conf     Config
	validate *validator.Validate
)

// ретраи
var (
	// для работы с БД
	DBRetry = retry.New(
		retry.WithMaxAttempts(5),
		retry.WithBackoffExponential(100*time.Millisecond, 2.0),
		retry.WithJitter(0.1),
		retry.RetryIf(retry.IsTransientSQLError),
	)

	// для внешних HTTP-клиентов
	APIRetry = retry.New(
		retry.WithMaxAttempts(3),
		retry.WithBackoffExponential(200*time.Millisecond, 1.5),
		retry.RetryIf(retry.Is5xxHTTPError),
	)
)

func init() {
	validate = validator.New()
}

func Load() (*Config, error) {
	if err := env.Parse(&Conf); err != nil {
		return nil, fmt.Errorf("env parse: %w", err)
	}
	if err := validate.Struct(Conf); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}
	return &Conf, nil
}
