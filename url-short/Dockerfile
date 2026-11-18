# Stage 1: Build the application         
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the remaining source code
COPY . .

# Compile the application
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /url-shortener ./cmd/api

# Stage 2: Small runtime image
FROM alpine:3.22

WORKDIR /app

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy compiled binary
COPY --from=builder /url-shortener .

# Application port
EXPOSE 8080

# Run the application
CMD ["./url-shortener"]
