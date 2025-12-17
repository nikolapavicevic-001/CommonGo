package grpcx

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// RegisterHealth registers the standard gRPC health service and sets it to SERVING.
func RegisterHealth(server *grpc.Server) *health.Server {
	hs := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, hs)
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	return hs
}

// RegisterReflection enables server reflection (useful for grpcurl, local debugging).
func RegisterReflection(server *grpc.Server) {
	reflection.Register(server)
}


