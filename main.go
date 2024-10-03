package main

import (
    "log"
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

    // Log the startup information
    log.Println("Starting monitors with the following intervals:")
    for _, target := range appConfig.Config.Targets {
        log.Printf("Monitor: %s, Interval: %d seconds\n", target.Monitor.Subreddit, target.Options.Interval)
    }

    for _, target := range appConfig.Config.Targets {
        go monitorSubreddit(target, accessToken, appConfig.Config.UserAgent, string(appConfig.Context))
    }

    // Prevent the main function from exiting
    select {}
}