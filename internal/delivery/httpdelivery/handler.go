package httpdelivery

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/hell-ecosystem/user-service/internal/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/users/me", h.withAuth(h.GetMe))
	mux.Handle("/users/", h.withAuth(h.GetByID))
	return mux
}

// withAuth проверяет метод и передаёт в f контекст с userID.
// В продакшене здесь же вешается middleware для валидации JWT.
func (h *Handler) withAuth(
	f func(ctx context.Context, w http.ResponseWriter, r *http.Request),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPut {
			WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "метод не разрешён")
			return
		}
		// В реальной логике из JWT-плагина gateway в контекст попадает "userID"
		f(r.Context(), w, r)
	}
}

// GetMe возвращает профиль текущего пользователя
func (h *Handler) GetMe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "неавторизованный запрос")
		return
	}

	u, err := h.svc.GetByID(ctx, userID)
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

// GetByID возвращает профиль пользователя по ID из пути /users/{id}
func (h *Handler) GetByID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	u, err := h.svc.GetByID(ctx, id)
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
