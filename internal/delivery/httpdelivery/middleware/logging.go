package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		// 1) Собираем Attr’ы
		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.status),
			slog.Duration("latency", time.Since(start)),
		}
		if rid := GetRequestID(r.Context()); rid != "" {
			attrs = append(attrs, slog.String("request_id", rid))
		}

		// 2) Конвертируем []slog.Attr → []any
		args := make([]any, len(attrs))
		for i, a := range attrs {
			args[i] = a
		}

		// 3) Вызываем Info
		slog.Info("http request", args...)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
