// Package grpcx provides gRPC server helpers (interceptors, health, reflection, OTEL).
package grpcx

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const requestIDHeader = "x-request-id"

// UnaryLoggingInterceptor logs unary RPCs using zerolog.
func UnaryLoggingInterceptor(log zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		code := status.Code(err)
		ev := eventForCode(log, code)

		ev.
			Str("grpc_method", info.FullMethod).
			Str("grpc_code", code.String()).
			Dur("duration", time.Since(start)).
			Str("request_id", requestIDFromIncomingContext(ctx)).
			Str("peer_ip", peerIP(ctx)).
			Msg("grpc request")

		return resp, err
	}
}

// StreamLoggingInterceptor logs stream RPCs using zerolog.
func StreamLoggingInterceptor(log zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)

		code := status.Code(err)
		ev := eventForCode(log, code)

		ev.
			Str("grpc_method", info.FullMethod).
			Bool("grpc_is_client_stream", info.IsClientStream).
			Bool("grpc_is_server_stream", info.IsServerStream).
			Str("grpc_code", code.String()).
			Dur("duration", time.Since(start)).
			Str("request_id", requestIDFromIncomingContext(ss.Context())).
			Str("peer_ip", peerIP(ss.Context())).
			Msg("grpc request")

		return err
	}
}

func eventForCode(log zerolog.Logger, code codes.Code) *zerolog.Event {
	switch {
	case code == codes.OK:
		return log.Info()
	case code == codes.Canceled || code == codes.DeadlineExceeded:
		return log.Warn()
	case code == codes.InvalidArgument || code == codes.NotFound || code == codes.AlreadyExists ||
		code == codes.PermissionDenied || code == codes.Unauthenticated || code == codes.FailedPrecondition ||
		code == codes.ResourceExhausted || code == codes.Aborted || code == codes.OutOfRange:
		return log.Warn()
	default:
		return log.Error()
	}
}

func requestIDFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	vals := md.Get(requestIDHeader)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func peerIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok || p.Addr == nil {
		return ""
	}
	// p.Addr.String() includes port for TCP; keep it since it's useful for debugging.
	return strings.TrimSpace(p.Addr.String())
}


