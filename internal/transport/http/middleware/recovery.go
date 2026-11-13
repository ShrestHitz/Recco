package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery catches panics in the request lifecycle and returns an HTTP 500 cleanly.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := string(debug.Stack())
				slog.Error("server panic recovered",
					"panic", rec,
					"path", r.URL.Path,
					"stack", stack,
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "internal server error",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
