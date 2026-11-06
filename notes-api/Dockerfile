# --------------------------------
# Stage 1: Build the application
# --------------------------------
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Copy dependency definitions first.
COPY go.mod go.sum ./

# Download dependencies.
RUN go mod download

# Copy the remaining source code.
COPY . .

# Compile the Go application.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /api \
    ./cmd/api


# --------------------------------
# Stage 2: Run the application
# --------------------------------
FROM alpine:3.22

WORKDIR /app

# Install trusted CA certificates for HTTPS connections.
RUN apk add --no-cache ca-certificates wget

# Create an unprivileged application user.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy only the compiled binary.
COPY --from=builder /api ./api

# Run as a non-root user.
USER appuser

# Default application port.
ENV PORT=8080

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:8080/health || exit 1

CMD ["./api"]