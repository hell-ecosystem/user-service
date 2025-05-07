// internal/repository/postgres/repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hell-ecosystem/user-service/internal/config"
	"github.com/hell-ecosystem/user-service/internal/model"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *sql.DB
}

// New создаёт репозиторий поверх sql.DB
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create вставляет нового пользователя с автоматическим retry на transient-ошибки
func (r *Repository) Create(ctx context.Context, u *model.User) error {
	return config.DBRetry.Do(ctx, func() error {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO users (id, email, telegram_id, created_at)
			VALUES ($1, $2, $3, $4)
		`, u.ID, u.Email, u.TelegramID, u.CreatedAt)
		return err
	})
}

// GetByID получает пользователя по ID с retry и обрабатывает sql.ErrNoRows
func (r *Repository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User

	err := config.DBRetry.Do(ctx, func() error {
		row := r.db.QueryRowContext(ctx, `
			SELECT id, email, telegram_id, created_at
			FROM users WHERE id = $1
		`, id)
		return row.Scan(&u.ID, &u.Email, &u.TelegramID, &u.CreatedAt)
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
