# Builder Image
FROM golang:1.25.4 AS builder

# Create the directory to match the structure and set it as the working directory
WORKDIR /opt/echo-server

# Copy go.mod and go.sum files needed for dependancies
COPY go.mod go.sum ./

# Download all dependencies using the go mod tool
RUN go mod download

# Copy the entire project directory
COPY . .

# Change directory to the binary directory
WORKDIR /opt/echo-server/cmd

# Build the Go app
# Output the binary to the root of /opt/echo-server so it's easy to find in the next stage
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../echo-server .

# Start a new stage using distroless for minimal attack surface
FROM gcr.io/distroless/static:latest

# Set the working directory to where you'll run your app
WORKDIR /opt/echo-server

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /opt/echo-server/echo-server .

# Execute the binary - PORT and LOG_LEVEL are read from environment variables at runtime
ENTRYPOINT ["/opt/echo-server/echo-server"]
