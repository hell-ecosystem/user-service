package httpdelivery

import (
	"context"
	"encoding/json"
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

// в проде сюда вставляется middleware, который заполняет Context значением userID
func (h *Handler) withAuth(f func(ctx context.Context, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// ctx := context.WithValue(r.Context(), "userID", ...)
		f(r.Context(), w, r)
	}
}

func (h *Handler) GetMe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID := ctx.Value("userID").(string)
	u, err := h.svc.GetByID(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func (h *Handler) GetByID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Выдергиваем ID из URL
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	u, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
