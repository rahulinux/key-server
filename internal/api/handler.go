package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rahulinux/key-server/internal/metrics"
)

type KeyHandler struct {
	maxSize int
	logger  *slog.Logger
}

type HealthHandler struct {
	logger *slog.Logger
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type KeyResponse struct {
	Key    string `json:"key"`
	Length int    `json:"length"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

func NewKeyHandler(maxSize int, logger *slog.Logger) *KeyHandler {
	return &KeyHandler{
		maxSize: maxSize,
		logger:  logger,
	}
}

func NewHealthHandler(logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

func (h *KeyHandler) HandleKey(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	vars := mux.Vars(r)
	lengthStr := vars["length"]

	// Parse and validate length
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid length parameter", err)
		return
	}

	if length <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "Length must be positive", nil)
		return
	}

	if length > h.maxSize {
		h.respondWithError(w, http.StatusBadRequest,
			fmt.Sprintf("Requested length %d exceeds maximum allowed size %d", length, h.maxSize), nil)
		return
	}

	// Generate random bytes
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		h.logger.Error("Failed to generate random bytes", "error", err, "length", length)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to generate key", err)
		return
	}

	// Create response
	key := hex.EncodeToString(bytes)
	response := KeyResponse{
		Key:    key,    // full hex string, length = 2 * bytes
		Length: length, // number of random bytes requested
	}

	// Record metrics
	metrics.KeyLengthHistogram.Observe(float64(length))
	metrics.HTTPStatusCounter.WithLabelValues("200").Inc()
	metrics.RequestDuration.WithLabelValues("key").Observe(time.Since(start).Seconds())

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
	}

	h.logger.Debug("Generated key", "length", length, "duration", time.Since(start))
}

func (h *KeyHandler) respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	metrics.HTTPStatusCounter.WithLabelValues(strconv.Itoa(statusCode)).Inc()

	errorResp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if encodeErr := json.NewEncoder(w).Encode(errorResp); encodeErr != nil {
		h.logger.Error("Failed to encode error response", "error", encodeErr)
	}

	if err != nil {
		h.logger.Error("Request failed", "status", statusCode, "message", message, "error", err)
	} else {
		h.logger.Warn("Request failed", "status", statusCode, "message", message)
	}
}

func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode health response", "error", err)
	}
}
