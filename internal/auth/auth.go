package auth

import (
	"time"

	authinfra "github.com/hell-ecosystem/auth-service/pkg/auth/infra"
	authsvc "github.com/hell-ecosystem/auth-service/pkg/auth/service"
	"github.com/hell-ecosystem/user-service/internal/config"
)

func InitAuth(cfg *config.Config) *authsvc.AuthService {
	redis := authinfra.NewRedisTokenStore(cfg.AuthRedisAddr)
	jwt := authinfra.NewJWTManager(cfg.AuthJWTSecret)
	return authsvc.NewAuthService(jwt, redis, 15*time.Minute)
}
