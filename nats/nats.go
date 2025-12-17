// Package nats provides NATS connection helpers and configuration.
package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// Config holds NATS connection configuration.
type Config struct {
	// URL is the NATS server URL (e.g., nats://localhost:4222)
	URL string

	// Name is the client connection name for identification
	Name string

	// ReconnectWait is the time to wait between reconnect attempts (default: 2s)
	ReconnectWait time.Duration

	// MaxReconnects is the maximum number of reconnect attempts (-1 for infinite, default: 60)
	MaxReconnects int

	// Timeout is the connection timeout (default: 5s)
	Timeout time.Duration

	// PingInterval is the interval for PING/PONG health checks (default: 2m)
	PingInterval time.Duration

	// MaxPingsOut is the max pending pings before connection is considered stale (default: 2)
	MaxPingsOut int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(url, name string) Config {
	return Config{
		URL:           url,
		Name:          name,
		ReconnectWait: 2 * time.Second,
		MaxReconnects: 60,
		Timeout:       5 * time.Second,
		PingInterval:  2 * time.Minute,
		MaxPingsOut:   2,
	}
}

// Connect establishes a connection to a NATS server using the provided configuration.
func Connect(cfg Config) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(cfg.Name),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.Timeout(cfg.Timeout),
		nats.PingInterval(cfg.PingInterval),
		nats.MaxPingsOutstanding(cfg.MaxPingsOut),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				// Connection lost - will attempt reconnect
			}
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			// Reconnected successfully
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			// Connection closed
		}),
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connecting to NATS: %w", err)
	}

	return nc, nil
}

// MustConnect is like Connect but panics on error.
func MustConnect(cfg Config) *nats.Conn {
	nc, err := Connect(cfg)
	if err != nil {
		panic(err)
	}
	return nc
}

// ConnectWithHandlers establishes a connection with custom disconnect/reconnect handlers.
func ConnectWithHandlers(
	cfg Config,
	onDisconnect func(*nats.Conn, error),
	onReconnect func(*nats.Conn),
	onClose func(*nats.Conn),
) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(cfg.Name),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.Timeout(cfg.Timeout),
		nats.PingInterval(cfg.PingInterval),
		nats.MaxPingsOutstanding(cfg.MaxPingsOut),
	}

	if onDisconnect != nil {
		opts = append(opts, nats.DisconnectErrHandler(onDisconnect))
	}
	if onReconnect != nil {
		opts = append(opts, nats.ReconnectHandler(onReconnect))
	}
	if onClose != nil {
		opts = append(opts, nats.ClosedHandler(onClose))
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connecting to NATS: %w", err)
	}

	return nc, nil
}

