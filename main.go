package main

import (
    "log"
    cfg "xenigo/internal/config"
    rdt "xenigo/internal/reddit"
)

func main() {
    appConfig, err := cfg.LoadConfig()
    if err != nil {
        log.Fatalf("Error reading config: %v", err)
    }

    log.Printf("Running in %s context", appConfig.Context)

    rdt.FetchData(appConfig)
}