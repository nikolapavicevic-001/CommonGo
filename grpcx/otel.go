package grpcx

import (
	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// OTELServerOptions returns gRPC server options to enable OpenTelemetry instrumentation.
//
// Note: You still need to configure a global OTel provider/exporter in your service.
func OTELServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}
}


