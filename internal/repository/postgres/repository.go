package postgres

import (
	"context"
	"database/sql"

	"github.com/hell-ecosystem/user-service/internal/model"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, u *model.User) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, email, password_hash, telegram_id, created_at) VALUES ($1, $2, $3, $4, NOW())`, u.ID, u.Email, u.Password, u.TelegramID)
	return err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, telegram_id, created_at FROM users WHERE email = $1`, email)
	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.TelegramID, &u.CreatedAt)
	return &u, err
}

func (r *Repository) GetByTelegramID(ctx context.Context, tgID int64) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, telegram_id, created_at FROM users WHERE telegram_id = $1`, tgID)
	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.TelegramID, &u.CreatedAt)
	return &u, err
}
