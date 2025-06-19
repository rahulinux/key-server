package main

import (
    "flag"
    "fmt"
    "os"
    "log"
    "net/http"
    "github.com/gorilla/handlers"
)

type Config struct {
    SrvPort string
    MaxSize int
}

func ParseFlags(args []string) (Config, error) {
    fs := flag.NewFlagSet("app", flag.ContinueOnError)
    srvPort := fs.String("srv-port", "8080", "Port to run the server on")
    maxSize := fs.Int("max-size", 1024, "Maximum number of bytes allowed")
    if err := fs.Parse(args); err != nil {
        return Config{}, err
    }
    return Config{
        SrvPort: *srvPort,
        MaxSize: *maxSize,
    }, nil
}

func RunServer(cfg Config) error {
    router := NewRouter(cfg.MaxSize)
    loggedRouter := handlers.LoggingHandler(os.Stdout, router) 

    addr := fmt.Sprintf(":%s", cfg.SrvPort)
    log.Printf("Listening on %s\n", addr)
    return http.ListenAndServe(addr, loggedRouter)
}

func main() {
    cfg, err := ParseFlags(flag.CommandLine.Args())
    if err != nil {
      log.Fatal(err)
    }
    if err = RunServer(cfg); err != nil {
        log.Fatal(err)
    }
}

