package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hell-ecosystem/user-service/internal/config"
	_ "github.com/lib/pq"
)

// Connect открывает sql.DB, настраивает пул, а затем пингует в retry-петле.
// Если по таймауту контекста (здесь 30s) всё ещё нет коннекта — возвращаем ошибку.
func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.DatabaseDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// настраиваем пул
	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.GetConnMaxLifetime())

	// 30-секундный контекст на все попытки
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := config.DBConnectRetry.Do(ctx, func() error {
		return db.PingContext(ctx)
	}); err != nil {
		// если не смогли ни разу запинговать — закрываем и возвращаем
		_ = db.Close()
		return nil, fmt.Errorf("db ping after retries: %w", err)
	}

	return db, nil
}
