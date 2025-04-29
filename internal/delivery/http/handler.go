package httpdelivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	authinfra "github.com/hell-ecosystem/auth-service/pkg/auth/infra"
	authsvc "github.com/hell-ecosystem/auth-service/pkg/auth/service"
	"github.com/hell-ecosystem/user-service/internal/config"
	"github.com/hell-ecosystem/user-service/internal/service"
)

type Handler struct {
	svc  *service.Service
	auth *authsvc.AuthService
}

func NewHandler(cfg *config.Config, svc *service.Service) *Handler {
	redis := authinfra.NewRedisTokenStore(cfg.AuthRedisAddr)
	jwt := authinfra.NewJWTManager(cfg.AuthJWTSecret)
	auth := authsvc.NewAuthService(jwt, redis, 15*time.Minute)
	return &Handler{svc: svc, auth: auth}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", h.wrap(h.Register))
	mux.HandleFunc("/login", h.wrap(h.Login))
	mux.HandleFunc("/telegram", h.wrap(h.TelegramLogin))
	return mux
}

func (h *Handler) wrap(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		f(w, r)
	}
}

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID, err := h.svc.RegisterUser(r.Context(), c.Email, c.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	accessToken, err := h.auth.Login(r.Context(), userID, "user")
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(accessToken))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID, err := h.svc.AuthenticateUser(r.Context(), c.Email, c.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	accessToken, err := h.auth.Login(r.Context(), userID, "user")
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(accessToken))
}

func (h *Handler) TelegramLogin(w http.ResponseWriter, r *http.Request) {
	tgIDRaw := r.URL.Query().Get("id")
	tgID, err := strconv.ParseInt(tgIDRaw, 10, 64)
	if err != nil {
		http.Error(w, "invalid telegram id", http.StatusBadRequest)
		return
	}
	userID, err := h.svc.AuthenticateTelegramUser(r.Context(), tgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	accessToken, err := h.auth.Login(r.Context(), userID, "user")
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(accessToken))
}
