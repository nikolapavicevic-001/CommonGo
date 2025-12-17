package httpx

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"github.com/nikolapavicevic-001/CommonGo/logger"
)

// RequestLogger returns a middleware that logs HTTP requests using zerolog.
// It logs method, path, status, duration, and correlates with request_id.
func RequestLogger(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Get or generate request ID
			requestID := middleware.GetReqID(r.Context())
			if requestID == "" {
				requestID = "unknown"
			}

			// Create a logger with request context
			reqLog := log.With().
				Str("request_id", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Logger()

			// Attach logger to context
			ctx := logger.With(r.Context(), reqLog)
			r = r.WithContext(ctx)

			// Wrap response writer to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Process request
			next.ServeHTTP(ww, r)

			// Log completion
			duration := time.Since(start)
			status := ww.Status()

			event := reqLog.Info()
			if status >= 500 {
				event = reqLog.Error()
			} else if status >= 400 {
				event = reqLog.Warn()
			}

			event.
				Int("status", status).
				Int("bytes", ww.BytesWritten()).
				Dur("duration", duration).
				Msg("request completed")
		})
	}
}

// RequestLoggerWithOptions returns a request logger with customizable options.
type RequestLoggerOptions struct {
	// SkipPaths are paths that should not be logged (e.g., /health, /metrics)
	SkipPaths []string

	// LogRequestBody logs the request body (use with caution for large payloads)
	LogRequestBody bool

	// LogResponseBody logs the response body (use with caution)
	LogResponseBody bool
}

// RequestLoggerWithOpts returns a middleware with custom options.
func RequestLoggerWithOpts(log zerolog.Logger, opts RequestLoggerOptions) func(http.Handler) http.Handler {
	skipMap := make(map[string]bool)
	for _, p := range opts.SkipPaths {
		skipMap[p] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip logging for specified paths
			if skipMap[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			requestID := middleware.GetReqID(r.Context())
			if requestID == "" {
				requestID = "unknown"
			}

			reqLog := log.With().
				Str("request_id", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Logger()

			ctx := logger.With(r.Context(), reqLog)
			r = r.WithContext(ctx)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			status := ww.Status()

			event := reqLog.Info()
			if status >= 500 {
				event = reqLog.Error()
			} else if status >= 400 {
				event = reqLog.Warn()
			}

			event.
				Int("status", status).
				Int("bytes", ww.BytesWritten()).
				Dur("duration", duration).
				Msg("request completed")
		})
	}
}

