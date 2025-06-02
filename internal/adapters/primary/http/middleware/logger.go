package middleware

import (
	"net/http"
	"time"

	"github.com/2rprbm/conta-med-backend/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

// Logger is a middleware that logs HTTP requests
func Logger(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom response writer to capture the status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Process the request
			next.ServeHTTP(ww, r)

			// Log the request details
			log.Info(
				"%s %s %s %d %s %s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				ww.Status(),
				time.Since(start),
				r.UserAgent(),
			)
		})
	}
}
