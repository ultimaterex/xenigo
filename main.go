package main

import (
    "log"
    "time"
    "xenigo/internal/config"
    "xenigo/internal/reddit"
)

func main() {
    appConfig, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    var accessToken string
    if appConfig.Context == config.ContextElevated {
        accessToken, err = reddit.GetAccessToken(appConfig.Config.OAuth)
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

    // Log the startup information
    log.Println("Starting monitors with the following intervals:")
    for _, target := range appConfig.Config.Targets {
        log.Printf("Monitor: %s, Interval: %d seconds\n", target.Monitor.Subreddit, target.Options.Interval)
    }

    for _, target := range appConfig.Config.Targets {
        go monitorSubreddit(target, accessToken, appConfig.Config.UserAgent, string(appConfig.Context), cache, sendInitial)
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