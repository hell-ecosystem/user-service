package middleware

import (
	"context"
	"net/http"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

// RequireUserIDMiddleware aborts with 401 if no userID is in ctx (or header).
func RequireUserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// допустим, gateway положил ID в заголовок X-User-ID
		uid := r.Header.Get("X-User-ID")
		if uid == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// или, если gateway уже положил его в контекст:
		// uid, _ := r.Context().Value(UserIDKey).(string)
		ctx := context.WithValue(r.Context(), UserIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
