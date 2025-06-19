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
    fs.SetOutput(os.Stderr)
    srvPort := fs.String("srv-port", "8080", "Port to run the server on")
    maxSize := fs.Int("max-size", 1024, "Maximum number of bytes allowed")
    fs.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
        fs.PrintDefaults()
    }

    if err := fs.Parse(args); err != nil {
        if err == flag.ErrHelp {
            fs.Usage()
            os.Exit(0)
        }
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
    cfg, err := ParseFlags(os.Args[1:])
    if err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
    if err = RunServer(cfg); err != nil {
        log.Fatal(err)
    }
}

