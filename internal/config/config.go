package config

import (
	"log"
	"os"
)

type Config struct {
	HTTPPort  string
	DBURL     string
	JWTSecret string
}

func Load() *Config {
	cfg := &Config{
		HTTPPort:  ":8080",
		DBURL:     os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
	if cfg.DBURL == "" || cfg.JWTSecret == "" {
		log.Fatal("Missing required env vars")
	}
	return cfg
}
