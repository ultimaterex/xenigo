package main

import (
    "log"
    "time"
    cfg "xenigo/internal/config"
    "xenigo/internal/reddit"
)

const appVersion string = "0.3.2"

func main() {
    log.Printf("Hello world from xenigo! (version %s)", appVersion)

    // Ensure config.yaml exists and is usable
    configFile := "config.yaml"
    if err := cfg.EnsureConfigFile(configFile); err != nil {
        log.Fatalf("Error ensuring config file: %v", err)
    }

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

    // Ensure xenigo.cache exists and is usable
    cacheFile := "xenigo.cache"
    cache := NewCache()
    if err := cache.EnsureUsable(cacheFile); err != nil {
        log.Fatalf("Error ensuring cache is usable: %v", err)
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

    // Start monitoring
    for _, target := range config.Targets {
        go monitorSubreddit(target, accessToken, config.UserAgent, string(appConfig.Context), cache, sendInitial, config.DeveloperFlags, config.OAuth)
    }

    // Periodically save the cache
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()
        for range ticker.C {
            if err := cache.Save(cacheFile); err != nil {
                log.Printf("Error saving cache: %v", err)
            }
        }
    }()

    // Prevent the main function from exiting
    select {}
}

