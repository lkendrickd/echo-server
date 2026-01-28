## Echo-Server

<img src="images/echo-server.webp" alt="Echo Server Logo" width="400"/>

A simple echo server showcasing Go's standard library HTTP server and routing capabilities without external frameworks.

### Features

- HTTP Server with production-ready timeouts
- Routing using Go's native `http.ServeMux`
- Metrics middleware with status code capture
- Structured JSON logging via `slog`
- Prometheus metrics with path, method, and status labels
- Flag and environment variable configuration
- Distroless Docker image for minimal attack surface
- Unit tests for handlers, middleware, and server

### Endpoints

- `GET /health` - Returns server health status
- `POST /api/v1/echo` - Returns the request body
- `GET /metrics` - Returns Prometheus metrics

### Quick Start

```bash
# Show available make targets
make

# Build and run locally
make build
make run

# Run with Docker
make docker-run

# Run tests
make test
```

### Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |

Environment variables take precedence over flags.

```bash
# Override defaults
PORT=9090 LOG_LEVEL=debug make run
```

### Curl Examples

Health Check:
```bash
curl http://localhost:8080/health
```

Echo:
```bash
curl -X POST http://localhost:8080/api/v1/echo -d '{"message":"Hello World"}'
```

Metrics:
```bash
curl http://localhost:8080/metrics
```

### Docker

The Docker image uses a multi-stage build with a distroless runtime image for security.

```bash
# Build and run
make docker-run

# Cleanup
make docker-clean
```

### Project Structure

```
.
├── cmd/                      # Application entrypoint
├── internal/
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # Metrics middleware
│   └── server/               # Server setup and routing
├── Dockerfile                # Multi-stage distroless build
└── Makefile                  # Build and run targets
```

### Extending

To add endpoints, modify `SetupRoutes()` in `internal/server/server.go` and add handlers in `internal/handlers/handlers.go`.
