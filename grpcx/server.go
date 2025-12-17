package grpcx

import (
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Options configures the gRPC server defaults provided by CommonGo.
type Options struct {
	// Logger is used by logging interceptors. If unset, logging interceptors are disabled.
	Logger zerolog.Logger

	// EnableHealth registers the standard gRPC health service.
	EnableHealth bool

	// EnableReflection enables gRPC server reflection.
	EnableReflection bool

	// EnableOTel enables OpenTelemetry gRPC instrumentation (stats handler).
	EnableOTel bool
}

// NewServer constructs a *grpc.Server with standard CommonGo interceptors and optional features enabled.
//
// extra options are appended after CommonGo's options, so callers can override as needed.
func NewServer(opts Options, extra ...grpc.ServerOption) (*grpc.Server, error) {
	var serverOpts []grpc.ServerOption

	// Interceptors
	// Note: we chain unary/stream interceptors so services can still add their own via extra opts.
	if reflect.ValueOf(opts.Logger).IsZero() {
		return nil, fmt.Errorf("creating grpc server: Options.Logger must be set (use logger.New(...) or zerolog.Nop())")
	}
	serverOpts = append(serverOpts,
		grpc.ChainUnaryInterceptor(UnaryLoggingInterceptor(opts.Logger)),
		grpc.ChainStreamInterceptor(StreamLoggingInterceptor(opts.Logger)),
	)

	// OpenTelemetry
	if opts.EnableOTel {
		serverOpts = append(serverOpts, OTELServerOptions()...)
	}

	serverOpts = append(serverOpts, extra...)

	s := grpc.NewServer(serverOpts...)

	if opts.EnableHealth {
		RegisterHealth(s)
	}
	if opts.EnableReflection {
		RegisterReflection(s)
	}

	// Basic sanity check: reflection without health is fine; no further validation required today.
	if s == nil {
		return nil, fmt.Errorf("creating grpc server: got nil")
	}
	return s, nil
}


