package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/hell-ecosystem/user-service/internal/config"
	"github.com/hell-ecosystem/user-service/internal/retry"
	_ "github.com/lib/pq"
)

// Connect открывает соединение с Postgres и пингует его,
// пока БД не станет доступна (с экспоненциальным backoff + jitter).
func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.DatabaseDSN()
	var db *sql.DB

	// Ретраер для старта БД: бесконечный → ждём, пока контейнер не заведётся
	r := retry.New(
		retry.WithMaxAttempts(0), // 0 = бесконечно
		retry.WithBackoffExponential(200*time.Millisecond, 1.5),
		retry.WithJitter(0.1),
		retry.RetryIf(retry.IsTransientSQLError), // ретраим только “транзиентные” SQL-ошибки
	)

	err := r.Do(context.Background(), func() error {
		var err error
		// открываем соединение (только один раз)
		if db == nil {
			db, err = sql.Open("postgres", dsn)
			if err != nil {
				return err
			}
			db.SetMaxOpenConns(cfg.DBMaxOpenConns)
			db.SetMaxIdleConns(cfg.DBMaxIdleConns)
			db.SetConnMaxLifetime(cfg.GetConnMaxLifetime())
		}
		// пингуем — если БД ещё не готова, вернётся ошибка и retry сработает
		return db.PingContext(context.Background())
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
