# Build stage
FROM golang:1.23.0 AS builder

# Set working directory
WORKDIR /app

# Install git (required for some Go modules)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls and wget for health check
RUN apk --no-cache add ca-certificates tzdata wget

# Set timezone
ENV TZ=UTC

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy .env file if exists (optional)
COPY --from=builder /app/.env* ./

# Change ownership to non-root user
RUN chown -R appuser:appuser /app/
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run
CMD ["./main"]