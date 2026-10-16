# Multi-stage Dockerfile for SatSet Backend (Go + Gin)

# Stage 1: Builder
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache ca-certificates tzdata git

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/satset-api ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:3.21

WORKDIR /app

# Install runtime dependencies (tzdata for Asia/Jakarta and wget for healthcheck)
RUN apk add --no-cache ca-certificates tzdata wget \
    && cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime \
    && echo "Asia/Jakarta" > /etc/timezone

# Prepare uploads folder
RUN mkdir -p /app/uploads

# Copy executable from builder stage
COPY --from=builder /app/satset-api /app/satset-api

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 -O - http://localhost:8080/healthz > /dev/null || exit 1

ENTRYPOINT ["/app/satset-api"]
