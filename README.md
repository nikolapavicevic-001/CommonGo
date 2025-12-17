# CommonGo

Shared Go utilities for HomeLab microservices.

## Installation

```bash
go get github.com/nikolapavicevic-001/CommonGo
```

## Packages

### logger

Zerolog-based logging with context propagation.

```go
import "github.com/nikolapavicevic-001/CommonGo/logger"

// Create a logger
log := logger.New("info", "my-service")

// Attach to context
ctx = logger.With(ctx, log)

// Retrieve from context
log := logger.From(ctx)
log.Info().Msg("hello world")

// Add request ID to context logger
ctx = logger.WithRequestID(ctx, "req-123")
```

### config

Environment variable helpers.

```go
import "github.com/nikolapavicevic-001/CommonGo/config"

// Load common config
cfg := config.LoadCommon()
// cfg.ServiceName, cfg.LogLevel, cfg.Environment

// Individual helpers
port := config.GetEnvInt("PORT", 8080)
timeout := config.GetEnvDuration("TIMEOUT", 30*time.Second)
debug := config.GetEnvBool("DEBUG", false)
```

### postgres

PostgreSQL connection pool helpers using pgxpool.

```go
import "github.com/nikolapavicevic-001/CommonGo/postgres"

cfg := postgres.DefaultConfig("postgres://user:pass@localhost:5432/mydb")
cfg.MaxConns = 20

pool, err := postgres.Open(ctx, cfg)
if err != nil {
    log.Fatal().Err(err).Msg("failed to connect to postgres")
}
defer pool.Close()
```

### nats

NATS connection helpers.

```go
import "github.com/nikolapavicevic-001/CommonGo/nats"

cfg := nats.DefaultConfig("nats://localhost:4222", "my-service")

nc, err := nats.Connect(cfg)
if err != nil {
    log.Fatal().Err(err).Msg("failed to connect to NATS")
}
defer nc.Close()

// With custom handlers
nc, err := nats.ConnectWithHandlers(cfg,
    func(c *nats.Conn, err error) { log.Warn().Err(err).Msg("disconnected") },
    func(c *nats.Conn) { log.Info().Msg("reconnected") },
    func(c *nats.Conn) { log.Info().Msg("connection closed") },
)
```

### httpx

Chi router utilities, middleware, and JSON response helpers.

#### Router Factory

```go
import "github.com/nikolapavicevic-001/CommonGo/httpx"

// Create router with defaults (RequestID, RealIP, Recoverer)
r := httpx.NewRouter()

// With options
r := httpx.NewRouter(
    httpx.WithTimeout(30 * time.Second),
    httpx.WithCORSDefaults(),
    httpx.WithHeartbeat("/ping"),
    httpx.WithCompression(5),
)
```

#### Request Logger Middleware

```go
log := logger.New("info", "my-service")

r := httpx.NewRouter(
    httpx.WithMiddleware(httpx.RequestLogger(log)),
)

// With options (skip health checks, etc.)
r := httpx.NewRouter(
    httpx.WithMiddleware(httpx.RequestLoggerWithOpts(log, httpx.RequestLoggerOptions{
        SkipPaths: []string{"/health", "/metrics"},
    })),
)
```

#### JSON Response Helpers

```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    user := User{ID: "123", Name: "John"}
    
    // Write with standard envelope: {"data": {...}, "request_id": "..."}
    httpx.WriteData(w, r, http.StatusOK, user)
    
    // With metadata (pagination, etc.)
    httpx.WriteDataWithMeta(w, r, http.StatusOK, users, Meta{Total: 100})
    
    // Raw JSON (no envelope)
    httpx.WriteJSON(w, r, http.StatusOK, user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Error responses: {"error": {"code": "...", "message": "..."}, "request_id": "..."}
    httpx.WriteBadRequest(w, r, "invalid email format")
    httpx.WriteNotFound(w, r, "user not found")
    httpx.WriteUnauthorized(w, r, "invalid token")
    httpx.WriteInternalError(w, r, "database error")
    
    // Custom error
    httpx.WriteError(w, r, http.StatusTeapot, "teapot", "I'm a teapot")
    
    // No content
    httpx.NoContent(w)
}
```

### grpcx

gRPC server helpers: server builder, zerolog interceptors, OTEL instrumentation, and health/reflection helpers.

```go
import (
  "net"

  "github.com/nikolapavicevic-001/CommonGo/grpcx"
  "github.com/nikolapavicevic-001/CommonGo/logger"
)

log := logger.New("info", "device-service")

srv, err := grpcx.NewServer(grpcx.Options{
  Logger:           log,
  EnableHealth:     true,
  EnableReflection: true,
  EnableOTel:       true,
})
if err != nil {
  log.Fatal().Err(err).Msg("failed to create grpc server")
}

lis, err := net.Listen("tcp", ":9090")
if err != nil {
  log.Fatal().Err(err).Msg("failed to listen")
}

// Register your service(s)...
// pb.RegisterMyServiceServer(srv, handler)

if err := srv.Serve(lis); err != nil {
  log.Fatal().Err(err).Msg("grpc serve failed")
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVICE_NAME` | Service identifier | `unknown` |
| `LOG_LEVEL` | Log level (trace/debug/info/warn/error) | `info` |
| `LOG_FORMAT` | Log format (`json` for JSON, anything else for console) | console |
| `ENVIRONMENT` | Environment name | `development` |

## License

MIT

