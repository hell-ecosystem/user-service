package httpdelivery

import (
	"errors"
	"net/http"

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
	mux.Handle("/users", http.HandlerFunc(h.GetByID))

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

// GetByID возвращает профиль пользователя по ID из заголовка X-User-ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-User-ID")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "MISSING_USER_ID", "заголовок X-User-ID не указан")
		return
	}

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
