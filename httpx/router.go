// Package httpx provides chi router utilities, middleware, and response helpers.
package httpx

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// RouterOption is a function that configures a chi.Mux router.
type RouterOption func(*chi.Mux)

// NewRouter creates a new chi.Mux with standard middlewares applied.
// Default middlewares: RequestID, RealIP, Recoverer.
// Use options to customize behavior (CORS, timeouts, etc.).
func NewRouter(opts ...RouterOption) *chi.Mux {
	r := chi.NewRouter()

	// Default middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Apply custom options
	for _, opt := range opts {
		opt(r)
	}

	return r
}

// WithTimeout adds a timeout middleware with the specified duration.
func WithTimeout(timeout time.Duration) RouterOption {
	return func(r *chi.Mux) {
		r.Use(middleware.Timeout(timeout))
	}
}

// WithCORS adds CORS middleware with the specified configuration.
func WithCORS(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) RouterOption {
	return func(r *chi.Mux) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   allowedOrigins,
			AllowedMethods:   allowedMethods,
			AllowedHeaders:   allowedHeaders,
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}
}

// WithCORSDefaults adds CORS middleware with permissive defaults.
// Allows all origins, common methods, and standard headers.
func WithCORSDefaults() RouterOption {
	return WithCORS(
		[]string{"*"},
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		[]string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
	)
}

// WithMiddleware adds custom middleware to the router.
func WithMiddleware(middlewares ...func(http.Handler) http.Handler) RouterOption {
	return func(r *chi.Mux) {
		for _, m := range middlewares {
			r.Use(m)
		}
	}
}

// WithHeartbeat adds a heartbeat endpoint at the specified path.
func WithHeartbeat(path string) RouterOption {
	return func(r *chi.Mux) {
		r.Use(middleware.Heartbeat(path))
	}
}

// WithStripSlashes adds middleware to strip trailing slashes from URLs.
func WithStripSlashes() RouterOption {
	return func(r *chi.Mux) {
		r.Use(middleware.StripSlashes)
	}
}

// WithCompression adds response compression middleware.
func WithCompression(level int) RouterOption {
	return func(r *chi.Mux) {
		r.Use(middleware.Compress(level))
	}
}

