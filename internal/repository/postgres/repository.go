package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hell-ecosystem/user-service/internal/model"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *model.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, email, telegram_id, created_at)
		VALUES ($1, $2, $3, $4)
	`, u.ID, u.Email, u.TelegramID, u.CreatedAt)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, telegram_id, created_at
		FROM users WHERE id = $1
	`, id)
	var u model.User
	if err := row.Scan(&u.ID, &u.Email, &u.TelegramID, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
