// internal/db/connect.go
package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/hell-ecosystem/user-service/internal/config"
	_ "github.com/lib/pq"
)

// Connect открывает sql.DB, настраивает его и пингует с retry.
// Если после всех попыток соединение не установилось, возвращает ошибку.
func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.DatabaseDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// Настраиваем пул
	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.GetConnMaxLifetime())

	// Пингуем с retry
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := config.DBRetry.Do(ctx, func() error {
		return db.PingContext(ctx)
	}); err != nil {
		// если не удалось соединиться — закрываем и возвращаем ошибку
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
