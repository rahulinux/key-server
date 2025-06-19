FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically-linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o key-server ./cmd/key-server

# --- Minimal runtime image using scratch ---
FROM scratch

# Set working dir
WORKDIR /

# Copy binary and any static assets
COPY --from=builder /app/key-server /

# Expose port
EXPOSE 8080

# Healthcheck (optional, remove if not needed)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/key-server", "-healthcheck"]

# Run binary
ENTRYPOINT ["/key-server"]
