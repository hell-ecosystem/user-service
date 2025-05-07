package service

import (
	"context"
	"errors"

	"github.com/hell-ecosystem/user-service/internal/model"
	"github.com/hell-ecosystem/user-service/internal/repository/postgres"
)

var ErrNotFound = errors.New("user not found")

type Repository interface {
	Create(ctx context.Context, u *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetByID возвращает пользователя по ID
func (s *Service) GetByID(ctx context.Context, id string) (*model.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// если репозиторий вернул ErrNotFound — пробрасываем сервисную ErrNotFound
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, ErrNotFound
		}
		// иначе — возвращаем любую другую ошибку как есть
		return nil, err
	}
	return u, nil
}
