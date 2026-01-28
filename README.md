## Echo-Server

<img src="images/echo-server.webp" alt="Echo Server Logo" width="400"/>

A simple echo server showcasing Go's standard library HTTP server and routing capabilities without external frameworks.

### Features

- HTTP Server with production-ready timeouts
- Routing using Go's native `http.ServeMux`
- API key authentication middleware
- Metrics middleware with status code capture
- Structured JSON logging via `slog`
- Prometheus metrics with path, method, and status labels
- 12-factor app configuration via environment variables
- Distroless Docker image for minimal attack surface
- Comprehensive unit tests with table-driven patterns

### Endpoints

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/health` | GET | No | Health check |
| `/metrics` | GET | No | Prometheus metrics |
| `/api/v1/echo` | POST | Yes* | Echo request body |

*When `AUTH_ENABLED=true`

### Quick Start

```bash
# Show available make targets
make

# Copy and configure environment
cp example.env .env

# Build and run locally
make build
make run

# Run with Docker
make docker-run

# Run tests
make test
```

### Configuration

All configuration is via environment variables (12-factor app compliant).

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `AUTH_ENABLED` | `false` | Enable API key authentication |
| `API_KEYS` | | Comma-separated list of valid API keys |

```bash
# Example: Run with authentication
AUTH_ENABLED=true API_KEYS="key1,key2" make run
```

### Authentication

When `AUTH_ENABLED=true`, protected endpoints (`/api/*`) require a valid API key in the `X-API-Key` header.

```bash
# Generate a secure API key
openssl rand -hex 32

# Request with API key
curl -X POST http://localhost:8080/api/v1/echo \
  -H "X-API-Key: your-api-key" \
  -d '{"message":"Hello World"}'
```

Unauthenticated requests return `401 Unauthorized`:
```json
{"error":"missing API key"}
```

Invalid keys return:
```json
{"error":"invalid API key"}
```

### Curl Examples

Health Check (no auth required):
```bash
curl http://localhost:8080/health
```

Echo (with auth):
```bash
curl -X POST http://localhost:8080/api/v1/echo \
  -H "X-API-Key: your-api-key" \
  -d '{"message":"Hello World"}'
```

Metrics (no auth required):
```bash
curl http://localhost:8080/metrics
```

### Docker

The Docker image uses a multi-stage build with a distroless runtime image for security.

```bash
# Build and run
make docker-run

# Run with auth enabled
AUTH_ENABLED=true API_KEYS="secret-key" make docker-run

# Cleanup
make docker-clean
```

### Project Structure

```
.
├── cmd/                      # Application entrypoint
├── internal/
│   ├── config/               # Environment configuration
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # Auth and metrics middleware
│   └── server/               # Server setup and routing
├── example.env               # Example environment file
├── Dockerfile                # Multi-stage distroless build
└── Makefile                  # Build and run targets
```

### Development

```bash
# Run tests
make test

# Run tests with coverage
make coverage

# Run linter
make lint

# Format code
make fmt
```

### Extending

To add endpoints:
1. Add handler in `internal/handlers/handlers.go`
2. Register route in `internal/server/server.go` `SetupRoutes()`
3. Add tests in `internal/handlers/handlers_test.go`
