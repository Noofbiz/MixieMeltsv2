# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /users-service ./cmd/server

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /users-service .

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./users-service"]
