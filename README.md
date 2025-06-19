# 🔐 key-server

A simple HTTP service to generate secure random keys of a given length. Useful for APIs, tokens, or one-time secrets. Includes Prometheus metrics and structured logging.

---

## 🚀 Features

- `/key/{length}`: Generate a hex-encoded random key representing {length} random bytes
- `/healthz`: Basic health check endpoint
- `/metrics`: Prometheus metrics (request durations, HTTP status codes, key lengths)
- Graceful shutdown with `SIGINT`/`SIGTERM`
- Configurable via CLI flags or environment variables

---

## 🛠️ Quick Start

### Build & Run (Go 1.21+)
```bash
go build -o key-server ./cmd/key-server
./key-server --srv-port=8080 --max-size=1024 --log-level=info
````

Or use Docker:

```bash
docker build -t key-server .
docker run -p 8080:8080 key-server
```

---

## ⚙️ Configuration

You can configure `key-server` using flags or environment variables.

| Flag          | Env Variable  | Default | Description                     |
| ------------- | ------------- | ------- | ------------------------------- |
| `--srv-port`  | `SERVER_PORT` | `8080`  | HTTP server port                |
| `--max-size`  | `MAX_SIZE`    | `1024`  | Max allowed key size (in bytes) |
| `--log-level` | `LOG_LEVEL`   | `info`  | Log level: debug, info, warn    |

---

## 📡 API Endpoints

| Method | Path            | Description                        |
| ------ | --------------- | ---------------------------------- |
| GET    | `/key/{length}` | Generate a hex-encoded random key representing {length} random bytes (key string length = 2 × {length})  |
| GET    | `/healthz`      | Health check                       |
| GET    | `/metrics`      | Prometheus metrics                 |

Note: The `key` field is a hex string. Its length is always `2 × length` because each byte is represented by two hex characters.

Example:

```bash
curl http://localhost:8080/key/16
```

Output:

```
{
  "key": "e3b0c44298fc1c149afbf4c8996fb924", // 32 hex chars for 16 bytes
  "length": 16
}
```

---

## 🧪 Running Tests

```bash
go test ./...
```

Includes table-driven tests for `/key/{length}` behavior and metrics initialization.

---

## 🧰 Developer Notes

* Prometheus histograms are carefully guarded against invalid bucket ordering.
* Logging uses Go’s `log/slog` with structured JSON output.
* Safe shutdown using `context.WithTimeout()` and signal listening.

---

## 📦 Deployment

Use the included Dockerfile:

```bash
docker build -t key-server .
docker run -p 8080:8080 key-server
```

Health check is already baked in:

```dockerfile
HEALTHCHECK CMD wget --spider --quiet http://localhost:8080/healthz || exit 1
```

Helm deployment

```bash
helm install key-server ./key-server

# OR

helm install key-server ./key-server \
  --set image.repository=yourrepo/key-server,image.tag=latest
```

