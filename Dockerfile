# Builder Image
FROM golang:1.22 AS builder

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

# Start a new stage from scratch for the runtime image
FROM alpine:latest

# Install the CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory to where you'll run your app
WORKDIR /opt/echo-server

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /opt/echo-server/echo-server .

# Set Arguments that can be passed during build time
ARG PORT
ARG LOG_LEVEL

# Set default environment variables
ENV PORT=$PORT LOG_LEVEL=$LOG_LEVEL

# Execute the binary directly, ensuring to respect the ENV variables for configuration.
# Could also utilize a make target
CMD ["sh", "-c", "./echo-server -port=${PORT} -logLevel=${LOG_LEVEL}"]
