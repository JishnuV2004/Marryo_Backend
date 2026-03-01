# ---------------------------
# Stage 1: Build the app
# ---------------------------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (required for go mod sometimes)
RUN apk add --no-cache git

# Copy go mod files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project
COPY . .

# Build the application from Cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./Cmd/main.go

# ---------------------------
# Stage 2: Run the app
# ---------------------------
FROM alpine:latest

WORKDIR /root/

# Install certificates (important for HTTPS, Firebase, etc.)
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/app .

# Copy Web templates (important for your SSR admin panel)
COPY --from=builder /app/Web ./Web

COPY --from=builder /app/Config ./Config

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./app"]