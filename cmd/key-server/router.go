package main

import (
    "github.com/gorilla/mux"
    "github.com/rahulinux/key-server/internal/api"
)

func NewRouter(maxSize int) *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/key/{length}", api.KeyHandler(maxSize)).Methods("GET")
    return r
}

