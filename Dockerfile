FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o key-server ./cmd/key-server

# Production image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Create non-root user
RUN addgroup -g 1001 -S keyserver && \
    adduser -u 1001 -S keyserver -G keyserver

# Copy binary
COPY --from=builder /app/key-server .

# Change ownership
RUN chown keyserver:keyserver key-server

# Switch to non-root user
USER keyserver

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

EXPOSE 8080

CMD ["./key-server"]
