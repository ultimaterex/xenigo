package main

import (
    "log"
    "time"
    cfg "xenigo/internal/config"
    "xenigo/internal/reddit"
)

func main() {
    appConfig, err := cfg.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }
    config := appConfig.Config

    var accessToken string
    if appConfig.Context == cfg.ContextElevated {
        accessToken, err = reddit.GetAccessToken(config.OAuth)
        if err != nil {
            log.Fatalf("Error getting access token: %v", err)
        }
    }

    // Load the cache
    cache := NewCache()
    if err := cache.Load("xenigo.cache"); err != nil {
        log.Fatalf("Error loading cache: %v", err)
    }

    // Determine if we should send the initial fetch data to Discord
    sendInitial := time.Since(cache.LastPersisted) <= 15*time.Minute
    if cfg.GetFlag(config.DeveloperFlags.ForceSendInitial) {
        sendInitial = true
    }

    // Log the startup information
    log.Println("Starting monitors with the following intervals:")
    for _, target := range config.Targets {
        log.Printf("Monitor: %s, Interval: %d seconds\n", target.Monitor.Subreddit, target.Options.Interval)
    }

    for _, target := range config.Targets {
        go monitorSubreddit(target, accessToken, config.UserAgent, string(appConfig.Context), cache, sendInitial, config.DeveloperFlags)
    }

    // Periodically save the cache
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()
        for range ticker.C {
            if err := cache.Save("xenigo.cache"); err != nil {
                log.Printf("Error saving cache: %v", err)
            }
        }
    }()

    // Prevent the main function from exiting
    select {}
}