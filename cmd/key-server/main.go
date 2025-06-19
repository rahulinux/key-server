package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "github.com/rahulinux/key-server/internal/api"

    "github.com/gorilla/mux"
)

func main() {
    srvPort := flag.String("srv-port", "8080", "Port to run the server on")
    maxSize := flag.Int("max-size", 1024, "Maximum number of bytes allowed")
    flag.Parse()

    r := mux.NewRouter()
    r.HandleFunc("/key/{length}", api.KeyHandler(*maxSize)).Methods("GET")

    addr := fmt.Sprintf(":%s", *srvPort)
    log.Printf("Listening on %s\n", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatal(err)
    }
}

