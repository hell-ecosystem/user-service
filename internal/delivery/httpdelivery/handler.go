// internal/delivery/httpdelivery/handler.go
package httpdelivery

import (
	"errors"
	"net/http"
	"strings"

	middleware "github.com/hell-ecosystem/user-service/internal/delivery/httpdelivery/middleware"
	"github.com/hell-ecosystem/user-service/internal/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) http.Handler {
	h := &Handler{svc}

	mux := http.NewServeMux()
	mux.Handle("/users/me", http.HandlerFunc(h.GetMe))
	mux.Handle("/users/", http.HandlerFunc(h.GetByID))

	return middleware.Chain(
		mux,
		middleware.RecoveryMiddleware,
		middleware.RequestIDMiddleware,
		middleware.LoggingMiddleware,
		middleware.CORSMiddleware,
	)
}

// GetMe возвращает профиль текущего пользователя.
// Требует, чтобы middleware уже положил в контекст middleware.UserIDKey.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || uid == "" {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "неавторизованный запрос")
		return
	}

	u, err := h.svc.GetByID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "USER_NOT_FOUND", "пользователь не найден")
		} else {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "внутренняя ошибка сервера")
		}
		return
	}

	WriteSuccess(w, u)
}

// GetByID возвращает профиль пользователя по ID из пути /users/{id}.
// Здесь мы доверяем, что авторизация (X-User-ID → контекст) уже сделана,
// но по сути GetByID доступен любому аутентифицированному клиенту.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	// вы можете при желании проверить, что r.Context().Value(middleware.UserIDKey) не пуст,
	// но если роутинг навешан через тот же RequireUserIDMiddleware — оно уже гарантировано.

	id := strings.TrimPrefix(r.URL.Path, "/users/")

	u, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "USER_NOT_FOUND", "пользователь не найден")
		} else {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "внутренняя ошибка сервера")
		}
		return
	}

	WriteSuccess(w, u)
}
