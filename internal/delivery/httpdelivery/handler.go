package httpdelivery

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	authsvc "github.com/hell-ecosystem/auth-service/pkg/auth/service"
	"github.com/hell-ecosystem/user-service/internal/service"
)

type Handler struct {
	svc  *service.Service
	auth *authsvc.AuthService
}

func NewHandler(svc *service.Service, auth *authsvc.AuthService) *Handler {
	return &Handler{svc: svc, auth: auth}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", h.withPost(h.Register))
	mux.HandleFunc("/login", h.withPost(h.Login))
	mux.HandleFunc("/telegram", h.withPost(h.TelegramLogin))
	return mux
}

func (h *Handler) withPost(f func(ctx context.Context, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		f(r.Context(), w, r)
	}
}

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID, err := h.svc.RegisterUser(ctx, c.Email, c.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	accessToken, err := h.auth.Login(ctx, userID, "user")
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(accessToken))
}

func (h *Handler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID, err := h.svc.AuthenticateUser(ctx, c.Email, c.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	accessToken, err := h.auth.Login(ctx, userID, "user")
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(accessToken))
}

func (h *Handler) TelegramLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tgIDRaw := r.URL.Query().Get("id")
	tgID, err := strconv.ParseInt(tgIDRaw, 10, 64)
	if err != nil {
		http.Error(w, "invalid telegram id", http.StatusBadRequest)
		return
	}
	userID, err := h.svc.AuthenticateTelegramUser(ctx, tgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	accessToken, err := h.auth.Login(ctx, userID, "user")
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(accessToken))
}
