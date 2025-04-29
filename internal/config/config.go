package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	AppPort string `env:"APP_PORT" envDefault:":8080" validate:"required"`

	DBHost    string `env:"DB_HOST" validate:"required,hostname|ip"`
	DBPort    string `env:"DB_PORT" envDefault:"5432" validate:"required,numeric"`
	DBUser    string `env:"DB_USER" validate:"required"`
	DBPass    string `env:"DB_PASS" validate:"required"`
	DBName    string `env:"DB_NAME" validate:"required"`
	DBSSLMode string `env:"DB_SSLMODE" envDefault:"disable" validate:"required,oneof=disable require"`

	ReadTimeout  int `env:"APP_READ_TIMEOUT" envDefault:"10" validate:"required,gte=1"`
	WriteTimeout int `env:"APP_WRITE_TIMEOUT" envDefault:"10" validate:"required,gte=1"`
	IdleTimeout  int `env:"APP_IDLE_TIMEOUT" envDefault:"120" validate:"required,gte=10"`

	JWTSecret string `env:"JWT_SECRET" validate:"required"`
	RedisAddr string `env:"REDIS_ADDR" envDefault:"localhost:6379"`

	SentryDSN        string  `env:"SENTRY_DSN"`
	SentryEnv        string  `env:"SENTRY_ENV" envDefault:"development"`
	SentrySampleRate float64 `env:"SENTRY_SAMPLE_RATE" envDefault:"1.0" validate:"gte=0,lte=1"`

	OtelExporterEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"localhost:4318" validate:"required"`

	ServiceName string `env:"SERVICE_NAME" envDefault:"user-service"`
}

var Conf Config
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Load() (*Config, error) {
	if err := env.Parse(&Conf); err != nil {
		return nil, fmt.Errorf("failed to load env vars: %w", err)
	}
	if err := validate.Struct(Conf); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}
	return &Conf, nil
}
