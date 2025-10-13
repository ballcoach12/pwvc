# Build stage
FROM golang:1.23.3-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies for building (including gcc for SQLite)
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o pairwise \
    ./cmd/server

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user
RUN addgroup -g 1001 pairwise && \
    adduser -D -u 1001 -G pairwise pairwise

# Set working directory
WORKDIR /home/pairwise

# Create data directory for SQLite
RUN mkdir -p /home/pairwise/data

# Copy binary from builder stage
COPY --from=builder /app/pairwise .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Set ownership
RUN chown -R pairwise:pairwise /home/pairwise

# Switch to non-root user
USER pairwise

# Expose port
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./pairwise"]