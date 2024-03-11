## Echo-Server

<img src="images/echo-server.webp" alt="Echo Server Logo" width="400"/>

This is a simple echo server that utilizes the new features in Go 1.22.
I wanted to showcase an example that now we do not need to use external libraries to handle the HTTP server and the routing.

### Features:

- [x] HTTP Server
- [x] Routing
- [x] Middleware
- [x] Structured Logging
- [x] Prometheus Metrics
- [x] Flag and Environment Variable Configuration

### Endpoints:

- `GET /health`: Returns the health of the server
- `POST api/v1/echo`: Returns the body of the request
- `GET /metrics`: Returns the metrics of the server

### Usage:

```bash
go run cmd/echo-service.go
```

### Configuration:

Note that environment variables for PORT and LOG_LEVEL take precedence over the flags.

### Make Native Go Execution:

```bash
make build
make run
```

#### Docker Execution:

```bash
make docker-run
```

#### Curl Examples:
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

### Cleanup - When done with Docker Execution:

```bash
make docker-clean
```

### Expansion:

To add on or remove an endpoint just manipulate this section under server/server.go

```go
func (s *Server) SetupRoutes() {}
```

Then add a handler for your route under handlers/handlers.go it's that simple.

This is to show that frameworks really are unnecessary for microservices with the new features in Go 1.22.