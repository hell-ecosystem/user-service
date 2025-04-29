package service

import (
	"context"
	"errors"

	"github.com/hell-ecosystem/user-service/internal/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(context.Context, *model.User) error
	GetByEmail(context.Context, string) (*model.User, error)
	GetByTelegramID(context.Context, int64) (*model.User, error)
}

type Service struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterUser(ctx context.Context, email, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Email:    &email,
		Password: ptr(string(hash)),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user.Password == nil {
		return "", errors.New("password not set")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	return user.ID, nil
}

func (s *Service) AuthenticateTelegramUser(ctx context.Context, tgID int64) (string, error) {
	user, err := s.repo.GetByTelegramID(ctx, tgID)
	if err == nil {
		return user.ID, nil
	}

	// Если пользователь не найден, создаём нового
	newUser := &model.User{
		ID:         uuid.New().String(),
		TelegramID: &tgID,
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return "", err
	}

	return newUser.ID, nil
}

func ptr[T any](v T) *T {
	return &v
}
