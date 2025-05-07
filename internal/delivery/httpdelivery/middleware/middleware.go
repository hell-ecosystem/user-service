package middleware

import "net/http"

// Middleware — функция, которая берёт http.Handler и возвращает обёрнутый.
type Middleware func(http.Handler) http.Handler

// Chain применяет список middleware к handler в порядке вызова.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}
