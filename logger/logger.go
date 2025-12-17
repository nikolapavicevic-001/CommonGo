// Package logger provides zerolog-based logging utilities with context propagation.
package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type ctxKey struct{}

// New creates a new zerolog.Logger with the specified level and service name.
// The logger outputs to os.Stdout with pretty console formatting in development
// or JSON formatting based on the LOG_FORMAT environment variable.
func New(level string, serviceName string) zerolog.Logger {
	var output io.Writer = os.Stdout

	// Use console writer for human-readable output if LOG_FORMAT != "json"
	if os.Getenv("LOG_FORMAT") != "json" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	lvl := parseLevel(level)

	return zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()
}

// With attaches the logger to the context for later retrieval.
func With(ctx context.Context, log zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

// From retrieves the logger from the context.
// If no logger is found, it returns a disabled logger that produces no output.
func From(ctx context.Context) zerolog.Logger {
	if log, ok := ctx.Value(ctxKey{}).(zerolog.Logger); ok {
		return log
	}
	return zerolog.Nop()
}

// WithFields returns a new context with additional fields added to the logger.
func WithFields(ctx context.Context, fields map[string]interface{}) context.Context {
	log := From(ctx)
	logCtx := log.With()
	for k, v := range fields {
		logCtx = logCtx.Interface(k, v)
	}
	return With(ctx, logCtx.Logger())
}

// WithRequestID adds a request_id field to the logger in context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	log := From(ctx).With().Str("request_id", requestID).Logger()
	return With(ctx, log)
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

