package middleware

import (
	"net/http"
)

// Chain chains two `http.Handler` together, returning a single `http.Handler`
func Chain(before, after http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		before.ServeHTTP(w, r)
		after.ServeHTTP(w, r)
	})
}
