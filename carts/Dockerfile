# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the Go modules files first to leverage layer caching
COPY go.mod go.sum ./

# Download modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /carts-service ./cmd/server

# Stage 2: Create the runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /carts-service .

# Copy migrations into the final image so they are available to migration tools
# (for example, when running migrations from inside the container). The builder
# places sources in /app, so migrations are at /app/migrations.
COPY --from=builder /app/migrations ./migrations

# Expose the port the service listens on
EXPOSE 8084

# Default migrations dir env points to the copied migrations so migration tools
# can access them inside the container.
ENV MIGRATIONS_DIR=/root/migrations

# Run the binary
CMD ["./carts-service"]
