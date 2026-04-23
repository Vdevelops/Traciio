# Build stage
FROM golang:1.25.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates, netcat, dan postgresql-client untuk healthcheck dan DB setup
RUN apk --no-cache add ca-certificates netcat-openbsd postgresql-client

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server /app/server

# Copy entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN sed -i 's/\r$//' /app/entrypoint.sh && chmod +x /app/entrypoint.sh

# SQL initialization files (if any needed for initdb)
COPY database_indexes.sql/database_indexes.sql /app/database_indexes.sql

# Create uploads directory
RUN mkdir -p /app/uploads && chmod 755 /app/uploads

# Expose port (must match docker-compose.yml)
EXPOSE 8080

# Set environment untuk production defaults
ENV GIN_MODE=release \
    PORT=8080 \
    DB_MAX_OPEN_CONNS=500 \
    DB_MAX_IDLE_CONNS=100 \
    DB_CONN_MAX_LIFETIME=5m \
    DB_CONN_MAX_IDLE_TIME=10m \
    REDIS_POOL_SIZE=200 \
    REDIS_MIN_IDLE_CONNS=50

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Use entrypoint script
ENTRYPOINT ["/app/entrypoint.sh"]

