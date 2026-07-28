# =====================================================
# BUILD STAGE
# =====================================================

FROM golang:1.25.4-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -a \
    -installsuffix cgo \
    -o server \
    ./cmd/server/main.go


# =====================================================
# RUNTIME STAGE
# =====================================================

FROM alpine:latest

RUN apk --no-cache add \
    ca-certificates \
    wget

WORKDIR /app

COPY --from=builder /app/server /app/server

RUN mkdir -p /app/uploads \
    && chmod 755 /app/uploads

EXPOSE 8080

ENV GIN_MODE=release \
    PORT=8080 \
    ENV=production \
    DB_MAX_OPEN_CONNS=100 \
    DB_MAX_IDLE_CONNS=25 \
    DB_CONN_MAX_LIFETIME=5m \
    DB_CONN_MAX_IDLE_TIME=10m \
    REDIS_POOL_SIZE=50 \
    REDIS_MIN_IDLE_CONNS=10

HEALTHCHECK \
    --interval=30s \
    --timeout=10s \
    --start-period=10s \
    --retries=3 \
    CMD wget \
        --no-verbose \
        --tries=1 \
        --spider \
        http://127.0.0.1:8080/health \
        || exit 1

CMD ["/app/server"]