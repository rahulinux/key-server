package api

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

func KeyHandler(maxSize int) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        lengthStr, ok := vars["length"]
        if !ok {
            http.NotFound(w, r)
            return
        }
        length, err := strconv.Atoi(lengthStr)
        if err != nil || length <= 0 {
            http.Error(w, "Invalid length", http.StatusBadRequest)
            return
        }
        if length > maxSize {
            http.Error(w, "Requested length exceeds max-size", http.StatusBadRequest)
            return
        }

        bytes := make([]byte, length)
        if _, err := rand.Read(bytes); err != nil {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        resp := map[string]string{"key": hex.EncodeToString(bytes)}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }
}

