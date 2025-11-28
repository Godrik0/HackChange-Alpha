package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
)

func Recovery(logger interfaces.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						"error", err,
						"stack", string(debug.Stack()),
						"path", r.URL.Path,
						"method", r.Method,
					)

					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
