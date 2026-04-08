# ── Build stage ─────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for go mod download with some private repos)
RUN apk add --no-cache git

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# ── Runtime stage ──────────────────────────────────────
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates for HTTPS and tzdata for timezone
RUN apk add --no-cache ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Jakarta

# Copy the binary from the builder
COPY --from=builder /app/server .

# Expose the app port
EXPOSE 3000

# Run the binary
CMD ["./server"]
